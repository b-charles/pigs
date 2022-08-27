package json

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"unicode/utf16"
	"unicode/utf8"
)

var (
	runesHex        = []rune(`0123456789abcdefABCDEF`)
	runesEscaped    = []rune(`"\/bfnrtu`)
	runesPlane      = []rune(`\u`)
	runesMinusDigit = []rune(`-0123456789`)
	runesDigit      = []rune(`0123456789`)
	runesTrue       = []rune(`true`)
	runesFalse      = []rune(`false`)
	runesNull       = []rune(`null`)
	runesObjectSep  = []rune(`,}`)
	runesArraySep   = []rune(`,]`)
)

type jsonParser struct {
	source      io.RuneReader
	currentPos  int
	currentRune rune
	previous    []rune
}

func newParser(source io.RuneReader) *jsonParser {
	parser := &jsonParser{
		source:     source,
		currentPos: -1,
		previous:   []rune{},
	}
	parser.next()
	return parser
}

func Parse(source io.RuneReader) (JsonNode, error) {
	parser := newParser(source)
	return parser.readJsonNode()
}

func ParseAll(source io.RuneReader) ([]JsonNode, error) {

	parser := newParser(source)

	nodes := []JsonNode{}

	for !parser.eof() {
		if n, err := parser.readJsonNode(); err != nil {
			return []JsonNode{}, err
		} else {
			nodes = append(nodes, n)
		}
		parser.skipWS()
	}

	return nodes, nil

}

func ParseString(source string) (JsonNode, error) {
	return Parse(strings.NewReader(source))
}

func ParseAllString(source string) ([]JsonNode, error) {
	return ParseAll(strings.NewReader(source))
}

func (self *jsonParser) wrapError(err error) error {
	return fmt.Errorf("Error at %d '%s': %w", self.currentPos, string(self.previous), err)
}

func (self *jsonParser) wrap(msg string, args ...any) error {
	return self.wrapError(fmt.Errorf(msg, args...))
}

func (self *jsonParser) next() rune {
	if self.eof() {
		return self.currentRune
	}
	if r, _, err := self.source.ReadRune(); err != nil {
		self.currentRune = utf8.RuneError
	} else {
		self.currentRune = r
		self.currentPos++
		if len(self.previous) < 25 {
			self.previous = append(self.previous, r)
		} else {
			self.previous = append(self.previous[1:], r)
		}
	}
	return self.currentRune
}

func (self *jsonParser) eof() bool {
	return self.currentRune == utf8.RuneError
}

func (self *jsonParser) expectAny(r []rune) (rune, error) {
	if self.eof() {
		return utf8.RuneError, self.wrap("Unexpected end of file (expected '%v').", string(r))
	}
	for _, e := range r {
		if e == self.currentRune {
			self.next()
			return e, nil
		}
	}
	return utf8.RuneError, self.wrap("Unexpected rune: '%v' (expected '%v')", string(self.currentRune), string(r))
}

func (self *jsonParser) expect(r rune) error {
	if self.eof() {
		return self.wrap("Unexpected end of file (expected '%v').", string(r))
	}
	if r != self.currentRune {
		return self.wrap("Unexpected rune: '%v' (expected '%v')", string(self.currentRune), string(r))
	} else {
		self.next()
		return nil
	}
}

func (self *jsonParser) expectAll(r []rune) error {
	for _, e := range r {
		if err := self.expect(e); err != nil {
			return err
		}
	}
	return nil
}

func (self *jsonParser) readHex() (int, error) {
	if h, err := self.expectAny(runesHex); err != nil {
		return 0, err
	} else if '0' <= h && h <= '9' {
		return int(h - '0'), nil
	} else if 'a' <= h && h <= 'f' {
		return 10 + int(h-'a'), nil
	} else { // if 'A' <= h && h <= 'F'
		return 10 + int(h-'A'), nil
	}
}

func (self *jsonParser) readCodePoint() (rune, error) {
	i := 0
	for n := 0; n < 4; n++ {
		if h, err := self.readHex(); err != nil {
			return utf8.RuneError, err
		} else {
			i = (i << 4) + h
		}
	}
	return rune(i), nil
}

func (self *jsonParser) readEscapedRune() (rune, error) {

	r := self.currentRune
	self.next()

	if r != '\\' {
		return r, nil
	}

	if e, err := self.expectAny(runesEscaped); err != nil {
		return utf8.RuneError, err
	} else {
		switch e {
		case '"':
			return '"', nil
		case '\\':
			return '\\', nil
		case '/':
			return '/', nil
		case 'b':
			return '\b', nil
		case 'f':
			return '\f', nil
		case 'n':
			return '\n', nil
		case 'r':
			return '\r', nil
		case 't':
			return '\t', nil
		case 'u':

			if r1, err := self.readCodePoint(); err != nil {
				return utf8.RuneError, err
			} else if utf16.IsSurrogate(r1) {

				self.expectAll(runesPlane)
				if r2, err := self.readCodePoint(); err != nil {
					return utf8.RuneError, err
				} else {
					return utf16.DecodeRune(r1, r2), nil
				}

			} else {
				return r1, nil
			}

		default:
			return utf8.RuneError, self.wrap("Unreachable")
		}
	}

}

func (self *jsonParser) readEscapedString() (string, error) {

	if err := self.expect('"'); err != nil {
		return "", err
	}

	var builder strings.Builder

	for self.currentRune != '"' {
		if self.eof() {
			return "", self.wrap("Unexpected end of file (expected '\"').")
		} else if r, err := self.readEscapedRune(); err != nil {
			return "", err
		} else {
			builder.WriteRune(r)
		}
	}

	self.next() // "

	return builder.String(), nil

}

func (self *jsonParser) startString() bool {
	return self.currentRune == '"'
}

func (self *jsonParser) readString() (JsonNode, error) {
	if str, err := self.readEscapedString(); err != nil {
		return JSON_NULL, err
	} else {
		return JsonString(str), nil
	}
}

func (self *jsonParser) readDigits(builder *strings.Builder) error {
	if self.eof() {
		return self.wrap("Unexpected end of file (expected '0123456789').")
	}
	for '0' <= self.currentRune && self.currentRune <= '9' {
		builder.WriteRune(self.currentRune)
		self.next()
	}
	return nil
}

func (self *jsonParser) startNumber() bool {
	return self.currentRune == '-' || ('0' <= self.currentRune && self.currentRune <= '9')
}

func (self *jsonParser) readNumber() (JsonNode, error) {

	var builder strings.Builder
	isInt := true

	if r, err := self.expectAny(runesMinusDigit); err != nil {
		return JSON_NULL, err
	} else {

		builder.WriteRune(r)

		if r == '-' {
			if r, err = self.expectAny(runesDigit); err != nil {
				return JSON_NULL, err
			}
			builder.WriteRune(r)
		}

		if r != '0' {
			for '0' <= self.currentRune && self.currentRune <= '9' {
				builder.WriteRune(self.currentRune)
				self.next()
			}
		}

	}

	if self.currentRune == '.' {
		isInt = false
		builder.WriteRune(self.currentRune)
		self.next()
		if err := self.readDigits(&builder); err != nil {
			return JSON_NULL, err
		}
	}

	if self.currentRune == 'e' || self.currentRune == 'E' {
		isInt = false
		builder.WriteRune(self.currentRune)
		self.next()
		if self.eof() {
			return JSON_NULL, self.wrap("Unexpected end of file (expected '+-0123456789').")
		}
		if self.currentRune == '+' || self.currentRune == '-' {
			builder.WriteRune(self.currentRune)
			self.next()
		}
		if err := self.readDigits(&builder); err != nil {
			return JSON_NULL, err
		}
	}

	str := builder.String()

	if isInt {
		if val, err := strconv.Atoi(str); err != nil {
			return JSON_NULL, self.wrapError(err)
		} else {
			return JsonInt(val), nil
		}
	} else {
		if val, err := strconv.ParseFloat(str, 64); err != nil {
			return JSON_NULL, self.wrapError(err)
		} else {
			return JsonFloat(val), nil
		}
	}

}

func (self *jsonParser) startConst() bool {
	return self.currentRune == 't' || self.currentRune == 'f' || self.currentRune == 'n'
}

func (self *jsonParser) readConst() (JsonNode, error) {

	switch self.currentRune {
	case 't':
		if err := self.expectAll(runesTrue); err != nil {
			return JSON_NULL, err
		} else {
			return JSON_TRUE, nil
		}
	case 'f':
		if err := self.expectAll(runesFalse); err != nil {
			return JSON_NULL, err
		} else {
			return JSON_FALSE, nil
		}
	case 'n':
		if err := self.expectAll(runesNull); err != nil {
			return JSON_NULL, err
		} else {
			return JSON_NULL, nil
		}
	default:
		return JSON_NULL, self.wrap("Unreachable")
	}

}

func (self *jsonParser) skipWS() {
	for self.currentRune == '\x20' || self.currentRune == '\x09' ||
		self.currentRune == '\x0a' || self.currentRune == '\x0d' {
		self.next()
	}
}

func (self *jsonParser) startObject() bool {
	return self.currentRune == '{'
}

func (self *jsonParser) readObject() (JsonNode, error) {

	var (
		s     rune
		err   error
		key   string
		value JsonNode
	)

	if err = self.expect('{'); err != nil {
		return JSON_NULL, err
	}

	members := newEmptySortedMap[JsonNode]()

	self.skipWS()

	s = self.currentRune
	for s != '}' {

		key, err = self.readEscapedString()
		if err != nil {
			return JSON_NULL, err
		}

		self.skipWS()

		err = self.expect(':')
		if err != nil {
			return JSON_NULL, err
		}

		self.skipWS()

		value, err = self.readJsonNode()
		if err != nil {
			return JSON_NULL, err
		}

		members.put(key, value)

		self.skipWS()

		s, err = self.expectAny(runesObjectSep)
		if err != nil {
			return JSON_NULL, err
		} else if s == ',' {
			self.skipWS()
		}

	}

	return &JsonObject{members}, nil

}

func (self *jsonParser) startArray() bool {
	return self.currentRune == '['
}

func (self *jsonParser) readArray() (JsonNode, error) {

	var (
		s     rune
		err   error
		value JsonNode
	)

	if err = self.expect('['); err != nil {
		return JSON_NULL, err
	}

	elements := make([]JsonNode, 0, 10)

	self.skipWS()

	s = self.currentRune
	for s != ']' {

		value, err = self.readJsonNode()
		if err != nil {
			return JSON_NULL, err
		}

		elements = append(elements, value)

		self.skipWS()

		s, err = self.expectAny(runesArraySep)
		if err != nil {
			return JSON_NULL, err
		} else if s == ',' {
			self.skipWS()
		}

	}

	return &JsonArray{elements}, nil

}

func (self *jsonParser) readJsonNode() (JsonNode, error) {

	self.skipWS()

	if self.startString() {
		return self.readString()
	} else if self.startObject() {
		return self.readObject()
	} else if self.startArray() {
		return self.readArray()
	} else if self.startConst() {
		return self.readConst()
	} else if self.startNumber() {
		return self.readNumber()
	} else {
		return JSON_NULL, self.wrap("Unexpected rune: '%v' (expected '\"{[tfn-0123456789')", string(self.currentRune))
	}

}

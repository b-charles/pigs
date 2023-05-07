package json

import (
	"unicode/utf16"
	"unicode/utf8"
)

type JsonString string

func (self JsonString) IsString() bool {
	return true
}

func (self JsonString) AsString() string {
	return string(self)
}

func (self JsonString) IsFloat() bool {
	return false
}

func (self JsonString) AsFloat() float64 {
	return 0.0
}

func (self JsonString) IsInt() bool {
	return false
}

func (self JsonString) AsInt() int {
	return 0
}

func (self JsonString) IsBool() bool {
	return false
}

func (self JsonString) AsBool() bool {
	return false
}

func (self JsonString) IsObject() bool {
	return false
}

func (self JsonString) GetKeys() []string {
	return []string{}
}

func (self JsonString) GetMember(string) JsonNode {
	return JSON_NULL
}

func (self JsonString) IsArray() bool {
	return false
}

func (self JsonString) GetLen() int {
	return 0
}

func (self JsonString) GetElement(int) JsonNode {
	return JSON_NULL
}

func (self JsonString) IsNull() bool {
	return false
}

func (self JsonString) Append(b []byte) []byte {
	return appendString(b, string(self))
}

func (self JsonString) String() string {
	b := make([]byte, 0, len(string(self))+2)
	b = self.Append(b)
	return string(b)
}

var hex = `0123456789abcdef`

var escaped = map[byte][]byte{
	'\b': []byte(`\b`),
	'\f': []byte(`\f`),
	'\n': []byte(`\n`),
	'\r': []byte(`\r`),
	'\t': []byte(`\t`),
}

func appendString(buf []byte, s string) []byte {

	buf = append(buf, '"')

	for i := 0; i < len(s); {

		b := s[i]

		if 0x020 <= b && b < utf8.RuneSelf {
			if b == '\\' || b == '"' {
				buf = append(buf, '\\')
			}
			buf = append(buf, b)
			i++
			continue
		}

		if esc, ok := escaped[b]; ok {
			buf = append(buf, esc...)
			i++
			continue
		}

		c, size := utf8.DecodeRuneInString(s[i:])

		if c == utf8.RuneError && size == 1 {

			buf = append(buf, []byte("\\ufffd")...)

		} else if size < 4 {

			buf = append(buf, []byte("\\u")...)
			buf = append(buf, hex[(c>>12)&0xF])
			buf = append(buf, hex[(c>>8)&0xF])
			buf = append(buf, hex[(c>>4)&0xF])
			buf = append(buf, hex[c&0xF])

		} else {

			c1, c2 := utf16.EncodeRune(c)

			buf = append(buf, []byte("\\u")...)
			buf = append(buf, hex[(c1>>12)&0xF])
			buf = append(buf, hex[(c1>>8)&0xF])
			buf = append(buf, hex[(c1>>4)&0xF])
			buf = append(buf, hex[c1&0xF])

			buf = append(buf, []byte("\\u")...)
			buf = append(buf, hex[(c2>>12)&0xF])
			buf = append(buf, hex[(c2>>8)&0xF])
			buf = append(buf, hex[(c2>>4)&0xF])
			buf = append(buf, hex[c2&0xF])

		}

		i += size

	}

	buf = append(buf, '"')

	return buf

}

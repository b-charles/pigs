package json

import (
	"strings"
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

func (self JsonString) String() string {
	return formatString(string(self))
}

var hex = `0123456789abcdef`

var escaped = map[byte]string{
	'\b': `\b`,
	'\f': `\f`,
	'\n': `\n`,
	'\r': `\r`,
	'\t': `\t`,
}

func formatString(s string) string {

	var builder strings.Builder
	builder.WriteByte('"')

	for i := 0; i < len(s); {

		b := s[i]

		if 0x020 <= b && b < utf8.RuneSelf {
			if b == '\\' || b == '"' {
				builder.WriteByte('\\')
			}
			builder.WriteByte(b)
			i++
			continue
		}

		if esc, ok := escaped[b]; ok {
			builder.WriteString(esc)
			i++
			continue
		}

		c, size := utf8.DecodeRuneInString(s[i:])

		if c == utf8.RuneError && size == 1 {

			builder.WriteString(`\ufffd`)

		} else if size < 4 {

			builder.WriteString(`\u`)
			builder.WriteByte(hex[(c>>12)&0xF])
			builder.WriteByte(hex[(c>>8)&0xF])
			builder.WriteByte(hex[(c>>4)&0xF])
			builder.WriteByte(hex[c&0xF])

		} else {

			c1, c2 := utf16.EncodeRune(c)

			builder.WriteString(`\u`)
			builder.WriteByte(hex[(c1>>12)&0xF])
			builder.WriteByte(hex[(c1>>8)&0xF])
			builder.WriteByte(hex[(c1>>4)&0xF])
			builder.WriteByte(hex[c1&0xF])

			builder.WriteString(`\u`)
			builder.WriteByte(hex[(c2>>12)&0xF])
			builder.WriteByte(hex[(c2>>8)&0xF])
			builder.WriteByte(hex[(c2>>4)&0xF])
			builder.WriteByte(hex[c2&0xF])

		}

		i += size

	}

	builder.WriteByte('"')
	return builder.String()

}

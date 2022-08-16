package core

import (
	"fmt"
	"strings"
)

type JsonArray struct {
	elements []JsonNode
}

func NewJsonArray(elements []JsonNode) *JsonArray {
	return &JsonArray{elements}
}

var JSON_EMPTY_ARRAY = &JsonArray{[]JsonNode{}}

func (self *JsonArray) IsString() bool {
	return false
}

func (self *JsonArray) AsString() string {
	return ""
}

func (self *JsonArray) IsFloat() bool {
	return false
}

func (self *JsonArray) AsFloat() float64 {
	return 0.0
}

func (self *JsonArray) IsInt() bool {
	return false
}

func (self *JsonArray) AsInt() int {
	return 0
}

func (self *JsonArray) IsBool() bool {
	return false
}

func (self *JsonArray) AsBool() bool {
	return false
}

func (self *JsonArray) IsObject() bool {
	return false
}

func (self *JsonArray) GetKeys() []string {
	return []string{}
}

func (self *JsonArray) GetMember(string) JsonNode {
	return JSON_NULL
}

func (self *JsonArray) IsArray() bool {
	return true
}

func (self *JsonArray) GetLen() int {
	return len(self.elements)
}

func (self *JsonArray) GetElement(i int) JsonNode {
	if i < 0 || i >= len(self.elements) {
		return JSON_NULL
	}
	if e := self.elements[i]; e == nil {
		return JSON_NULL
	} else {
		return e
	}
}

func (self *JsonArray) IsNull() bool {
	return false
}

func (self *JsonArray) String() string {

	l := len(self.elements)
	if l == 0 {
		return "[]"
	}

	var b strings.Builder
	fmt.Fprint(&b, "[")

	b.WriteString(self.GetElement(0).String())
	for i := 1; i < l; i++ {
		b.WriteRune(',')
		b.WriteString(self.GetElement(i).String())
	}

	fmt.Fprint(&b, "]")
	return b.String()

}

package core

import (
	"fmt"
	"math"
	"strings"
)

type JsonArray struct {
	elements []JsonNode
}

func newJsonArray() *JsonArray {
	return &JsonArray{make([]JsonNode, 0, 10)}
}

var JSON_EMPTY_ARRAY = newJsonArray()

func (self *JsonArray) set(i int, elt JsonNode) {

	if i < 0 {
		panic(fmt.Errorf("Invalid argument: index %d must not be negative.", i))
	}

	if cap(self.elements) <= i {

		c := int(math.Exp2(math.Ceil(math.Log2(float64(i + 1)))))
		l := int(math.Max(float64(len(self.elements)), float64(i+1)))

		newSlice := make([]JsonNode, l, c)
		for i, e := range self.elements {
			newSlice[i] = e
		}
		self.elements = newSlice

	}

	if len(self.elements) <= i {
		self.elements = self.elements[:i+1]
	}

	self.elements[i] = elt

}

func (self *JsonArray) append(elt JsonNode) {
	self.elements = append(self.elements, elt)
}

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

	l := self.GetLen()
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

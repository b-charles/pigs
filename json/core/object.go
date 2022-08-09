package core

import (
	"fmt"
	"strings"
)

type JsonObject struct {
	members map[string]JsonNode
	sorted  []string
}

func newJsonObject() *JsonObject {
	return &JsonObject{map[string]JsonNode{}, []string{}}
}

var JSON_EMPTY_OBJECT = newJsonObject()

func (self *JsonObject) set(key string, elt JsonNode) {
	if _, p := self.members[key]; !p {
		self.sorted = append(self.sorted, key)
	}
	self.members[key] = elt
}

func (self *JsonObject) IsString() bool {
	return false
}

func (self *JsonObject) AsString() string {
	return ""
}

func (self *JsonObject) IsFloat() bool {
	return false
}

func (self *JsonObject) AsFloat() float64 {
	return 0.0
}

func (self *JsonObject) IsInt() bool {
	return false
}

func (self *JsonObject) AsInt() int {
	return 0
}

func (self *JsonObject) IsBool() bool {
	return false
}

func (self *JsonObject) AsBool() bool {
	return false
}

func (self *JsonObject) IsObject() bool {
	return true
}

func (self *JsonObject) GetKeys() []string {
	return self.sorted
}

func (self *JsonObject) GetMember(key string) JsonNode {
	if elt, p := self.members[key]; !p {
		return JSON_NULL
	} else {
		return elt
	}
}

func (self *JsonObject) IsArray() bool {
	return false
}

func (self *JsonObject) GetLen() int {
	return 0
}

func (self *JsonObject) GetElement(int) JsonNode {
	return JSON_NULL
}

func (self *JsonObject) IsNull() bool {
	return false
}

func (self *JsonObject) String() string {

	if len(self.sorted) == 0 {
		return "{}"
	}

	var b strings.Builder
	fmt.Fprint(&b, "{")

	b.WriteString(formatString(self.sorted[0]))
	b.WriteRune(':')
	b.WriteString(self.members[self.sorted[0]].String())

	for _, k := range self.sorted[1:] {
		b.WriteRune(',')
		b.WriteString(formatString(k))
		b.WriteRune(':')
		b.WriteString(self.members[k].String())
	}

	fmt.Fprint(&b, "}")
	return b.String()

}

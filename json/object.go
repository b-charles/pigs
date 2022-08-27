package json

import (
	"fmt"
	"strings"
)

type JsonObject struct {
	members *sortedMap[JsonNode]
}

func NewJsonObjectSorted(members map[string]JsonNode, sorted []string) *JsonObject {
	return &JsonObject{newSortedMap(members, sorted)}
}

func NewJsonObject(members map[string]JsonNode) *JsonObject {

	sorted := make([]string, 0, len(members))
	for k := range members {
		sorted = append(sorted, k)
	}

	return &JsonObject{&sortedMap[JsonNode]{members, sorted}}

}

var JSON_EMPTY_OBJECT = &JsonObject{&sortedMap[JsonNode]{map[string]JsonNode{}, []string{}}}

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
	return self.members.keys()
}

func (self *JsonObject) GetMember(key string) JsonNode {
	if elt, p := self.members.get(key); !p {
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

	if self.members.len() == 0 {
		return "{}"
	}

	var b strings.Builder
	fmt.Fprint(&b, "{")

	for i, k := range self.members.s {

		if i > 0 {
			b.WriteRune(',')
		}

		b.WriteString(formatString(k))
		b.WriteRune(':')
		b.WriteString(self.GetMember(k).String())

	}

	fmt.Fprint(&b, "}")
	return b.String()

}

package json

import (
	"sort"
)

type JsonObject struct {
	members *sortedMap[JsonNode]
}

func newJsonObject(members map[string]JsonNode, sorted []string) *JsonObject {
	return &JsonObject{&sortedMap[JsonNode]{members, sorted}}
}

func NewJsonObjectMappedSorted[T any](
	members map[string]T,
	mapper func(v T) JsonNode,
	less func(a, b string) bool) *JsonObject {

	sorted := make([]string, 0, len(members))
	for k := range members {
		sorted = append(sorted, k)
	}
	sort.Slice(sorted, func(i, j int) bool {
		return less(sorted[i], sorted[j])
	})

	nodes := make(map[string]JsonNode, len(members))
	for k := range members {
		nodes[k] = mapper(members[k])
	}

	return newJsonObject(nodes, sorted)

}

func NewJsonObjectMapped[T any](members map[string]T, mapper func(v T) JsonNode) *JsonObject {
	return NewJsonObjectMappedSorted(members, mapper, stringsLess)
}

func NewJsonObjectSorted(members map[string]JsonNode, less func(a, b string) bool) *JsonObject {
	return NewJsonObjectMappedSorted(members, jsonToJson, less)
}

func NewJsonObject(members map[string]JsonNode) *JsonObject {
	return NewJsonObjectMappedSorted(members, jsonToJson, stringsLess)
}

func NewJsonObjectStringsSorted(members map[string]string, less func(a, b string) bool) *JsonObject {
	return NewJsonObjectMappedSorted(members, stringToJson, less)
}

func NewJsonObjectStrings(members map[string]string) *JsonObject {
	return NewJsonObjectMappedSorted(members, stringToJson, stringsLess)
}

var JSON_EMPTY_OBJECT = newJsonObject(map[string]JsonNode{}, []string{})

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

func (self *JsonObject) Append(b []byte) []byte {

	if self.members.len() == 0 {
		return append(b, '{', '}')
	}

	b = append(b, '{')

	for i, k := range self.members.s {

		if i > 0 {
			b = append(b, ',')
		}

		b = appendString(b, k)
		b = append(b, ':')
		b = self.GetMember(k).Append(b)

	}

	b = append(b, '}')

	return b

}

func (self *JsonObject) String() string {
	b := make([]byte, 0, 50)
	b = self.Append(b)
	return string(b)
}

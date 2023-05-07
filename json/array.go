package json

type JsonArray struct {
	elements []JsonNode
}

func NewJsonArrayMapped[T any](elements []T, mapper func(v T) JsonNode) *JsonArray {
	elts := make([]JsonNode, 0, len(elements))
	for _, k := range elements {
		elts = append(elts, mapper(k))
	}
	return &JsonArray{elts}
}

func NewJsonArray(elements []JsonNode) *JsonArray {
	return &JsonArray{elements}
}

func NewJsonArrayStrings(elements []string) *JsonArray {
	return NewJsonArrayMapped(elements, stringToJson)
}

func NewJsonArrayFloats(elements []float64) *JsonArray {
	return NewJsonArrayMapped(elements, floatToJson)
}

func NewJsonArrayInts(elements []int) *JsonArray {
	return NewJsonArrayMapped(elements, intToJson)
}

func NewJsonArrayBools(elements []bool) *JsonArray {
	return NewJsonArrayMapped(elements, boolToJson)
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

func (self *JsonArray) Append(b []byte) []byte {

	l := len(self.elements)
	if l == 0 {
		return append(b, '[', ']')
	}

	b = append(b, '[')
	b = self.GetElement(0).Append(b)
	for i := 1; i < l; i++ {
		b = append(b, ',')
		b = self.GetElement(i).Append(b)
	}

	b = append(b, ']')
	return b

}

func (self *JsonArray) String() string {
	b := make([]byte, 0, 50)
	b = self.Append(b)
	return string(b)
}

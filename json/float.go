package json

import "fmt"

type JsonFloat float64

func (self JsonFloat) IsString() bool {
	return false
}

func (self JsonFloat) AsString() string {
	return ""
}

func (self JsonFloat) IsFloat() bool {
	return true
}

func (self JsonFloat) AsFloat() float64 {
	return float64(self)
}

func (self JsonFloat) IsInt() bool {
	return false
}

func (self JsonFloat) AsInt() int {
	return 0
}

func (self JsonFloat) IsBool() bool {
	return false
}

func (self JsonFloat) AsBool() bool {
	return false
}

func (self JsonFloat) IsObject() bool {
	return false
}

func (self JsonFloat) GetKeys() []string {
	return []string{}
}

func (self JsonFloat) GetMember(string) JsonNode {
	return JSON_NULL
}

func (self JsonFloat) IsArray() bool {
	return false
}

func (self JsonFloat) GetLen() int {
	return 0
}

func (self JsonFloat) GetElement(int) JsonNode {
	return JSON_NULL
}

func (self JsonFloat) IsNull() bool {
	return false
}

func (self JsonFloat) String() string {
	return fmt.Sprintf("%g", float64(self))
}

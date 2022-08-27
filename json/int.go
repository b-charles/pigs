package json

import "fmt"

type JsonInt int

func (self JsonInt) IsString() bool {
	return false
}

func (self JsonInt) AsString() string {
	return ""
}

func (self JsonInt) IsFloat() bool {
	return true
}

func (self JsonInt) AsFloat() float64 {
	return float64(self)
}

func (self JsonInt) IsInt() bool {
	return true
}

func (self JsonInt) AsInt() int {
	return int(self)
}

func (self JsonInt) IsBool() bool {
	return false
}

func (self JsonInt) AsBool() bool {
	return false
}

func (self JsonInt) IsObject() bool {
	return false
}

func (self JsonInt) GetKeys() []string {
	return []string{}
}

func (self JsonInt) GetMember(string) JsonNode {
	return JSON_NULL
}

func (self JsonInt) IsArray() bool {
	return false
}

func (self JsonInt) GetLen() int {
	return 0
}

func (self JsonInt) GetElement(int) JsonNode {
	return JSON_NULL
}

func (self JsonInt) IsNull() bool {
	return false
}

func (self JsonInt) String() string {
	return fmt.Sprintf("%d", int(self))
}

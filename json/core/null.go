package core

type JsonNull struct{}

func (self JsonNull) IsString() bool {
	return false
}

func (self JsonNull) AsString() string {
	return ""
}

func (self JsonNull) IsFloat() bool {
	return false
}

func (self JsonNull) AsFloat() float64 {
	return 0.0
}

func (self JsonNull) IsInt() bool {
	return false
}

func (self JsonNull) AsInt() int {
	return 0
}

func (self JsonNull) IsBool() bool {
	return false
}

func (self JsonNull) AsBool() bool {
	return false
}

func (self JsonNull) IsObject() bool {
	return false
}

func (self JsonNull) GetKeys() []string {
	return []string{}
}

func (self JsonNull) GetMember(string) JsonNode {
	return self
}

func (self JsonNull) IsArray() bool {
	return false
}

func (self JsonNull) GetLen() int {
	return 0
}

func (self JsonNull) GetElement(int) JsonNode {
	return self
}

func (self JsonNull) IsNull() bool {
	return true
}

func (self JsonNull) String() string {
	return "null"
}

var JSON_NULL = JsonNull{}

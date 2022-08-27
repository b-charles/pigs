package json

type JsonBool bool

func (self JsonBool) IsString() bool {
	return false
}

func (self JsonBool) AsString() string {
	return ""
}

func (self JsonBool) IsFloat() bool {
	return false
}

func (self JsonBool) AsFloat() float64 {
	return 0.0
}

func (self JsonBool) IsInt() bool {
	return false
}

func (self JsonBool) AsInt() int {
	return 0
}

func (self JsonBool) IsBool() bool {
	return true
}

func (self JsonBool) AsBool() bool {
	return bool(self)
}

func (self JsonBool) IsObject() bool {
	return false
}

func (self JsonBool) GetKeys() []string {
	return []string{}
}

func (self JsonBool) GetMember(string) JsonNode {
	return JSON_NULL
}

func (self JsonBool) IsArray() bool {
	return false
}

func (self JsonBool) GetLen() int {
	return 0
}

func (self JsonBool) GetElement(int) JsonNode {
	return JSON_NULL
}

func (self JsonBool) IsNull() bool {
	return false
}

func (self JsonBool) String() string {
	if bool(self) {
		return "true"
	} else {
		return "false"
	}
}

var JSON_TRUE = JsonBool(true)
var JSON_FALSE = JsonBool(false)

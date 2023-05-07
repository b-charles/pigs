package json

type JsonNode interface {
	IsString() bool
	AsString() string

	IsFloat() bool
	AsFloat() float64

	IsInt() bool
	AsInt() int

	IsBool() bool
	AsBool() bool

	IsObject() bool
	GetKeys() []string
	GetMember(string) JsonNode

	IsArray() bool
	GetLen() int
	GetElement(int) JsonNode

	IsNull() bool

	Append([]byte) []byte
	String() string
}

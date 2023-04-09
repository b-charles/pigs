package json

import (
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/b-charles/pigs/ioc"
)

var errorType = reflect.TypeOf(func(error) {}).In(0)
var stringType = reflect.TypeOf(func(string) {}).In(0)
var jsonType = reflect.TypeOf(func(JsonNode) {}).In(0)

type Jsons interface {
	Marshal(any) (JsonNode, error)
	MarshalToString(any) (string, error)
	Unmarshal(json JsonNode, callback any) error
	UnmarshalFromString(json string, callback any) error
}

type JsonMapper struct {
	marshallers     map[reflect.Type]*jsonValueMarshaller
	int_marshallers map[reflect.Type]*jsonValueMarshaller
	unmarshallers   map[reflect.Type]*jsonValueUnmarshaller
}

func typeComp(t1, t2 reflect.Type) int {
	if t1 == t2 {
		return 0
	} else if t1.Kind() == reflect.Pointer {
		if c := typeComp(t1.Elem(), t2); c == -1 {
			return -1
		} else {
			return 1
		}
	} else if t2.Kind() == reflect.Pointer {
		if c := typeComp(t1, t2.Elem()); c == -1 || c == 0 {
			return -1
		} else {
			return 1
		}
	} else if str1, str2 := t1.String(), t2.String(); str1 < str2 {
		return -1
	} else if str1 == str2 {
		return 0
	} else {
		return 1
	}
}

func sortedKeys[T any](m map[reflect.Type]T) []string {

	l := len(m)

	sorted := make([]reflect.Type, 0, l)
	for k := range m {
		sorted = append(sorted, k)
	}

	sort.Slice(sorted, func(i, j int) bool {
		return typeComp(sorted[i], sorted[j]) < 0
	})

	strs := make([]string, l)
	for i, t := range sorted {
		strs[i] = t.String()
	}

	return strs

}

func (self *JsonMapper) String() string {

	var b strings.Builder

	b.WriteString("[ marshallers:")
	for i, t := range sortedKeys(self.marshallers) {
		if i > 0 {
			b.WriteString(",")
		}
		b.WriteString(t)
	}

	b.WriteString(" unmarshallers:")
	for i, t := range sortedKeys(self.unmarshallers) {
		if i > 0 {
			b.WriteString(",")
		}
		b.WriteString(t)
	}

	b.WriteString(" ]")

	return b.String()

}

func (self *JsonMapper) PostInit(marshallers []JsonMarshaller, unmarshallers []JsonUnmarshaller) error {

	self.marshallers = map[reflect.Type]*jsonValueMarshaller{}
	self.int_marshallers = map[reflect.Type]*jsonValueMarshaller{}
	for _, marshaller := range marshallers {
		if target, value, err := valueMarshaller(marshaller); err != nil {
			return fmt.Errorf("Error during registring %v as a JsonMarshaller: %w", marshaller, err)
		} else {
			if target.Kind() == reflect.Interface {
				self.int_marshallers[target] = value
			} else {
				self.marshallers[target] = value
			}
		}
	}

	self.unmarshallers = map[reflect.Type]*jsonValueUnmarshaller{}
	for _, unmarshaller := range unmarshallers {
		if target, value, err := valueUnmarshaller(unmarshaller); err != nil {
			return fmt.Errorf("Error during registring %v as a JsonUnmarshaller: %w", unmarshaller, err)
		} else {
			self.unmarshallers[target] = value
		}
	}

	return nil

}

func (self *JsonMapper) getMarshaller(target reflect.Type) (*jsonValueMarshaller, error) {

	if marshaller, ok := self.marshallers[target]; ok {
		return marshaller, nil
	}

	for t, marshaller := range self.int_marshallers {
		if target.Implements(t) {
			self.marshallers[target] = marshaller
			return marshaller, nil
		}
	}

	if target.Kind() == reflect.Pointer {
		marshaller := &jsonValueMarshaller{}
		self.marshallers[target] = marshaller
		if err := newPointerMarshaller(self, target, marshaller); err != nil {
			delete(self.marshallers, target)
			return nil, err
		} else {
			return marshaller, nil
		}
	}

	if target.Kind() == reflect.Struct {
		marshaller := &jsonValueMarshaller{}
		self.marshallers[target] = marshaller
		if err := newStructMarshaller(self, target, marshaller); err != nil {
			delete(self.marshallers, target)
			return nil, err
		} else {
			return marshaller, nil
		}
	}

	if target.Kind() == reflect.Slice {
		marshaller := &jsonValueMarshaller{}
		self.marshallers[target] = marshaller
		if err := newSliceMarshaller(self, target, marshaller); err != nil {
			delete(self.marshallers, target)
			return nil, err
		} else {
			return marshaller, nil
		}
	}

	if target.Kind() == reflect.Map && target.Key() == stringType {
		marshaller := &jsonValueMarshaller{}
		self.marshallers[target] = marshaller
		if err := newMapMarshaller(self, target, marshaller); err != nil {
			delete(self.marshallers, target)
			return nil, err
		} else {
			return marshaller, nil
		}
	}

	return nil, fmt.Errorf("No Json marshaller found for type %v.", target)

}

func (self *JsonMapper) Marshal(value any) (JsonNode, error) {

	target := reflect.TypeOf(value)

	if marshaller, err := self.getMarshaller(target); err != nil {
		return JSON_NULL, err
	} else if json, err := marshaller.f(reflect.ValueOf(value)); err != nil {
		return JSON_NULL, err
	} else {
		return json, nil
	}

}

func (self *JsonMapper) MarshalToString(value any) (string, error) {
	if json, err := self.Marshal(value); err != nil {
		return "", nil
	} else {
		return json.String(), nil
	}
}

func (self *JsonMapper) getUnmarshaller(target reflect.Type) (*jsonValueUnmarshaller, error) {

	if unmarshaller, ok := self.unmarshallers[target]; ok {
		return unmarshaller, nil
	}

	if target.Kind() == reflect.Pointer {
		unmarshaller := &jsonValueUnmarshaller{}
		self.unmarshallers[target] = unmarshaller
		if err := newPointerUnmarshaller(self, target, unmarshaller); err != nil {
			delete(self.unmarshallers, target)
			return nil, err
		} else {
			return unmarshaller, nil
		}
	}

	if target.Kind() == reflect.Struct {
		unmarshaller := &jsonValueUnmarshaller{}
		self.unmarshallers[target] = unmarshaller
		if err := newStructUnmarshaller(self, target, unmarshaller); err != nil {
			delete(self.unmarshallers, target)
			return nil, err
		} else {
			return unmarshaller, nil
		}
	}

	if target.Kind() == reflect.Slice {
		unmarshaller := &jsonValueUnmarshaller{}
		self.unmarshallers[target] = unmarshaller
		if err := newSliceUnmarshaller(self, target, unmarshaller); err != nil {
			delete(self.unmarshallers, target)
			return nil, err
		} else {
			return unmarshaller, nil
		}
	}

	if target.Kind() == reflect.Map && target.Key() == stringType {
		unmarshaller := &jsonValueUnmarshaller{}
		self.unmarshallers[target] = unmarshaller
		if err := newMapUnMarshaller(self, target, unmarshaller); err != nil {
			delete(self.unmarshallers, target)
			return nil, err
		} else {
			return unmarshaller, nil
		}
	}

	return nil, fmt.Errorf("No Json unmarshaller found for type %v.", target)

}

func (self *JsonMapper) Unmarshal(json JsonNode, callback any) error {

	t := reflect.TypeOf(callback)
	if t.Kind() != reflect.Func {
		return fmt.Errorf("Invalid argument: the callback should be a function.")
	} else if t.NumIn() != 1 {
		return fmt.Errorf("Invalid argument: the callback funcion should take one input.")
	} else if t.NumOut() != 0 {
		return fmt.Errorf("Invalid argument: the callback funcion should take no output.")
	}

	if unmarshaller, err := self.getUnmarshaller(t.In(0)); err != nil {
		return err
	} else if value, err := unmarshaller.f(json); err != nil {
		return err
	} else {
		reflect.ValueOf(callback).Call([]reflect.Value{value})
		return nil
	}

}

func (self *JsonMapper) UnmarshalFromString(str string, callback any) error {
	if json, err := ParseString(str); err != nil {
		return err
	} else {
		return self.Unmarshal(json, callback)
	}
}

func init() {
	ioc.PutFactory(func() (*JsonMapper, error) {
		return &JsonMapper{}, nil
	}, func(Jsons) {})
}

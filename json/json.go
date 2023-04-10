package json

import (
	"fmt"
	"reflect"
	"sort"

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
	marshallers     map[reflect.Type]*wrappedMarshaller
	int_marshallers []*wrappedMarshaller
	unmarshallers   map[reflect.Type]*wrappedUnmarshaller
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

func sortByKeys[T any](m map[reflect.Type]T) []T {

	l := len(m)

	types := make([]reflect.Type, 0, l)
	for k := range m {
		types = append(types, k)
	}

	sort.Slice(types, func(i, j int) bool {
		return typeComp(types[i], types[j]) < 0
	})

	sorted := make([]T, 0, l)
	for _, t := range types {
		sorted = append(sorted, m[t])
	}

	return sorted

}

func (self *JsonMapper) JsonNode() JsonNode {

	b := NewJsonBuilder()

	b.SetEmptyObject("marshallers")
	for _, m := range sortByKeys(self.marshallers) {
		b.SetString(fmt.Sprintf("marshallers.%s", m.t.String()), m.String())
	}

	b.SetEmptyObject("interface_marshallers")
	for _, m := range self.int_marshallers {
		b.SetString(fmt.Sprintf("interface_marshallers.%s", m.t.String()), m.String())
	}

	b.SetEmptyObject("unmarshallers")
	for _, m := range sortByKeys(self.unmarshallers) {
		b.SetString(fmt.Sprintf("unmarshallers.%s", m.t.String()), m.String())
	}

	return b.Build()

}

func (self *JsonMapper) String() string {
	return self.JsonNode().String()
}

func (self *JsonMapper) insertIntWrappedMarshaller(wrapped *wrappedMarshaller) error {

	target := wrapped.t

	for i, m := range self.int_marshallers {

		if target == m.t {
			return fmt.Errorf("Marshaller %v already defined for type %v.", m, m.t)
		} else if target.Implements(m.t) {

			if m.t.Implements(target) {
				return fmt.Errorf("Interfaces %v and %v (marshalled by %v) are identical.", target, m, m.t)
			}

			self.int_marshallers = append(self.int_marshallers[:i+1], self.int_marshallers[i:]...)
			self.int_marshallers[i] = wrapped
			return nil

		}
	}

	self.int_marshallers = append(self.int_marshallers, wrapped)
	return nil

}

func (self *JsonMapper) insertMarshaller(marshaller JsonMarshaller) error {

	if wrapped, err := wrapMarshaller(marshaller); err != nil {
		return err
	} else {

		target := wrapped.t
		if target.Kind() == reflect.Interface {

			return self.insertIntWrappedMarshaller(wrapped)

		} else {

			if m, ok := self.marshallers[target]; ok {
				return fmt.Errorf("Marshaller %v already defined for type %v.", m, target)
			}
			self.marshallers[target] = wrapped

		}

	}

	return nil

}

func (self *JsonMapper) insertUnmarshaller(unmarshaller JsonUnmarshaller) error {

	if wrapped, err := wrapUnmarshaller(unmarshaller); err != nil {
		return err
	} else {

		target := wrapped.t
		if m, ok := self.unmarshallers[target]; ok {
			return fmt.Errorf("Unmarshaller %v already defined for type %v.", m, target)
		}
		self.unmarshallers[target] = wrapped

	}

	return nil

}

func (self *JsonMapper) PostInit(marshallers []JsonMarshaller, unmarshallers []JsonUnmarshaller) error {

	self.marshallers = map[reflect.Type]*wrappedMarshaller{}
	self.int_marshallers = make([]*wrappedMarshaller, 0)
	for _, marshaller := range marshallers {
		if err := self.insertMarshaller(marshaller); err != nil {
			return fmt.Errorf("Error during registring %v as a JsonMarshaller: %w", marshaller, err)
		}
	}

	self.unmarshallers = map[reflect.Type]*wrappedUnmarshaller{}
	for _, unmarshaller := range unmarshallers {
		if err := self.insertUnmarshaller(unmarshaller); err != nil {
			return fmt.Errorf("Error during registring %v as a JsonUnmarshaller: %w", unmarshaller, err)
		}
	}

	return nil

}

func (self *JsonMapper) getMarshaller(target reflect.Type) (*wrappedMarshaller, error) {

	if marshaller, ok := self.marshallers[target]; ok {
		return marshaller, nil
	}

	for _, marshaller := range self.int_marshallers {
		if target.Implements(marshaller.t) {
			self.marshallers[target] = marshaller
			return marshaller, nil
		}
	}

	if target.Kind() == reflect.Pointer {
		marshaller := &wrappedMarshaller{}
		self.marshallers[target] = marshaller
		if err := newPointerMarshaller(self, target, marshaller); err != nil {
			delete(self.marshallers, target)
			return nil, err
		} else {
			return marshaller, nil
		}
	}

	if target.Kind() == reflect.Struct {
		marshaller := &wrappedMarshaller{}
		self.marshallers[target] = marshaller
		if err := newStructMarshaller(self, target, marshaller); err != nil {
			delete(self.marshallers, target)
			return nil, err
		} else {
			return marshaller, nil
		}
	}

	if target.Kind() == reflect.Slice {
		marshaller := &wrappedMarshaller{}
		self.marshallers[target] = marshaller
		if err := newSliceMarshaller(self, target, marshaller); err != nil {
			delete(self.marshallers, target)
			return nil, err
		} else {
			return marshaller, nil
		}
	}

	if target.Kind() == reflect.Map && target.Key() == stringType {
		marshaller := &wrappedMarshaller{}
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

func (self *JsonMapper) getUnmarshaller(target reflect.Type) (*wrappedUnmarshaller, error) {

	if unmarshaller, ok := self.unmarshallers[target]; ok {
		return unmarshaller, nil
	}

	if target.Kind() == reflect.Pointer {
		unmarshaller := &wrappedUnmarshaller{}
		self.unmarshallers[target] = unmarshaller
		if err := newPointerUnmarshaller(self, target, unmarshaller); err != nil {
			delete(self.unmarshallers, target)
			return nil, err
		} else {
			return unmarshaller, nil
		}
	}

	if target.Kind() == reflect.Struct {
		unmarshaller := &wrappedUnmarshaller{}
		self.unmarshallers[target] = unmarshaller
		if err := newStructUnmarshaller(self, target, unmarshaller); err != nil {
			delete(self.unmarshallers, target)
			return nil, err
		} else {
			return unmarshaller, nil
		}
	}

	if target.Kind() == reflect.Slice {
		unmarshaller := &wrappedUnmarshaller{}
		self.unmarshallers[target] = unmarshaller
		if err := newSliceUnmarshaller(self, target, unmarshaller); err != nil {
			delete(self.unmarshallers, target)
			return nil, err
		} else {
			return unmarshaller, nil
		}
	}

	if target.Kind() == reflect.Map && target.Key() == stringType {
		unmarshaller := &wrappedUnmarshaller{}
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

	ioc.Put(func(jsonMapper *JsonMapper) (JsonNode, error) {
		return jsonMapper.JsonNode(), nil
	}, func(JsonMarshaller) {})

}

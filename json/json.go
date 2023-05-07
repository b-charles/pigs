package json

import (
	"fmt"
	"reflect"

	"github.com/b-charles/pigs/ioc"
)

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

func (self *JsonMapper) Json() JsonNode {

	b := NewJsonBuilder()

	m := make([]reflect.Type, 0, len(self.marshallers))
	for k := range self.marshallers {
		m = append(m, k)
	}
	b.Set("marshallers", ReflectTypeSliceToJson(m))

	i := make([]reflect.Type, 0, len(self.int_marshallers))
	for _, k := range self.int_marshallers {
		i = append(i, k.t)
	}
	b.Set("interface_marshallers", ReflectTypeSliceToJson(i))

	u := make([]reflect.Type, 0, len(self.unmarshallers))
	for k := range self.unmarshallers {
		u = append(u, k)
	}
	b.Set("unmarshallers", ReflectTypeSliceToJson(u))

	return b.Build()

}

func (self *JsonMapper) String() string {
	return self.Json().String()
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

	if target.Implements(jsonerType) {
		marshaller := &wrappedMarshaller{}
		self.marshallers[target] = marshaller
		if err := newJsonerMarshaller(self, target, marshaller); err != nil {
			delete(self.marshallers, target)
			return nil, err
		} else {
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

	ioc.PutNamedFactory("Json mapper", func() (*JsonMapper, error) {
		return &JsonMapper{}, nil
	}, func(Jsons) {})

}

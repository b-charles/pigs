package json

import (
	"fmt"
	"reflect"

	"github.com/b-charles/pigs/ioc"
	"github.com/b-charles/pigs/memfun"
)

type Jsons interface {
	Marshal(any) (JsonNode, error)
	MarshalToString(any) (string, error)
	Unmarshal(json JsonNode, callback any) error
	UnmarshalFromString(json string, callback any) error
}

func append_int_mashallers(marshallers []*wrappedMarshaller, wrapped *wrappedMarshaller) ([]*wrappedMarshaller, error) {

	target := wrapped.t

	for i, m := range marshallers {

		if target == m.t {
			return marshallers, fmt.Errorf("Marshaller %v already defined for type %v.", m, m.t)
		} else if target.Implements(m.t) {

			if m.t.Implements(target) {
				return marshallers, fmt.Errorf("Interfaces %v and %v (marshalled by %v) are identical.", target, m, m.t)
			}

			marshallers = append(marshallers[:i+1], marshallers[i:]...)
			marshallers[i] = wrapped
			return marshallers, nil

		}

	}

	return append(marshallers, wrapped), nil

}

type JsonMapper struct {
	marshallers     memfun.MemFun[reflect.Type, *wrappedMarshaller]
	int_marshallers []*wrappedMarshaller
	unmarshallers   memfun.MemFun[reflect.Type, *wrappedUnmarshaller]
}

func (self *JsonMapper) Json() JsonNode {

	b := NewJsonBuilder()

	m := self.marshallers.Keys()
	b.Set("marshallers", ReflectTypeSliceToJson(m))

	i := make([]reflect.Type, 0, len(self.int_marshallers))
	for _, k := range self.int_marshallers {
		i = append(i, k.t)
	}
	b.Set("interface_marshallers", ReflectTypeSliceToJson(i))

	u := self.unmarshallers.Keys()
	b.Set("unmarshallers", ReflectTypeSliceToJson(u))

	return b.Build()

}

func (self *JsonMapper) String() string {
	return self.Json().String()
}

func (self *JsonMapper) PostInit(marshallers []JsonMarshaller, unmarshallers []JsonUnmarshaller) error {

	self.int_marshallers = make([]*wrappedMarshaller, 0)
	self.marshallers = memfun.NewMemFun(
		func(target reflect.Type, rec func(reflect.Type) (*wrappedMarshaller, error)) (*wrappedMarshaller, error) {

			for _, marshaller := range self.int_marshallers {
				if target.Implements(marshaller.t) {
					return marshaller, nil
				}
			}

			if target.Implements(jsonerType) {
				return newJsonerMarshaller(target)
			}

			if target.Kind() == reflect.Pointer {
				return newPointerMarshaller(target, rec)
			}

			if target.Kind() == reflect.Struct {
				return newStructMarshaller(target, rec)
			}

			if target.Kind() == reflect.Slice {
				return newSliceMarshaller(target, rec)
			}

			if target.Kind() == reflect.Map && target.Key() == stringType {
				return newMapMarshaller(target, rec)
			}

			return nil, fmt.Errorf("No Json marshaller found for type %v.", target)

		})

	for _, marshaller := range marshallers {

		if wrapped, err := wrapMarshaller(marshaller); err != nil {
			return err
		} else {

			target := wrapped.t
			if target.Kind() == reflect.Interface {

				if self.int_marshallers, err = append_int_mashallers(self.int_marshallers, wrapped); err != nil {
					return err
				}

			} else {

				if m, ok, err := self.marshallers.Lookup(target); err != nil {
					return err
				} else if ok {
					return fmt.Errorf("Marshaller %v already defined for type %v.", m, target)
				} else {
					self.marshallers.Store(target, wrapped)
				}

			}

		}

	}

	self.unmarshallers = memfun.NewMemFun(
		func(target reflect.Type, rec func(reflect.Type) (*wrappedUnmarshaller, error)) (*wrappedUnmarshaller, error) {

			if target.Kind() == reflect.Pointer {
				return newPointerUnmarshaller(target, rec)
			}

			if target.Kind() == reflect.Struct {
				return newStructUnmarshaller(target, rec)
			}

			if target.Kind() == reflect.Slice {
				return newSliceUnmarshaller(target, rec)
			}

			if target.Kind() == reflect.Map && target.Key() == stringType {
				return newMapUnMarshaller(target, rec)
			}

			return nil, fmt.Errorf("No Json unmarshaller found for type %v.", target)

		})

	for _, unmarshaller := range unmarshallers {

		if wrapped, err := wrapUnmarshaller(unmarshaller); err != nil {
			return err
		} else {

			target := wrapped.t
			if m, ok, err := self.unmarshallers.Lookup(target); err != nil {
				return err
			} else if ok {
				return fmt.Errorf("Unmarshaller %v already defined for type %v.", m, target)
			} else {
				self.unmarshallers.Store(target, wrapped)
			}

		}

	}

	return nil

}

func (self *JsonMapper) Marshal(value any) (JsonNode, error) {

	target := reflect.TypeOf(value)

	if marshaller, err := self.marshallers.Get(target); err != nil {
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

func (self *JsonMapper) Unmarshal(json JsonNode, callback any) error {

	t := reflect.TypeOf(callback)
	if t.Kind() != reflect.Func {
		return fmt.Errorf("Invalid argument: the callback should be a function.")
	} else if t.NumIn() != 1 {
		return fmt.Errorf("Invalid argument: the callback funcion should take one input.")
	} else if t.NumOut() != 0 {
		return fmt.Errorf("Invalid argument: the callback funcion should take no output.")
	}

	if unmarshaller, err := self.unmarshallers.Get(t.In(0)); err != nil {
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

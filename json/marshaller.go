package json

import (
	"fmt"
	"reflect"
)

// type JsonMarshaller func(T) (JsonNode, error)
type JsonMarshaller any

type wrappedMarshaller struct {
	src JsonMarshaller
	t   reflect.Type
	f   func(reflect.Value) (JsonNode, error)
}

func (self *wrappedMarshaller) String() string {
	return fmt.Sprint(self.src)
}

func wrapMarshaller(marshaller JsonMarshaller) (*wrappedMarshaller, error) {

	v := reflect.ValueOf(marshaller)
	t := v.Type()

	if t.Kind() != reflect.Func {
		return nil, fmt.Errorf("Invalid marshaller: %v is not a function.", v)
	} else if t.NumIn() != 1 {
		return nil, fmt.Errorf("Invalid marshaller (wrong inputs): %v, expected func(T) (JsonNode, error)", t)
	} else if t.NumOut() != 2 || t.Out(0) != jsonType || !t.Out(1).AssignableTo(errorType) {
		return nil, fmt.Errorf("Invalid marshaller (wrong outputs): %v, expected func(T) (JsonNode, error)", t)
	}

	f := func(value reflect.Value) (JsonNode, error) {
		outs := v.Call([]reflect.Value{value})
		if err := outs[1].Interface(); err != nil {
			return outs[0].Interface().(JsonNode), err.(error)
		} else {
			return outs[0].Interface().(JsonNode), nil
		}
	}

	return &wrappedMarshaller{marshaller, t.In(0), f}, nil

}

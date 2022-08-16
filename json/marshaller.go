package json

import (
	"fmt"
	"reflect"

	"github.com/b-charles/pigs/json/core"
)

// type JsonMarshaller func(T) (core.JsonNode, error)
type JsonMarshaller any

type jsonValueMarshaller struct {
	f func(reflect.Value) (core.JsonNode, error)
}

func valueMarshaller(marshaller JsonMarshaller) (reflect.Type, *jsonValueMarshaller, error) {

	v := reflect.ValueOf(marshaller)
	t := v.Type()

	if t.Kind() != reflect.Func {
		return nil, nil, fmt.Errorf("Invalid marshaller: %v is not a function.", v)
	} else if t.NumIn() != 1 {
		return nil, nil, fmt.Errorf("Invalid marshaller (wrong inputs): %v, expected func(T) (core.JsonNode, error)", t)
	} else if t.NumOut() != 2 || t.Out(0) != jsonType || !t.Out(1).AssignableTo(errorType) {
		return nil, nil, fmt.Errorf("Invalid marshaller (wrong outputs): %v, expected func(T) (core.JsonNode, error)", t)
	}

	f := func(value reflect.Value) (core.JsonNode, error) {
		outs := v.Call([]reflect.Value{value})
		if err := outs[1].Interface(); err != nil {
			return outs[0].Interface().(core.JsonNode), err.(error)
		} else {
			return outs[0].Interface().(core.JsonNode), nil
		}
	}

	return t.In(0), &jsonValueMarshaller{f}, nil

}

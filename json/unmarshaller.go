package json

import (
	"fmt"
	"reflect"

	"github.com/b-charles/pigs/json/core"
)

// type JsonUnmarshaller func(core.JsonNode) (T, error)
type JsonUnmarshaller any

type jsonValueUnmarshaller struct {
	f func(core.JsonNode) (reflect.Value, error)
}

func valueUnmarshaller(unmarshaller JsonUnmarshaller) (reflect.Type, *jsonValueUnmarshaller, error) {

	v := reflect.ValueOf(unmarshaller)
	t := v.Type()

	if t.Kind() != reflect.Func {
		return nil, nil, fmt.Errorf("Invalid unmarshaller: %v is not a function.", v)
	} else if t.NumIn() != 1 || t.In(0) != jsonType {
		return nil, nil, fmt.Errorf("Invalid unmarshaller (wrong inputs): %v, expected func(core.JsonNode) (T, error)", t)
	} else if t.NumOut() != 2 || !t.Out(1).AssignableTo(errorType) {
		return nil, nil, fmt.Errorf("Invalid unmarshaller (wrong outputs): %v, expected func(core.JsonNode) (T, error)", t)
	}

	f := func(json core.JsonNode) (reflect.Value, error) {
		outs := v.Call([]reflect.Value{reflect.ValueOf(json)})
		if err := outs[1].Interface(); err != nil {
			return reflect.Value{}, err.(error)
		} else {
			return outs[0], nil
		}
	}

	return t.Out(0), &jsonValueUnmarshaller{f}, nil

}

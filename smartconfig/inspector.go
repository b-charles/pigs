package smartconfig

import (
	"fmt"
	"reflect"
)

// type Inspector func(NavConfig) (any, error)
type Inspector any

func inspectorConfigurer(inspector Inspector) (*configurer, error) {

	value := reflect.ValueOf(inspector)

	typ := value.Type()
	if typ.Kind() != reflect.Func {
		return nil, fmt.Errorf("An inspector should be a function.")
	}
	if numIn := typ.NumIn(); numIn != 1 {
		return nil, fmt.Errorf("An inspector should have only one input, not %d.", numIn)
	}
	if in := typ.In(0); in != navconfig_type {
		return nil, fmt.Errorf("An inspector's only input should be a NavConfig, not a %v.", in)
	}
	if numOut := typ.NumOut(); numOut != 2 {
		return nil, fmt.Errorf("An inspector should have two outputs, not %d.", numOut)
	}
	if out2 := typ.Out(1); !out2.AssignableTo(error_type) {
		return nil, fmt.Errorf("Then inspector's second output should be an error, not a %v.", out2)
	}

	return &configurer{
		target: typ.Out(0),
		setter: func(config NavConfig, receiver reflect.Value) error {

			outs := value.Call([]reflect.Value{reflect.ValueOf(config)})
			if !outs[1].IsNil() {
				return outs[1].Interface().(error)
			}

			receiver.Set(outs[0])
			return nil

		}}, nil

}

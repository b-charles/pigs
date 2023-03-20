package smartconfig

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/b-charles/pigs/ioc"
)

// type Parser func(string) (any, error)
type Parser any

func init() {

	ioc.Put(func(value string) (string, error) {
		return value, nil
	}, func(Parser) {})
	ioc.Put(strconv.ParseBool, func(Parser) {})
	ioc.Put(strconv.Atoi, func(Parser) {})

}

func parserConfigurer(parser Parser) (*configurer, error) {

	value := reflect.ValueOf(parser)

	typ := value.Type()
	if typ.Kind() != reflect.Func {
		return nil, fmt.Errorf("A parser should be a function.")
	}
	if numIn := typ.NumIn(); numIn != 1 {
		return nil, fmt.Errorf("A parser should have only one input, not %d.", numIn)
	}
	if in := typ.In(0); in != string_type {
		return nil, fmt.Errorf("A parser's only input should be a string, not a %v.", in)
	}
	if numOut := typ.NumOut(); numOut != 2 {
		return nil, fmt.Errorf("A parser should have two outputs, not %d.", numOut)
	}
	if out2 := typ.Out(1); !out2.AssignableTo(error_type) {
		return nil, fmt.Errorf("The parser's second output should be an error, not a %v.", out2)
	}

	return &configurer{
		source: parser,
		target: typ.Out(0),
		setter: func(config NavConfig, receiver reflect.Value) error {

			outs := value.Call([]reflect.Value{reflect.ValueOf(config.Value())})
			if !outs[1].IsNil() {
				return outs[1].Interface().(error)
			}

			receiver.Set(outs[0])
			return nil

		}}, nil

}

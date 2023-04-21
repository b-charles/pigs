package ioc

import (
	"fmt"
	"reflect"
)

// A component is a singleton recording, managed by a Container. At most one of
// 'value' or 'factory' field can be valid. Fields 'name' and 'main' are only
// useful for debugging (see container status).
type component struct {
	name       string
	main       reflect.Type
	value      reflect.Value
	factory    reflect.Value
	signatures []reflect.Type
}

// checkFactory checks if the input is an acceptable factory.
func checkFactory(factory any) (reflect.Type, reflect.Value, error) {

	val := reflect.ValueOf(factory)
	typ := val.Type()

	if typ.Kind() != reflect.Func {
		return nil, val, fmt.Errorf("The factory should be a function, not a %v.", typ)
	}

	if nout := typ.NumOut(); nout == 0 {
		return nil, val, fmt.Errorf("The function should return at least one value.")
	} else if nout == 2 && !typ.Out(1).AssignableTo(error_type) {
		return nil, val, fmt.Errorf("The second output of the factory (%v) should be assignable to error.", typ.Out(1))
	} else if nout > 2 {
		return nil, val, fmt.Errorf("The factory should return one or two outputs, not %d.", nout)
	}

	return typ.Out(0), val, nil

}

// extractSignatures extracts all signatures from a slice of signature
// function.
func extractSignatures(main reflect.Type, signFuncs []any) ([]reflect.Type, error) {

	unique := make(map[reflect.Type]bool)
	unique[main] = true

	for _, f := range signFuncs {

		typ := reflect.TypeOf(f)
		if typ.Kind() != reflect.Func {
			return nil, fmt.Errorf("Signatures should be defined in a function, not a %v.", typ)
		} else if nout := typ.NumOut(); nout > 0 {
			return nil, fmt.Errorf("A signatures function should return no output, not %d.", nout)
		}

		for i := 0; i < typ.NumIn(); i++ {
			sig := typ.In(i)
			if _, ok := unique[sig]; !ok {
				if !main.AssignableTo(sig) {
					return nil, fmt.Errorf("The main type '%v' doesn't implement the signature type '%v'.", main, sig)
				}
				unique[sig] = true
			}
		}

	}

	signatures := make([]reflect.Type, 0, len(unique))
	for sig := range unique {
		signatures = append(signatures, sig)
	}

	return signatures, nil

}

// newComponent returns a new component defined by its name, its value or its
// factory and the signatures functions. At least one of the arguments 'value'
// or 'factory' should be nil.
func newComponent(name string, value any, factory any, signFuncs []any) (*component, error) {

	if value != nil && factory != nil {
		return nil, fmt.Errorf("The component '%v' can not be defined with a value and a factory.", name)
	}

	var (
		main     reflect.Type
		factory_ reflect.Value
		value_   reflect.Value
	)

	if factory != nil {

		var err error
		if main, factory_, err = checkFactory(factory); err != nil {
			return nil, fmt.Errorf("Error during the registration of '%v': %w", name, err)
		}

	} else if value != nil {

		value_ = reflect.ValueOf(value)
		main = value_.Type()

	} else {

		return &component{name, main, value_, factory_, []reflect.Type{}}, nil

	}

	// signatures
	signatures, err := extractSignatures(main, signFuncs)
	if err != nil {
		return nil, fmt.Errorf("Error during the registration of '%v': %w", name, err)
	}

	// return
	return &component{name, main, value_, factory_, signatures}, nil

}

// String returns the name of the component.
func (self *component) String() string {
	return self.name
}

package ioc

import (
	"fmt"
	"reflect"
)

var error_type = reflect.TypeOf(func(error) {}).In(0)

// A component is a singleton recording, managed by a Container.
type component struct {
	container  *Container
	main       reflect.Type
	signatures []reflect.Type
	factory    reflect.Value
}

// addErrorOutput take typ and val, the type and the value of a function, and
// returns a function with the same inputs but with an error (nil) added in the
// outputs.
func addErrorOutput(typ reflect.Type, val reflect.Value) reflect.Value {

	in := make([]reflect.Type, 0, typ.NumIn())
	for i := 0; i < typ.NumIn(); i++ {
		in = append(in, typ.In(i))
	}

	out := []reflect.Type{typ.Out(0), error_type}

	wTyp := reflect.FuncOf(in, out, typ.IsVariadic())
	wVal := reflect.MakeFunc(wTyp, func(args []reflect.Value) []reflect.Value {
		return append(val.Call(args), reflect.Zero(error_type))
	})

	return wVal

}

// checkFactory checks if the input is an acceptable factory.
func checkFactory(factory reflect.Value) (reflect.Type, reflect.Value, error) {

	typ := factory.Type()

	if typ.Kind() != reflect.Func {
		return nil, factory, fmt.Errorf("The type %v should be a function.", typ)
	}

	if nout := typ.NumOut(); nout == 0 {
		return nil, factory, fmt.Errorf("The function should return at least one value.")
	} else if nout == 1 {
		factory = addErrorOutput(typ, factory)
	} else if nout == 2 && !typ.Out(1).AssignableTo(error_type) {
		return nil, factory, fmt.Errorf("The second output of the factory (%v) should be assignable to %v.", typ.Out(1), error_type)
	} else if nout > 2 {
		return nil, factory, fmt.Errorf("The factory should return one or two outputs, not %d.", nout)
	}

	return typ.Out(0), factory, nil

}

// extractSignature extracts signatures from one signature function.
func extractSignature(signFunc any) ([]reflect.Type, error) {

	typ := reflect.TypeOf(signFunc)
	if typ.Kind() != reflect.Func {
		return nil, fmt.Errorf("A signature wrapper should be a function, not a '%v'.", typ)
	} else if nout := typ.NumOut(); nout > 0 {
		return nil, fmt.Errorf("A signature wrapper should return no output, not %d.", nout)
	}

	list := make([]reflect.Type, 0, typ.NumIn())
	for i := 0; i < typ.NumIn(); i++ {
		list = append(list, typ.In(i))
	}

	return list, nil

}

// extractSignatures extracts all signatures from a slice of signature
// function.
func extractSignatures(signFuncs []any) ([]reflect.Type, error) {

	list := []reflect.Type{}

	for _, f := range signFuncs {

		if i, err := extractSignature(f); err != nil {
			return nil, err
		} else {
			list = append(list, i...)
		}

	}

	return list, nil

}

// newComponent returns a new component, build from the containing container,
// the factory and the signature functions.
func newComponent(container *Container, factory reflect.Value, signFuncs []any) (*component, error) {

	// factory
	main, factoryValue, err := checkFactory(factory)
	if err != nil {
		return nil, err
	}

	// signatures

	signs, err := extractSignatures(signFuncs)
	if err != nil {
		return nil, fmt.Errorf("Error during registration of '%v': %w", main, err)
	}

	uniqueMap := make(map[reflect.Type]bool)
	uniqueMap[main] = true
	for _, sig := range signs {
		if !main.AssignableTo(sig) {
			return nil, fmt.Errorf("The component '%v' is defined with the signature '%v' but doesn't implements it.", main, sig)
		}
		uniqueMap[sig] = true
	}

	signatures := make([]reflect.Type, 0, len(uniqueMap))
	for sig := range uniqueMap {
		signatures = append(signatures, sig)
	}

	// return
	return &component{container, main, signatures, factoryValue}, nil

}

// instanciate returns an instance of the component.
func (self *component) instanciate(stack *componentStack) (*instance, error) {

	if self == nil {
		return nil, nil
	}

	if err := stack.push(self); err != nil {
		return nil, err
	}

	outs, err := self.container.callInjected(self.factory, stack)

	stack.pop(self)

	if err != nil {
		return nil, fmt.Errorf("Error during call of factory of '%v': %w", self, err)
	} else if !outs[1].IsNil() {
		return nil, fmt.Errorf("Error during instanciation of '%v': %w", self, outs[1].Interface().(error))
	}

	comp := outs[0]
	if comp.Kind() == reflect.Interface {
		comp = comp.Elem()
	}

	return &instance{self, comp}, nil

}

// String returns a string representation of the component.
func (self *component) String() string {
	if self == nil {
		return "nil"
	} else {
		return self.main.String()
	}
}

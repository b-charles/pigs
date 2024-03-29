package ioc

import (
	"fmt"
	"io"
	"reflect"
)

// instance represents a component instance.
type instance struct {
	component *component
	value     reflect.Value
}

// newInstance returns an instance of the component.
func newInstance(container *Container, component *component, stack *componentStack) (*instance, error) {

	if err := stack.push(component); err != nil {
		return nil, err
	}
	defer stack.pop(component)

	if component.factory.IsValid() {

		outs, err := container.callInjected(component.factory, stack)

		if err != nil {
			return nil, fmt.Errorf("Error during call of factory of '%v': %w", component, err)
		} else if len(outs) == 2 && !outs[1].IsNil() {
			return nil, fmt.Errorf("Error during instanciation of '%v': %w", component, outs[1].Interface().(error))
		}

		value := outs[0]
		if value.IsValid() && value.Kind() == reflect.Interface {
			value = value.Elem()
		}

		return &instance{component, value}, nil

	} else if component.value.IsValid() {

		return &instance{component, component.value}, nil

	} else {

		return &instance{component, component.value}, nil

	}

}

func (self *instance) isNil() bool {

	if !self.value.IsValid() {
		return true
	}

	kind := self.value.Kind()
	if kind == reflect.Chan ||
		kind == reflect.Func ||
		kind == reflect.Map ||
		kind == reflect.Pointer ||
		kind == reflect.UnsafePointer ||
		kind == reflect.Interface ||
		kind == reflect.Slice {

		return self.value.IsNil()

	}

	return false

}

// initialize initializes the instance: if the instance is a struct or a
// pointer to a struct, each tagged 'inject' field is injected.
func (self *instance) initialize(container *Container, stack *componentStack) error {

	if self.isNil() {
		return nil
	}

	value := self.value
	if value.Kind() == reflect.Pointer {
		value = value.Elem()
	}
	if value.Kind() != reflect.Struct {
		return nil
	}

	typ := value.Type()

	for i := 0; i < value.NumField(); i++ {

		field := value.Field(i)
		structField := typ.Field(i)

		if _, ok := structField.Tag.Lookup("inject"); !ok {
			continue
		} else if !field.CanSet() {
			return fmt.Errorf("The field '%v' of %v is not settable.", structField.Name, self)
		} else if fieldValue, err := container.getValue(structField.Type, stack); err != nil {
			return fmt.Errorf("Can not inject field '%v': %w", structField.Name, err)
		} else {
			field.Set(fieldValue)
		}

	}

	return nil

}

// postInit call the PostInit method (if defined).
func (self *instance) postInit(container *Container, stack *componentStack) error {

	if self.isNil() {
		return nil
	}

	postInit := self.value.MethodByName("PostInit")
	if !postInit.IsValid() {
		return nil
	}

	if out, err := container.callInjected(postInit, stack); err != nil {

		return fmt.Errorf("Error during PostInit call of '%v': %w", self, err)

	} else if len(out) > 1 {

		return fmt.Errorf("The PostInit method of '%v' should return none or one output, not %d.", self, len(out))

	} else if len(out) == 1 {

		if outType := postInit.Type().Out(0); !outType.AssignableTo(error_type) {
			return fmt.Errorf("The output of the PostInit method of '%v' should be an error, not a '%v'.", self, outType)
		} else if err := out[0].Interface(); err != nil {
			return fmt.Errorf("Error returned by PostInit of '%v': %w", self, err.(error))
		} else {
			return nil
		}

	} else {

		return nil

	}

}

// castToCloser checks if the instance can be casted to an io.Closer.
func (self *instance) castToCloser() (io.Closer, bool) {
	if self.isNil() {
		return nil, false
	} else {
		closer, closable := self.value.Interface().(io.Closer)
		return closer, closable
	}
}

// isClosable returns true if the instance implements io.Closer.
func (self *instance) isClosable() bool {
	_, closable := self.castToCloser()
	return closable
}

// Close calls the Close method (if defined).
func (self *instance) Close() {
	if closer, closable := self.castToCloser(); closable {
		defer func() { recover() }()
		closer.Close()
	}
}

// String returns a string representation of the instance.
func (self *instance) String() string {
	if self.isNil() {
		return "nil"
	} else {
		return fmt.Sprintf("%v", self.value.Interface())
	}
}

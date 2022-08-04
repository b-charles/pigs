package ioc

import (
	"fmt"
	"io"
	"reflect"
)

// instance represents a component instance. It's created by the factory of the
// component.
type instance struct {
	producer *component
	value    reflect.Value
}

// initialize initializes the instance: if the instance is a struct or a
// pointer to a struct, each tagged 'inject' field is injected.
func (self *instance) initialize(stack *componentStack) error {

	if self == nil {
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
		} else if fieldValue, err := self.producer.container.getValue(structField.Type, stack); err != nil {
			return fmt.Errorf("Can not inject field '%v' of %v: %w", structField.Name, self, err)
		} else {
			field.Set(fieldValue)
		}

	}

	return nil

}

// postInit call the PostInit method (if defined).
func (self *instance) postInit(stack *componentStack) error {

	if self == nil {
		return nil
	}

	postInit := self.value.MethodByName("PostInit")
	if !postInit.IsValid() {
		return nil
	}

	if out, err := self.producer.container.callInjected(postInit, stack); err != nil {

		return fmt.Errorf("Error during PostInit call of '%v': %w", self, err)

	} else if len(out) > 1 {

		return fmt.Errorf("The PostInit method of '%v' should return none or one output, not %d.", self, len(out))

	} else if len(out) == 1 {

		if err, ok := out[0].Interface().(error); !ok {
			return fmt.Errorf("The output of the PostInit method of '%v' should be an error, not a '%v'.", self, out[0].Type())
		} else if err != nil {
			return fmt.Errorf("Error returned by PostInit of '%v': %w", self, err)
		} else {
			return nil
		}

	} else {

		return nil

	}

}

// castToCloser checks if the instance can be casted to an io.Closer.
func (self *instance) castToCloser() (io.Closer, bool) {
	if self == nil {
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
	if self == nil {
		return "nil"
	} else {
		return fmt.Sprintf("%v", self.value.Interface())
	}
}

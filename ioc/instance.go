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

// isNil returns whether the value is nil or not.
func (self *instance) isNil() bool {

	if self == nil {
		return true
	}

	if !self.value.IsValid() {
		return true
	}

	kind := self.value.Kind()
	if kind != reflect.Chan &&
		kind != reflect.Func &&
		kind != reflect.Interface &&
		kind != reflect.Map &&
		kind != reflect.Pointer &&
		kind != reflect.Slice {
		return false
	}

	return self.value.IsNil()

}

// initialize initializes the instance: if the instance is a struct or a
// pointer to a struct, each field is injected.
func (self *instance) initialize() error {

	if self == nil {
		return nil
	} else if err := self.producer.container.inject(self.value, true); err != nil {
		return fmt.Errorf("Error during initialisation of '%v': %w", self.producer, err)
	} else {
		return nil
	}

}

// postInit call the PostInit method (if defined).
func (self *instance) postInit() error {

	if self == nil {
		return nil
	}

	method := self.value.MethodByName("PostInit")
	if !method.IsValid() {
		return nil
	}

	wrap := func(err error) error {
		return fmt.Errorf("Error during post-initialization of '%v': %w", self.producer, err)
	}

	if out, err := self.producer.container.callInjected(method); err != nil {

		return wrap(err)

	} else if len(out) > 1 {

		return wrap(fmt.Errorf("The post-init method should return none or one output, not %d.", len(out)))

	} else if len(out) == 1 {

		if err, ok := out[0].Interface().(error); !ok {
			return wrap(fmt.Errorf("The output of the post-init method should be an error, not a '%v'.", out[0].Type()))
		} else if err != nil {
			return wrap(err)
		} else {
			return nil
		}

	} else {

		return nil

	}

}

// precall call the precall method (panics if not defined).
func (self *instance) precall(method reflect.Value) error {

	if self == nil {
		return nil
	} else if awared, ok := self.value.Interface().(PreCallAwared); !ok {
		panic(fmt.Sprintf("The component '%v' (%v) should be a PreCallAwared.", self, self.value.Type()))
	} else if err := awared.Precall(method); err != nil {
		return fmt.Errorf("Error during call of precall of '%v': %w", self, err)
	} else {
		return nil
	}

}

// postinst call the postinst method (panics if not defined).
func (self *instance) postinst(method reflect.Value, args []reflect.Value) error {

	if self == nil {
		return nil
	} else if awared, ok := self.value.Interface().(PostInstAwared); !ok {
		panic(fmt.Sprintf("The component '%v' should be a PostInstAwared.", self))
	} else if err := awared.Postinst(method, args); err != nil {
		return fmt.Errorf("Error during call of postinst of '%v': %w", self, err)
	} else {
		return nil
	}

}

// preclose close the preclose method (panics if not defined).
func (self *instance) preclose() {

	if self == nil {
		return
	} else if awared, ok := self.value.Interface().(PreCloseAwared); !ok {
		panic(fmt.Sprintf("The component '%v' should be a PreCloseAwared.", self))
	} else {

		defer func() {
			_ = recover()
		}()

		awared.Preclose()

	}

}

// postclose close the postclose method (panics if not defined).
func (self *instance) postclose() {

	if self == nil {
		return
	} else if awared, ok := self.value.Interface().(PostCloseAwared); !ok {
		panic(fmt.Sprintf("The component '%v' should be a PostCloseAwared.", self))
	} else {

		defer func() {
			_ = recover()
		}()

		awared.Postclose()

	}

}

func (self *instance) isClosable() bool {
	if self == nil {
		return false
	} else {
		_, closable := self.value.Interface().(io.Closer)
		return closable
	}
}

// close call the Close method (if defined).
func (self *instance) close() {

	if self == nil {
		return
	} else if closer, ok := self.value.Interface().(io.Closer); ok {

		defer func() {
			_ = recover()
		}()

		_ = closer.Close()

	}

}

func (self *instance) String() string {
	if self == nil {
		return "nil"
	} else {
		return self.producer.name
	}
}

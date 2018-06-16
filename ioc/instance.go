package ioc

import (
	"io"
	"reflect"

	"github.com/pkg/errors"
)

// instance represents a Component instance. It's created by the factory of the
// Component.
type Instance struct {
	value    reflect.Value
	producer *Component
}

// voidInstance creates an instance without value.
func voidInstance(producer *Component) *Instance {
	return &Instance{reflect.Value{}, producer}
}

// newInstance creates a new instance.
func newInstance(instance reflect.Value, producer *Component) *Instance {
	return &Instance{instance, producer}
}

// isNil returns whether the value is nil or not.
func (self *Instance) isNil() bool {

	if !self.value.IsValid() {
		return true
	}

	kind := self.value.Kind()
	if kind != reflect.Chan &&
		kind != reflect.Func &&
		kind != reflect.Interface &&
		kind != reflect.Map &&
		kind != reflect.Ptr &&
		kind != reflect.Slice {
		return false
	}

	return self.value.IsNil()

}

// initialize initializes the instance: if the instance is a struct or a
// pointer to a struct, each field is injected.
func (self *Instance) initialize(container *Container) error {

	err := container.inject(self.value, true)
	return errors.Wrapf(err, "Error during initialisation of '%v'", self.producer)

}

// postInit call the PostInit method (if defined).
func (self *Instance) postInit(container *Container) error {

	wrap := func(err error) error {
		return errors.Wrapf(err, "Error during post-initialization of '%v'", self.producer)
	}

	method := self.value.MethodByName("PostInit")
	if !method.IsValid() {
		return nil
	}

	out, err := container.callInjected(method)
	if err != nil {
		return wrap(err)
	}

	if len(out) == 1 {
		if err, ok := out[0].Interface().(error); !ok {
			return wrap(errors.Errorf("The output of the post-init method should be an error, not a '%v'.", out[0].Type()))
		} else if err != nil {
			return wrap(err)
		}
	}

	if len(out) > 1 {
		return wrap(errors.Errorf("The post-init method should return none or one output, not %d.", len(out)))
	}

	return nil

}

// close call the Close method (if defined).
func (self *Instance) close(container *Container) {

	if closer, ok := self.value.Interface().(io.Closer); ok {
		_ = closer.Close()
	}

}

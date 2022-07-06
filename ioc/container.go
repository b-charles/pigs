package ioc

import (
	"fmt"
	"reflect"
)

// CONTAINER

type Container struct {
	coreComponents map[string][]*Component
	testComponents map[string][]*Component
	instances      map[*Component]*Instance
	closables      []*Instance
}

// NewContainer creates a pointer to new Container. This new container contains
// already one component: itself with the name 'ApplicationContainer'.
func NewContainer() *Container {

	container := &Container{
		make(map[string][]*Component),
		make(map[string][]*Component),
		make(map[*Component]*Instance),
		make([]*Instance, 0)}

	container.Put(container)

	return container

}

// REGISTRATION

// putIn creates a Component (by its factory, name and aliases) and put it in
// the given map.
func (self *Container) putIn(
	factory interface{},
	name string,
	aliases []string,
	components map[string][]*Component) error {

	component, err := NewComponent(factory, name, aliases)
	if err != nil {
		return err
	}

	for _, alias := range component.aliases {

		list, ok := components[alias]
		if !ok {
			list = make([]*Component, 0)
		}

		components[alias] = append(list, component)

	}

	return nil

}

// wrap converts an object into a function which returns this object.
func wrap[T any](object T) func() T {
	return func() T {
		return object
	}
}

// defaultName returns the default name of a component type.
func defaultComponentName(typ reflect.Type) string {

	if typ.Kind() == reflect.Pointer {
		typ = typ.Elem()
	}

	return typ.String()

}

// defaultName returns the default name of a component type.
func defaultFactoryName(t reflect.Type) (string, error) {

	if t.Kind() != reflect.Func {
		return "", fmt.Errorf("The type %v should be a function.", t)
	}

	o := t.NumOut()
	if o == 0 {
		return "", fmt.Errorf("The function should return at least one value.")
	}

	return defaultComponentName(t.Out(0)), nil

}

// PutNamedFactory records a new Component by its factory and its name.
func (self *Container) PutNamedFactory(factory interface{}, name string, aliases ...string) error {
	return self.putIn(factory, name, aliases, self.coreComponents)
}

// PutNamed records a new Component by its value and its name.
func (self *Container) PutNamed(object interface{}, name string, aliases ...string) error {
	return self.PutNamedFactory(wrap(object), name, aliases...)
}

// PutFactory records a new Component by its factory.
func (self *Container) PutFactory(factory interface{}, aliases ...string) error {
	name, err := defaultFactoryName(reflect.TypeOf(factory))
	if err != nil {
		return fmt.Errorf("Can not register factory '%v': %w", factory, err)
	}
	return self.putIn(factory, name, aliases, self.coreComponents)
}

// Put records a new Component by its value.
func (self *Container) Put(object interface{}, aliases ...string) error {
	name := defaultComponentName(reflect.TypeOf(object))
	return self.PutNamedFactory(wrap(object), name, aliases...)
}

// TestPutNamedFactory records a new Component for tests by its factory and its name.
func (self *Container) TestPutNamedFactory(factory interface{}, name string, aliases ...string) error {
	return self.putIn(factory, name, aliases, self.testComponents)
}

// TestNamedPut records a new Component for tests by its value and its name.
func (self *Container) TestPutNamed(object interface{}, name string, aliases ...string) error {
	return self.TestPutNamedFactory(wrap(object), name, aliases...)
}

// TestPutFactory records a new Component for tests by its factory.
func (self *Container) TestPutFactory(factory interface{}, aliases ...string) error {
	name, err := defaultFactoryName(reflect.TypeOf(factory))
	if err != nil {
		return fmt.Errorf("Can not register factory '%v': %w", factory, err)
	}
	return self.putIn(factory, name, aliases, self.testComponents)
}

// TestPut records a new Component for tests by its value.
func (self *Container) TestPut(object interface{}, aliases ...string) error {
	name := defaultComponentName(reflect.TypeOf(object))
	return self.TestPutNamedFactory(wrap(object), name, aliases...)
}

// INTERNAL RESOLUTION

// getComponentInstance gets the instance of the given component. If the
// instance is already created, the instance is returned. If not, the instance is
// created, recorded, initialized, post-initialized and returned.
func (self *Container) getComponentInstance(component *Component) (*Instance, error) {

	if instance, ok := self.instances[component]; ok {
		return instance, nil
	}

	instance, err := component.instanciate(self)
	if err != nil {
		return instance, err
	}

	self.instances[component] = instance

	if err := instance.initialize(self); err != nil {
		return instance, err
	}
	if err := instance.postInit(self); err != nil {
		return instance, err
	}

	if instance.isClosable() {
		self.closables = append(self.closables, instance)
	}

	return instance, nil

}

// extractInstances get not nil instances from a name and a map of Components.
func (self *Container) extractInstances(name string, components map[string][]*Component) ([]*Instance, error) {

	list, ok := components[name]
	if !ok {
		return []*Instance{}, nil
	}

	instances := make([]*Instance, 0, len(list))

	for _, component := range list {

		instance, err := self.getComponentInstance(component)
		if err != nil {
			return instances, err
		}

		if !instance.isNil() {
			instances = append(instances, instance)
		}

	}

	return instances, nil

}

// getComponentInstances get the instance from the test map Component, or the
// core map if nothing is found.
func (self *Container) getComponentInstances(name string) ([]*Instance, error) {

	instances, err := self.extractInstances(name, self.testComponents)
	if err != nil {
		return instances, err
	}
	if len(instances) > 0 {
		return instances, nil
	}

	instances, err = self.extractInstances(name, self.coreComponents)
	if err != nil {
		return instances, err
	}
	return instances, nil

}

// INJECTION

var string_type reflect.Type = reflect.TypeOf("")

// packInstance converts one instance to target type (direct, ptr or addr).
func packInstance(instance *Instance, target reflect.Type) (reflect.Value, error) {

	value := instance.value
	typ := value.Type()

	if typ.AssignableTo(target) {
		return value, nil
	}

	if typ.Kind() == reflect.Pointer && typ.Elem().AssignableTo(target) {
		return value.Elem(), nil
	}

	if value.CanAddr() && reflect.PointerTo(typ).AssignableTo(target) {
		return value.Addr(), nil
	}

	return reflect.Value{}, fmt.Errorf("The component '%v' (%v) can not be assigned to %v.", instance, typ, target)

}

// pack converts instance slice (not empty) to target type (direct, ptr, addr, slice or map).
func packInstances(instances []*Instance, target reflect.Type) (reflect.Value, error) {

	if len(instances) == 0 {
		return reflect.Value{}, fmt.Errorf("No component found.")
	}

	if len(instances) == 1 && instances[0].value.Type().AssignableTo(target) {
		return instances[0].value, nil
	}

	if target.Kind() == reflect.Slice {

		target = target.Elem()

		slice := reflect.MakeSlice(reflect.SliceOf(target), 0, len(instances))
		for _, instance := range instances {

			value, err := packInstance(instance, target)
			if err != nil {
				return reflect.Value{}, err
			}

			slice = reflect.Append(slice, value)

		}

		return slice, nil

	}

	if target.Kind() == reflect.Map && target.Key() == string_type {

		target = target.Elem()

		m := reflect.MakeMapWithSize(reflect.MapOf(string_type, target), len(instances))
		for _, instance := range instances {

			value, err := packInstance(instance, target)
			if err != nil {
				return reflect.Value{}, err
			}

			m.SetMapIndex(reflect.ValueOf(instance.producer.name), value)

		}

		return m, nil

	}

	if len(instances) > 1 {

		producers := make([]*Component, 0)
		for _, instance := range instances {
			producers = append(producers, instance.producer)
		}

		return reflect.Value{}, fmt.Errorf("Too many components found: %v.", producers)

	}

	return packInstance(instances[0], target)

}

// inject each field of a value, only if the value is a struct or a pointer to a struct.
func (self *Container) inject(value reflect.Value, onlyTagged bool) error {

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

		nameField := structField.Name
		typeField := structField.Type

		name, ok := structField.Tag.Lookup("inject")
		if !ok && onlyTagged {
			continue
		}

		if !field.CanSet() {
			return fmt.Errorf("The field '%v' of %v is not settable.", nameField, typ)
		}

		var instances []*Instance
		var err error

		if name != "" {

			instances, err = self.getComponentInstances(name)
			if err != nil {
				return err
			}

			if len(instances) == 0 {
				return fmt.Errorf("Can not inject field %v of %v: No component '%v' found.", nameField, typ, name)
			}

		} else {

			name = defaultComponentName(typeField)
			instances, err = self.getComponentInstances(name)
			if err != nil {
				return err
			}

			if len(instances) == 0 {
				instances, err = self.getComponentInstances(nameField)
				if err != nil {
					return err
				}
			}

			if len(instances) == 0 {
				return fmt.Errorf("Can not inject field %v of %v: No component '%v' or '%v' found.", nameField, typ, name, nameField)
			}

		}

		if val, err := packInstances(instances, structField.Type); err != nil {
			return fmt.Errorf("Can not inject field %v of %v: %w", nameField, typ, err)
		} else {
			field.Set(val)
		}

	}

	return nil

}

// recoveredCall calls a method and recover in case of panic.
func recoveredCall(method reflect.Value, args []reflect.Value) (out []reflect.Value, err error) {

	defer func() {
		if r, ok := recover().(error); ok {
			err = fmt.Errorf("Error while calling %v, %w", method, r)
		}
	}()

	return method.Call(args), nil

}

// callInjected call the given method, injecting its arguments.
func (self *Container) callInjected(method reflect.Value) ([]reflect.Value, error) {

	if method.Kind() != reflect.Func {
		return nil, fmt.Errorf("Can not use '%v' as a function.", method)
	}

	methodType := method.Type()

	numIn := methodType.NumIn()
	args := make([]reflect.Value, numIn)

	if numIn == 1 {

		argType := methodType.In(0)
		argName := defaultComponentName(argType)

		instances, err := self.getComponentInstances(argName)
		if err != nil {
			return nil, err
		}

		if len(instances) != 0 {

			args[0], err = packInstances(instances, argType)
			if err != nil {
				return nil, err
			}

		} else {

			ptr := argType.Kind() == reflect.Pointer
			if ptr {
				argType = argType.Elem()
			}

			if argType.Kind() == reflect.Struct {

				arg := reflect.New(argType).Elem()
				if err := self.inject(arg, false); err != nil {
					return nil, err
				}

				if ptr {
					arg = arg.Addr()
				}

				args[0] = arg

			} else {

				return nil, fmt.Errorf("No component '%v' found.", argName)

			}

		}

	} else {

		for i := 0; i < numIn; i++ {

			argType := methodType.In(i)
			argName := defaultComponentName(argType)

			instances, err := self.getComponentInstances(argName)
			if err != nil {
				return nil, err
			}

			if len(instances) == 0 {
				return nil,
					fmt.Errorf("Can not inject parameter #%d: No component '%v' found.", i, argName)
			}

			args[i], err = packInstances(instances, argType)
			if err != nil {
				return nil,
					fmt.Errorf("Can not inject parameter #%d: %w", i, err)
			}

		}

	}

	return recoveredCall(method, args)

}

// EXTERNAL RESOLUTION

// CallInjected call the given method, injecting its arguments.
func (self *Container) CallInjected(method interface{}) error {

	out, err := self.callInjected(reflect.ValueOf(method))
	if err != nil {
		return err
	}

	if len(out) == 1 {
		if err, ok := out[0].Interface().(error); !ok {
			return fmt.Errorf("The output of the method should be an error, not a '%v'.", out[0].Type())
		} else if err != nil {
			return err
		}
	}

	if len(out) > 1 {
		return fmt.Errorf("The method should return none or one output, not %d.", len(out))
	}

	return nil

}

// CLOSE

// Close close all components.
func (self *Container) Close() {

	for _, instance := range self.closables {

		instance.close(self)

		delete(self.instances, instance.producer)

	}

	self.closables = nil

}

// ClearTests close all components and delete all registered Component for tests.
func (self *Container) ClearTests() {

	self.Close()

	self.testComponents = make(map[string][]*Component)

}

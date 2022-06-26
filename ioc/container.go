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
	sorted         []*Instance
}

// NewContainer creates a pointer to new Container. This new container contains
// already one component: itself with the name 'ApplicationContainer'.
func NewContainer() *Container {

	container := &Container{
		make(map[string][]*Component),
		make(map[string][]*Component),
		make(map[*Component]*Instance),
		make([]*Instance, 0)}

	container.Put(container, "ApplicationContainer")

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

// wrap convert an object into a function which returns this object.
func (self *Container) wrap(object interface{}) func() interface{} {
	return func() interface{} {
		return object
	}
}

// PutFactory records a new Component by its factory.
func (self *Container) PutFactory(factory interface{}, name string, aliases ...string) error {
	return self.putIn(factory, name, aliases, self.coreComponents)
}

// Put records a new Component by its value.
func (self *Container) Put(object interface{}, name string, aliases ...string) error {
	return self.PutFactory(self.wrap(object), name, aliases...)
}

// PutFactory records a new Component for tests by its factory.
func (self *Container) TestPutFactory(factory interface{}, name string, aliases ...string) error {
	return self.putIn(factory, name, aliases, self.testComponents)
}

// Put records a new Component for tests by its value.
func (self *Container) TestPut(object interface{}, name string, aliases ...string) error {
	return self.TestPutFactory(self.wrap(object), name, aliases...)
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

	self.sorted = append(self.sorted, instance)

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

	wrap := func(err error) error {
		return fmt.Errorf("Error during instanciation of '%s', %w", name, err)
	}

	instances, err := self.extractInstances(name, self.testComponents)
	if err != nil {
		return instances, wrap(err)
	}
	if len(instances) > 0 {
		return instances, nil
	}

	instances, err = self.extractInstances(name, self.coreComponents)
	if err != nil {
		return instances, wrap(err)
	}
	return instances, nil

}

// INJECTION

var string_type reflect.Type = reflect.TypeOf("")

// resolve get Components with the given name and construct a value of the
// given type.
func (self *Container) resolve(name string, class reflect.Type) (reflect.Value, error) {

	instances, err := self.getComponentInstances(name)
	if err != nil {
		return reflect.Value{}, err
	} else if len(instances) == 0 {
		return reflect.Value{}, fmt.Errorf("No producer found for '%s'.", name)
	}

	if instanceType := instances[0].value.Type(); !instanceType.AssignableTo(class) {

		if class.Kind() == reflect.Slice {

			class = class.Elem()

			slice := reflect.MakeSlice(reflect.SliceOf(class), 0, len(instances))
			for _, instance := range instances {
				slice = reflect.Append(slice, instance.value)
			}

			return slice, nil

		} else if class.Kind() == reflect.Map {

			if keyClass := class.Key(); keyClass != string_type {
				return reflect.Value{}, fmt.Errorf("Unsupported key type for a map injection: %v, only 'string' is valid.", keyClass)
			}

			class = class.Elem()

			m := reflect.MakeMapWithSize(reflect.MapOf(string_type, class), len(instances))
			for _, instance := range instances {
				m.SetMapIndex(reflect.ValueOf(instance.producer.name), instance.value)
			}

			return m, nil

		}

	}

	if len(instances) > 1 {

		producers := make([]*Component, 0)
		for _, instance := range instances {
			producers = append(producers, instance.producer)
		}

		return reflect.Value{}, fmt.Errorf("Too many producers found for '%s': %v.", name, producers)

	}

	return instances[0].value, nil

}

// inject each field of a value, only if the value is a struct or a pointer to
// a struct.
func (self *Container) inject(value reflect.Value, onlyTagged bool) error {

	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}
	if value.Kind() != reflect.Struct {
		return nil
	}

	class := value.Type()

	for i := 0; i < value.NumField(); i++ {

		field := value.Field(i)
		fieldType := class.Field(i)

		name, ok := fieldType.Tag.Lookup("inject")
		if !ok && onlyTagged {
			continue
		}

		if name == "" {
			name = fieldType.Name
		}

		if !field.CanSet() {
			return fmt.Errorf("The field %v of %v is not settable.", fieldType, class)
		}

		if val, err := self.resolve(name, fieldType.Type); err != nil {
			return err
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
		return []reflect.Value{}, fmt.Errorf("The argument is not a function: %v", method)
	}

	methodType := method.Type()

	numIn := methodType.NumIn()
	args := make([]reflect.Value, numIn)

	out := []reflect.Value{}

	injected := false

	if numIn == 1 {

		argType := methodType.In(0)

		ptr := argType.Kind() == reflect.Ptr
		if ptr {
			argType = argType.Elem()
		}

		name := argType.Name()
		_, testPresence := self.testComponents[name]
		_, corePresence := self.coreComponents[name]
		if !testPresence && !corePresence {

			arg := reflect.New(argType).Elem()

			if err := self.inject(arg, false); err != nil {
				return out, err
			}

			if ptr {
				arg = arg.Addr()
			}

			args[0] = arg

			injected = true

		}

	}

	if !injected {

		for i := 0; i < numIn; i++ {

			argType := methodType.In(i)

			argName := argType.Name()
			if argType.Kind() == reflect.Ptr {
				argName = argType.Elem().Name()
			}

			if resolved, err := self.resolve(argName, argType); err != nil {
				return out, err
			} else {
				args[i] = resolved
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

	for _, instance := range self.sorted {

		instance.close(self)

		delete(self.instances, instance.producer)

	}

	self.sorted = nil

}

// ClearTests close all components and delete all registered Component for tests.
func (self *Container) ClearTests() {

	self.Close()

	self.testComponents = make(map[string][]*Component)

}

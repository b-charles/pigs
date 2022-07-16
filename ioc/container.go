package ioc

import (
	"fmt"
	"reflect"
	"time"
)

// Awareness

var preCallAwaredName = unsafeDefaultAlias(func(PreCallAwared) {})
var postInstAwaredName = unsafeDefaultAlias(func(PostInstAwared) {})
var preCloseAwaredName = unsafeDefaultAlias(func(PreCloseAwared) {})
var postCloseAwaredName = unsafeDefaultAlias(func(PostCloseAwared) {})

// CONTAINER

type Container struct {
	coreComponents  map[string][]*component
	testComponents  map[string][]*component
	creationTime    time.Time
	instances       map[*component]*instance
	preCloseAwared  []*instance
	postCloseAwared []*instance
	closables       []*instance
}

// NewContainer creates a pointer to a new initialized Container.
func NewContainer() *Container {

	container := &Container{
		coreComponents:  make(map[string][]*component),
		testComponents:  make(map[string][]*component),
		creationTime:    time.Now(),
		instances:       map[*component]*instance{},
		preCloseAwared:  []*instance{},
		postCloseAwared: []*instance{},
		closables:       []*instance{}}

	container.Put(container)

	container.Put(&noopPreCallAwaredHandler{}, preCallAwaredName)
	container.Put(&noopPostInstAwaredHandler{}, postInstAwaredName)
	container.Put(&noopPreCloseAwaredHandler{}, preCloseAwaredName)
	container.Put(&noopPostCloseAwaredHandler{}, postCloseAwaredName)

	return container

}

// CreationTime returns the creation time of the container.
func (self *Container) CreationTime() time.Time {
	return self.creationTime
}

// REGISTRATION

// putIn creates a component (by its factory, name and aliases) and put it in
// the given map.
func (self *Container) putIn(factory any, name string, aliases []string, core bool) error {

	var sink map[string][]*component
	if core {
		sink = self.coreComponents
	} else {
		sink = self.testComponents
	}

	if old, present := sink[name]; present {
		return fmt.Errorf("Two components '%v' and '%v' are registered with the same main name '%s'.", old, factory, name)
	}

	comp, err := newComponent(self, factory, name, aliases)
	if err != nil {
		return fmt.Errorf("Error during registration of '%s': %w", name, err)
	}

	sink[name] = []*component{comp}

	for _, alias := range comp.aliases {

		list, ok := sink[alias]
		if !ok {
			list = make([]*component, 0)
		}

		sink[alias] = append(list, comp)

	}

	return nil

}

// wrap converts a component into a factory which returns this component.
func wrap[T any](component T) func() T {
	return func() T {
		return component
	}
}

// PutNamedFactory records a new component by its factory and its name.
func (self *Container) PutNamedFactory(factory any, name string, aliases ...any) error {

	allAliases, err := defaultAliases(aliases...)
	if err != nil {
		return fmt.Errorf("Can not register factory '%v': %w", name, err)
	}

	return self.putIn(factory, name, allAliases, true)

}

// PutNamed records a new component by its value and its name.
func (self *Container) PutNamed(component any, name string, aliases ...any) error {

	allAliases, err := defaultAliases(aliases...)
	if err != nil {
		return fmt.Errorf("Can not register component '%v': %w", name, err)
	}

	return self.putIn(wrap(component), name, allAliases, true)

}

// PutFactory records a new component by its factory.
func (self *Container) PutFactory(factory any, aliases ...any) error {

	name, err := defaultFactoryName(reflect.TypeOf(factory))
	if err != nil {
		return fmt.Errorf("Can not register factory '%v': %w", factory, err)
	}

	allAliases, err := defaultAliases(aliases...)
	if err != nil {
		return fmt.Errorf("Can not register factory '%v': %w", name, err)
	}

	return self.putIn(factory, name, allAliases, true)

}

// Put records a new component by its value.
func (self *Container) Put(component any, aliases ...any) error {

	name := defaultComponentName(reflect.TypeOf(component))

	allAliases, err := defaultAliases(aliases...)
	if err != nil {
		return fmt.Errorf("Can not register component '%v': %w", name, err)
	}

	return self.putIn(wrap(component), name, allAliases, true)

}

// TestPutNamedFactory records a new component for tests by its factory and its name.
func (self *Container) TestPutNamedFactory(factory any, name string, aliases ...any) error {

	allAliases, err := defaultAliases(aliases...)
	if err != nil {
		return fmt.Errorf("Can not register test factory '%v': %w", name, err)
	}

	return self.putIn(factory, name, allAliases, false)

}

// TestNamedPut records a new component for tests by its value and its name.
func (self *Container) TestPutNamed(component any, name string, aliases ...any) error {

	allAliases, err := defaultAliases(aliases...)
	if err != nil {
		return fmt.Errorf("Can not register test component '%v': %w", name, err)
	}

	return self.putIn(wrap(component), name, allAliases, false)

}

// TestPutFactory records a new component for tests by its factory.
func (self *Container) TestPutFactory(factory any, aliases ...any) error {

	name, err := defaultFactoryName(reflect.TypeOf(factory))
	if err != nil {
		return fmt.Errorf("Can not register test factory '%v': %w", factory, err)
	}

	allAliases, err := defaultAliases(aliases...)
	if err != nil {
		return fmt.Errorf("Can not register test factory '%v': %w", name, err)
	}

	return self.putIn(factory, name, allAliases, false)

}

// TestPut records a new component for tests by its value.
func (self *Container) TestPut(component any, aliases ...any) error {

	name := defaultComponentName(reflect.TypeOf(component))

	allAliases, err := defaultAliases(aliases...)
	if err != nil {
		return fmt.Errorf("Can not register test component '%v': %w", name, err)
	}

	return self.putIn(wrap(component), name, allAliases, false)

}

// INTERNAL RESOLUTION

// getComponentInstance gets the instance of the given component. If the
// instance is already created, the instance is returned. If not, the instance is
// created, recorded, initialized, post-initialized and returned.
func (self *Container) getComponentInstance(component *component) (*instance, error) {

	if instance, ok := self.instances[component]; ok {
		return instance, nil
	}

	instance, err := component.instanciate()
	if err != nil {
		return instance, err
	}

	self.instances[component] = instance

	if err := instance.initialize(); err != nil {
		return instance, err
	}
	if err := instance.postInit(); err != nil {
		return instance, err
	}

	if instance.isClosable() {
		self.closables = append(self.closables, instance)
	}

	return instance, nil

}

// extractInstances get not nil instances from a name and a map of Components.
func (self *Container) extractInstances(name string, components map[string][]*component) ([]*instance, error) {

	list, ok := components[name]
	if !ok {
		return []*instance{}, nil
	}

	instances := make([]*instance, 0, len(list))

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

// getComponentInstances get the instance from the test map component, or the
// core map if nothing is found.
func (self *Container) getComponentInstances(name string) ([]*instance, error) {

	instances, err := self.extractInstances(name, self.testComponents)
	if err != nil || len(instances) > 0 {
		return instances, err
	}

	return self.extractInstances(name, self.coreComponents)

}

// INJECTION

var string_type reflect.Type = reflect.TypeOf("")
var error_type = reflect.TypeOf(func(error) {}).In(0)

// packInstance converts one instance to target type (direct, ptr or addr).
func packInstance(instance *instance, target reflect.Type) (reflect.Value, error) {

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
func packInstances(instances []*instance, target reflect.Type) (reflect.Value, error) {

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

		producers := make([]*component, 0)
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

		var instances []*instance
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

// getInjectedArguments returns initialized and injected arguments to call the given method.
func (self *Container) getArguments(method reflect.Value) ([]reflect.Value, error) {

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

	return args, nil

}

// callInjected call the given method, injecting its arguments.
func (self *Container) callInjected(method reflect.Value) ([]reflect.Value, error) {

	args, err := self.getArguments(method)
	if err != nil {
		return nil, err
	}

	return method.Call(args), nil

}

// EXTERNAL RESOLUTION

// callEntryPoint call the given method, injecting its arguments.
func (self *Container) CallInjected(method any) error {

	// input checks

	methodValue := reflect.ValueOf(method)

	typ := methodValue.Type()
	if typ.Kind() != reflect.Func {
		return fmt.Errorf("The input should a function, not a %v.", typ)
	}
	numOut := typ.NumOut()
	if numOut > 1 {
		return fmt.Errorf("The method should return none or one output, not %d.", numOut)
	} else if numOut == 1 {
		if outType := typ.Out(0); outType != error_type {
			return fmt.Errorf("The output of the method should be an error, not a '%v'.", outType)
		}
	}

	// get life cycle aware instances and injected arguments

	preCallAwared, err := self.getComponentInstances(preCallAwaredName)
	if err != nil {
		self.closeInstances()
		return err
	}
	for _, awared := range preCallAwared {
		e := awared.precall(methodValue)
		if e != nil {
			self.closeInstances()
			return e
		}
	}

	postInstAwared, err := self.getComponentInstances(postInstAwaredName)
	if err != nil {
		self.closeInstances()
		return err
	}

	self.preCloseAwared, err = self.getComponentInstances(preCloseAwaredName)
	if err != nil {
		self.closeInstances()
		return err
	}

	self.postCloseAwared, err = self.getComponentInstances(postCloseAwaredName)
	if err != nil {
		self.closeInstances()
		return err
	}

	args, err := self.getArguments(methodValue)
	if err != nil {
		return err
	}

	// release unused instances
	self.instances = make(map[*component]*instance)
	if len(self.testComponents) == 0 {
		self.coreComponents = make(map[string][]*component)
	} else {
		self.testComponents = make(map[string][]*component)
	}

	// postinst
	for _, awared := range postInstAwared {
		e := awared.postinst(methodValue, args)
		if e != nil {
			self.closeInstances()
			return e
		}
	}

	// calling
	outs := methodValue.Call(args)

	// close
	self.closeInstances()

	// output

	if numOut == 1 {
		if err := outs[0].Interface().(error); err != nil {
			return err
		}
	}

	return nil

}

func (self *Container) closeInstances() {

	for _, awared := range self.preCloseAwared {
		awared.preclose()
	}

	for _, instance := range self.closables {
		instance.close()
	}

	for _, awared := range self.postCloseAwared {
		awared.postclose()
	}

}

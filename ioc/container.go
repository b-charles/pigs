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
	coreComponents map[string][]*Component
	testComponents map[string][]*Component
	instances      map[*Component]*Instance
	closables      []*Instance
	creationTime   time.Time
}

// NewContainer creates a pointer to a new initialized Container.
func NewContainer() *Container {

	container := &Container{
		make(map[string][]*Component),
		make(map[string][]*Component),
		make(map[*Component]*Instance),
		make([]*Instance, 0),
		time.Now()}

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

// DEFAULT NAMING

// defaultName returns the default name of a component type.
func defaultComponentName(typ reflect.Type) string {

	if typ.Kind() == reflect.Pointer || typ.Kind() == reflect.Slice ||
		(typ.Kind() == reflect.Map && typ.Key() == string_type) {
		typ = typ.Elem()
	}

	name := typ.Name()
	if name == "" {
		return typ.String()
	}

	pkg := typ.PkgPath()
	if pkg == "" {
		return name
	}

	return fmt.Sprintf("%s/%s", pkg, name)

}

// defaultName returns the default name of a component type.
func defaultFactoryName(typ reflect.Type) (string, error) {

	if typ.Kind() != reflect.Func {
		return "", fmt.Errorf("The type %v should be a function.", typ)
	}

	o := typ.NumOut()
	if o == 0 {
		return "", fmt.Errorf("The function should return at least one value.")
	}

	return defaultComponentName(typ.Out(0)), nil

}

// defaultAlias returns the aliases of the given argument.
func defaultAlias(alias any) ([]string, error) {

	value := reflect.ValueOf(alias)
	typ := value.Type()

	if typ == string_type {
		return []string{value.String()}, nil
	}

	if typ.Kind() == reflect.Func {
		list := make([]string, 0)
		for i := 0; i < typ.NumIn(); i++ {
			list = append(list, defaultComponentName(typ.In(i)))
		}
		return list, nil
	}

	return nil, fmt.Errorf("Can not guess aliases of '%v'", alias)

}

// unsafeDefaultAlias returns the first alias, and panics in case of error.
func unsafeDefaultAlias(alias any) string {
	aliases, err := defaultAlias(alias)
	if err != nil {
		panic(err)
	}
	return aliases[0]
}

// defaultAliases returns the aliases of given arguments.
func defaultAliases(aliases ...any) ([]string, error) {

	list := make([]string, 0)

	for _, elt := range aliases {

		l, err := defaultAlias(elt)
		if err != nil {
			return nil, err
		}

		list = append(list, l...)

	}

	return list, nil

}

// REGISTRATION

// putIn creates a Component (by its factory, name and aliases) and put it in
// the given map.
func putIn(
	factory any,
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

// wrap converts a component into a function which returns this component.
func wrap[T any](component T) func() T {
	return func() T {
		return component
	}
}

// PutNamedFactory records a new Component by its factory and its name.
func (self *Container) PutNamedFactory(factory any, name string, aliases ...any) error {

	allAliases, err := defaultAliases(aliases...)
	if err != nil {
		return fmt.Errorf("Can not register factory '%v': %w", name, err)
	}

	return putIn(factory, name, allAliases, self.coreComponents)

}

// PutNamed records a new Component by its value and its name.
func (self *Container) PutNamed(component any, name string, aliases ...any) error {

	allAliases, err := defaultAliases(aliases...)
	if err != nil {
		return fmt.Errorf("Can not register component '%v': %w", name, err)
	}

	return putIn(wrap(component), name, allAliases, self.coreComponents)

}

// PutFactory records a new Component by its factory.
func (self *Container) PutFactory(factory any, aliases ...any) error {

	name, err := defaultFactoryName(reflect.TypeOf(factory))
	if err != nil {
		return fmt.Errorf("Can not register factory '%v': %w", factory, err)
	}

	allAliases, err := defaultAliases(aliases...)
	if err != nil {
		return fmt.Errorf("Can not register factory '%v': %w", name, err)
	}

	return putIn(factory, name, allAliases, self.coreComponents)

}

// Put records a new Component by its value.
func (self *Container) Put(component any, aliases ...any) error {

	name := defaultComponentName(reflect.TypeOf(component))

	allAliases, err := defaultAliases(aliases...)
	if err != nil {
		return fmt.Errorf("Can not register component '%v': %w", name, err)
	}

	return putIn(wrap(component), name, allAliases, self.coreComponents)

}

// TestPutNamedFactory records a new Component for tests by its factory and its name.
func (self *Container) TestPutNamedFactory(factory any, name string, aliases ...any) error {

	allAliases, err := defaultAliases(aliases...)
	if err != nil {
		return fmt.Errorf("Can not register test factory '%v': %w", name, err)
	}

	return putIn(factory, name, allAliases, self.testComponents)

}

// TestNamedPut records a new Component for tests by its value and its name.
func (self *Container) TestPutNamed(component any, name string, aliases ...any) error {

	allAliases, err := defaultAliases(aliases...)
	if err != nil {
		return fmt.Errorf("Can not register test component '%v': %w", name, err)
	}

	return putIn(wrap(component), name, allAliases, self.testComponents)

}

// TestPutFactory records a new Component for tests by its factory.
func (self *Container) TestPutFactory(factory any, aliases ...any) error {

	name, err := defaultFactoryName(reflect.TypeOf(factory))
	if err != nil {
		return fmt.Errorf("Can not register test factory '%v': %w", factory, err)
	}

	allAliases, err := defaultAliases(aliases...)
	if err != nil {
		return fmt.Errorf("Can not register test factory '%v': %w", name, err)
	}

	return putIn(factory, name, allAliases, self.testComponents)

}

// TestPut records a new Component for tests by its value.
func (self *Container) TestPut(component any, aliases ...any) error {

	name := defaultComponentName(reflect.TypeOf(component))

	allAliases, err := defaultAliases(aliases...)
	if err != nil {
		return fmt.Errorf("Can not register test component '%v': %w", name, err)
	}

	return putIn(wrap(component), name, allAliases, self.testComponents)

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
	if err != nil || len(instances) > 0 {
		return instances, err
	}

	return self.extractInstances(name, self.coreComponents)

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

func (self *Container) callInjected(method reflect.Value) ([]reflect.Value, error) {

	args, err := self.getArguments(method)
	if err != nil {
		return nil, err
	}

	return method.Call(args), nil

}

// EXTERNAL RESOLUTION

// CallInjected call the given method, injecting its arguments.
func (self *Container) CallInjected(method any) error {

	methodValue := reflect.ValueOf(method)

	preCallAwared, err := self.getComponentInstances(preCallAwaredName)
	if err != nil {
		self.closeInstances([]*Instance{}, []*Instance{})
		return err
	}

	postInstAwared, err := self.getComponentInstances(postInstAwaredName)
	if err != nil {
		self.closeInstances([]*Instance{}, []*Instance{})
		return err
	}

	preCloseAwared, err := self.getComponentInstances(preCloseAwaredName)
	if err != nil {
		self.closeInstances([]*Instance{}, []*Instance{})
		return err
	}

	postCloseAwared, err := self.getComponentInstances(postCloseAwaredName)
	if err != nil {
		self.closeInstances(preCloseAwared, []*Instance{})
		return err
	}

	for _, awared := range preCallAwared {
		e := awared.precall(methodValue)
		if e != nil {
			return e
		}
	}

	args, err := self.getArguments(methodValue)
	if err != nil {
		return err
	}

	for _, awared := range postInstAwared {
		e := awared.postinst(methodValue, args)
		if e != nil {
			return e
		}
	}

	outs := methodValue.Call(args)

	self.closeInstances(preCloseAwared, postCloseAwared)

	if len(outs) == 1 {
		if err, ok := outs[0].Interface().(error); !ok {
			return fmt.Errorf("The output of the method should be an error, not a '%v'.", outs[0].Type())
		} else if err != nil {
			return err
		}
	}

	if len(outs) > 1 {
		return fmt.Errorf("The method should return none or one output, not %d.", len(outs))
	}

	return nil

}

func (self *Container) closeInstances(preCloseAwared []*Instance, postCloseAwared []*Instance) {

	for _, awared := range preCloseAwared {
		awared.preclose()
	}

	for _, instance := range self.closables {
		instance.close(self)
		delete(self.instances, instance.producer)
	}

	for _, awared := range postCloseAwared {
		awared.postclose()
	}

}

// ClearTests close all components and delete all registered Component for tests.
func (self *Container) ClearTests() {

	self.instances = make(map[*Component]*Instance)
	self.testComponents = make(map[string][]*Component)

}

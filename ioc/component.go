package ioc

import (
	"fmt"
	"reflect"

	"github.com/pkg/errors"
)

type PostInitAwarable interface {
	PostInit()
}

type Closeable interface {
	Close()
}

type Component struct {
	name    string
	aliases []string
	factory reflect.Value
}

var STRING_TYPE reflect.Type = reflect.TypeOf("")

func checkInInjectedMethod(method reflect.Value) error {

	methodType := method.Type()

	numIn := methodType.NumIn()
	if numIn != 0 && numIn != 1 {
		return fmt.Errorf("The function should take none or one argument (not %d).", numIn)
	}

	if numIn == 1 {
		in := methodType.In(0)
		kindIn := in.Kind()
		if kindIn == reflect.Ptr {
			kindIn = in.Elem().Kind()
		}
		if kindIn != reflect.Struct {
			return fmt.Errorf("The function should only take a struct or a pointer to a struct as argument (not a %v)", in)
		}
	}

	return nil

}

func NewComponent(factory interface{}, name string, aliases []string) (*Component, error) {

	wrap := func(err error) error {
		return errors.Wrapf(err, "Error during registration of '%s'", name)
	}
	wrapf := func(msg string, args ...interface{}) error {
		return wrap(fmt.Errorf(msg, args...))
	}

	// get unique aliases

	uniqueAliases := make(map[string]bool)
	uniqueAliases[name] = true
	for _, alias := range aliases {
		if _, ok := uniqueAliases[alias]; ok {
			return nil, wrapf("Alias specified more than once: '%s'.", alias)
		}
		uniqueAliases[alias] = true
	}

	allAliases := make([]string, 0, len(uniqueAliases))
	for alias := range uniqueAliases {
		allAliases = append(allAliases, alias)
	}

	// check factory

	factoValue := reflect.ValueOf(factory)
	if factoValue.Kind() != reflect.Func {
		return nil, wrapf("The factory should be a function, not a %v.", factoValue.Kind())
	}

	// check in

	if err := checkInInjectedMethod(factoValue); err != nil {
		return nil, wrap(err)
	}

	// check out

	factoType := factoValue.Type()

	if numOut := factoType.NumOut(); numOut != 1 {
		return nil, wrapf("The output number of the factory should be 1, not %d.", numOut)
	}

	outKind := factoType.Out(0).Kind()
	if outKind != reflect.Chan &&
		outKind != reflect.Func &&
		outKind != reflect.Interface &&
		outKind != reflect.Map &&
		outKind != reflect.Ptr &&
		outKind != reflect.Slice {
		return nil, wrapf("The output type should be a chan, a function, an interface, a map, a pointer or a slice, not a %v.", outKind)
	}

	// return

	return &Component{name, allAliases, factoValue}, nil

}

func (self *Component) resolve(container *Container, name string, class reflect.Type) (reflect.Value, error) {

	instances, producers, err := container.getComponentInstances(name)
	if err != nil {
		return reflect.Value{}, err
	} else if len(instances) == 0 {
		return reflect.Value{}, fmt.Errorf("No producer found for '%s'.", name)
	}

	if instanceType := instances[0].Type(); !instanceType.AssignableTo(class) {

		if class.Kind() == reflect.Slice {

			class = class.Elem()

			slice := reflect.MakeSlice(reflect.SliceOf(class), 0, len(instances))
			slice = reflect.Append(slice, instances...)

			return slice, nil

		} else if class.Kind() == reflect.Map {

			if keyClass := class.Key(); keyClass != STRING_TYPE {
				panic(fmt.Errorf("Unsupported key type for a map injection: %v, only 'string' is valid.", keyClass))
			}

			class = class.Elem()

			m := reflect.MakeMapWithSize(reflect.MapOf(STRING_TYPE, class), len(instances))
			for i := 0; i < len(instances); i++ {
				m.SetMapIndex(reflect.ValueOf(producers[i].name), instances[i])
			}

			return m, nil

		}

	}

	if len(instances) > 1 {
		return reflect.Value{}, fmt.Errorf("Too many producers found for '%s': %v.", name, producers)
	}

	return instances[0], nil

}

func (self *Component) inject(container *Container, value reflect.Value, onlyTagged bool) error {

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

		if val, err := self.resolve(container, name, fieldType.Type); err != nil {
			return err
		} else {
			field.Set(val)
		}

	}

	return nil

}

func (self *Component) callInjected(container *Container, method reflect.Value) (out []reflect.Value, err error) {

	defer func() {
		if r, ok := recover().(error); ok {
			err = errors.Wrapf(r, "Error during calling %v", method)
		}
	}()

	methodType := method.Type()

	numIn := methodType.NumIn()
	args := make([]reflect.Value, numIn)
	if numIn == 1 {

		argType := methodType.In(0)

		ptr := argType.Kind() == reflect.Ptr
		if ptr {
			argType = argType.Elem()
		}

		arg := reflect.New(argType).Elem()

		if err = self.inject(container, arg, false); err != nil {
			return
		}

		if ptr {
			arg = arg.Addr()
		}

		args[0] = arg

	}

	return method.Call(args), nil

}

func (self *Component) Instanciate(container *Container) (reflect.Value, error) {

	out, err := self.callInjected(container, self.factory)
	if err != nil {
		return reflect.Value{}, errors.Wrapf(err, "Error during instanciation of '%s'", self.name)
	}

	instance := out[0]
	if instance.Kind() == reflect.Interface {
		instance = instance.Elem()
	}

	return instance, nil

}

func (self *Component) Initialize(container *Container, instance reflect.Value) error {

	if instance.IsNil() {
		return nil
	}
	if instance.Kind() == reflect.Ptr {
		instance = instance.Elem()
	}
	if instance.Kind() != reflect.Struct {
		return nil
	}

	err := self.inject(container, instance, true)
	return errors.Wrapf(err, "Error during initialisation of '%s'", self.name)

}

func (self *Component) PostInit(container *Container, instance reflect.Value) error {

	postInit := instance.MethodByName("PostInit")
	if !postInit.IsValid() {
		return nil
	}

	if err := checkInInjectedMethod(postInit); err != nil {
		return err
	}

	postInitType := postInit.Type()

	if numOut := postInitType.NumOut(); numOut != 0 {
		panic(fmt.Errorf("The method should return nothing (found %d outputs).", numOut))
	}

	_, err := self.callInjected(container, postInit)
	return errors.Wrapf(err, "Error during post-initialization of '%s'", self.name)

}

func (self *Component) Close(container *Container, instance reflect.Value) {

	if closeable, ok := instance.Interface().(Closeable); ok {
		closeable.Close()
	}

}

func (self *Component) String() string {
	return self.name
}

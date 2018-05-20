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

func checkInInjectedMethod(method reflect.Value) {

	methodType := method.Type()

	numIn := methodType.NumIn()
	if numIn != 0 && numIn != 1 {
		panic(fmt.Errorf("The function should take none or one argument (not %d).", numIn))
	}

	if numIn == 1 {
		in := methodType.In(0)
		kindIn := in.Kind()
		if kindIn == reflect.Ptr {
			kindIn = in.Elem().Kind()
		}
		if kindIn != reflect.Struct {
			panic(fmt.Errorf("The function should only take a struct or a pointer to a struct as argument (not a %v)", in))
		}
	}

}

func NewComponent(factory interface{}, name string, aliases []string) *Component {

	defer func() {
		if err, ok := recover().(error); ok {
			panic(errors.Wrapf(err, "Error during registration of '%s'", name))
		}
	}()

	// get unique aliases

	uniqueAliases := make(map[string]bool)
	uniqueAliases[name] = true
	for _, alias := range aliases {
		if _, ok := uniqueAliases[alias]; ok {
			panic(fmt.Errorf("Alias specified more than once: '%s'.", alias))
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
		panic(fmt.Errorf("The factory should be a function, not a %v.", factoValue.Kind()))
	}

	// check in

	checkInInjectedMethod(factoValue)

	// check out

	factoType := factoValue.Type()

	if numOut := factoType.NumOut(); numOut != 1 {
		panic(fmt.Errorf("The output number of the factory should be 1, not %d.", numOut))
	}

	outKind := factoType.Out(0).Kind()
	if outKind != reflect.Chan &&
		outKind != reflect.Func &&
		outKind != reflect.Interface &&
		outKind != reflect.Map &&
		outKind != reflect.Ptr &&
		outKind != reflect.Slice {
		panic(fmt.Errorf("The output type should be a chan, a function, an interface, a map, a pointer or a slice, not a %v.", outKind))
	}

	// return

	return &Component{name, allAliases, factoValue}

}

func (self *Component) resolve(container *Container, name string, class reflect.Type) reflect.Value {

	instances, producers := container.getComponentInstances(name)

	if len(instances) == 0 {
		panic(fmt.Errorf("No producer found for '%s'.", name))
	}

	if instanceType := instances[0].Type(); !instanceType.AssignableTo(class) {

		if class.Kind() == reflect.Slice {

			class = class.Elem()

			slice := reflect.MakeSlice(reflect.SliceOf(class), 0, len(instances))
			slice = reflect.Append(slice, instances...)

			return slice

		} else if class.Kind() == reflect.Map {

			if keyClass := class.Key(); keyClass != STRING_TYPE {
				panic(fmt.Errorf("Unsupported key type for a map injection: %v, only 'string' is valid.", keyClass))
			}

			class = class.Elem()

			m := reflect.MakeMapWithSize(reflect.MapOf(STRING_TYPE, class), len(instances))
			for i := 0; i < len(instances); i++ {
				m.SetMapIndex(reflect.ValueOf(producers[i].name), instances[i])
			}

			return m

		}

	}

	if len(instances) > 1 {
		panic(fmt.Errorf("Too many producers found for '%s': %v.", name, producers))
	}

	return instances[0]

}

func (self *Component) inject(container *Container, value reflect.Value, onlyTagged bool) {

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
			panic(fmt.Errorf("The field %v of %v is not settable.", fieldType, class))
		}

		field.Set(self.resolve(container, name, fieldType.Type))

	}

}

func (self *Component) callInjected(container *Container, method reflect.Value) []reflect.Value {

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

		self.inject(container, arg, false)

		if ptr {
			arg = arg.Addr()
		}

		args[0] = arg

	}

	return method.Call(args)

}

func (self *Component) Instanciate(container *Container) reflect.Value {

	defer func() {
		if err, ok := recover().(error); ok {
			panic(errors.Wrapf(err, "Error during instanciation of '%s'", self.name))
		}
	}()

	instance := self.callInjected(container, self.factory)[0]
	if instance.Kind() == reflect.Interface {
		instance = instance.Elem()
	}

	return instance

}

func (self *Component) Initialize(container *Container, instance reflect.Value) {

	defer func() {
		if err, ok := recover().(error); ok {
			panic(errors.Wrapf(err, "Error during initialisation of '%s'", self.name))
		}
	}()

	if instance.IsNil() {
		return
	}
	if instance.Kind() == reflect.Ptr {
		instance = instance.Elem()
	}
	if instance.Kind() != reflect.Struct {
		return
	}

	self.inject(container, instance, true)

}

func (self *Component) PostInit(container *Container, instance reflect.Value) {

	postInit := instance.MethodByName("PostInit")
	if !postInit.IsValid() {
		return
	}

	checkInInjectedMethod(postInit)

	postInitType := postInit.Type()

	if numOut := postInitType.NumOut(); numOut != 0 {
		panic(fmt.Errorf("The method should return nothing (found %d outputs).", numOut))
	}

	self.callInjected(container, postInit)

}

func (self *Component) Close(container *Container, instance reflect.Value) {

	if closeable, ok := instance.Interface().(Closeable); ok {
		closeable.Close()
	}

}

func (self *Component) String() string {
	return self.name
}

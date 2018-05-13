package ioc

import (
	"reflect"
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
	inputs  []string
}

var STRING_TYPE reflect.Type = reflect.TypeOf("")

func NewComponent(factory interface{}, inputs []string, name string, aliases []string) *Component {

	defer func() {
		if err := recover(); err != nil {
			panic(&ComponentRegistrationError{name, err})
		}
	}()

	// get unique aliases

	uniqueAliases := make(map[string]bool)
	uniqueAliases[name] = true
	for _, alias := range aliases {
		if _, ok := uniqueAliases[alias]; ok {
			panic(&NotUniqueAliasError{alias})
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
		panic(&NotFunctionFactoryError{factoValue.Kind()})
	}

	factoType := factoValue.Type()
	numIn := factoType.NumIn()
	numOut := factoType.NumOut()

	if numIn != len(inputs) {
		panic(&InvalidInputNumberError{numIn, len(inputs)})
	}

	if numOut != 1 {
		panic(&InvalidOutputNumberError{numOut})
	}

	outKind := factoType.Out(0).Kind()
	correctKind := outKind == reflect.Chan
	correctKind = correctKind || outKind == reflect.Func
	correctKind = correctKind || outKind == reflect.Interface
	correctKind = correctKind || outKind == reflect.Map
	correctKind = correctKind || outKind == reflect.Ptr
	correctKind = correctKind || outKind == reflect.Slice
	if !correctKind {
		panic(&InvalidOutputKindError{outKind})
	}

	return &Component{name, allAliases, factoValue, inputs}

}

func (self *Component) resolve(container *Container, name string, class reflect.Type) reflect.Value {

	instances, producers := container.getComponentInstances(name)

	if class.Kind() == reflect.Slice {

		class = class.Elem()

		slice := reflect.MakeSlice(reflect.SliceOf(class), 0, len(instances))
		slice = reflect.Append(slice, instances...)

		return slice

	} else if class.Kind() == reflect.Map {

		if keyClass := class.Key(); keyClass != STRING_TYPE {
			panic(&InvalidKeyTypeError{keyClass})
		}

		class = class.Elem()

		m := reflect.MakeMapWithSize(reflect.MapOf(STRING_TYPE, class), len(instances))
		for i := 0; i < len(instances); i++ {
			m.SetMapIndex(reflect.ValueOf(producers[i].name), instances[i])
		}

		return m

	} else {

		if len(instances) == 0 {
			panic(&NoProducerError{name})
		} else if len(instances) > 1 {
			panic(&TooManyProducersError{name, producers})
		}

		return instances[0]

	}

}

func (self *Component) Instanciate(container *Container) reflect.Value {

	defer func() {
		if err := recover(); err != nil {
			panic(&ComponentInstanciationError{self, err})
		}
	}()

	factoType := self.factory.Type()
	numIn := factoType.NumIn()

	args := make([]reflect.Value, numIn)

	for i := 0; i < numIn; i++ {
		args[i] = self.resolve(container, self.inputs[i], factoType.In(i))
	}

	instance := self.factory.Call(args)[0]
	if instance.Kind() == reflect.Interface {
		instance = instance.Elem()
	}

	return instance

}

func (self *Component) Initialize(container *Container, instance reflect.Value) {

	defer func() {
		if err := recover(); err != nil {
			panic(&ComponentInitialisationError{self, err})
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

	class := instance.Type()

	for i := 0; i < instance.NumField(); i++ {

		field := instance.Field(i)
		fieldType := class.Field(i)

		name, ok := fieldType.Tag.Lookup("inject")
		if !ok {
			continue
		}

		if name == "" {
			name = fieldType.Name
		}

		field.Set(self.resolve(container, name, field.Type()))

	}

}

func (self *Component) PostInit(container *Container, instance reflect.Value) {

	if postInitAwarable, ok := instance.Interface().(PostInitAwarable); ok {
		postInitAwarable.PostInit()
	}

}

func (self *Component) Close(container *Container, instance reflect.Value) {

	if closeable, ok := instance.Interface().(Closeable); ok {
		closeable.Close()
	}

}

func (self *Component) String() string {
	return self.name
}

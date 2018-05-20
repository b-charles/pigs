package ioc

import (
	"fmt"
	"reflect"

	"github.com/pkg/errors"
)

type Container struct {
	coreComponents map[string][]*Component
	testComponents map[string][]*Component
	instances      map[*Component]reflect.Value
	sorted         []*Component
}

func NewContainer() *Container {

	container := &Container{
		make(map[string][]*Component),
		make(map[string][]*Component),
		make(map[*Component]reflect.Value),
		make([]*Component, 0)}

	container.Put(container, "ApplicationContainer")

	return container

}

// PUT

func (self *Container) putIn(
	factory interface{},
	name string,
	aliases []string,
	components map[string][]*Component) {

	component := NewComponent(factory, name, aliases)

	for _, alias := range component.aliases {

		list, ok := components[alias]
		if !ok {
			list = make([]*Component, 0)
		}

		components[alias] = append(list, component)

	}

}

func (self *Container) wrap(object interface{}) func() interface{} {
	return func() interface{} {
		return object
	}
}

func (self *Container) PutFactory(factory interface{}, name string, aliases ...string) {
	self.putIn(factory, name, aliases, self.coreComponents)
}

func (self *Container) Put(object interface{}, name string, aliases ...string) {
	self.PutFactory(self.wrap(object), name, aliases...)
}

func (self *Container) TestPutFactory(factory interface{}, name string, aliases ...string) {
	self.putIn(factory, name, aliases, self.testComponents)
}

func (self *Container) TestPut(object interface{}, name string, aliases ...string) {
	self.TestPutFactory(self.wrap(object), name, aliases...)
}

// INTERNAL RESOLUTION

func (self *Container) getComponentInstance(component *Component) reflect.Value {

	instance, ok := self.instances[component]
	if ok {
		return instance
	}

	instance = component.Instanciate(self)
	self.instances[component] = instance

	component.Initialize(self, instance)
	component.PostInit(self, instance)

	self.sorted = append(self.sorted, component)

	return instance

}

func (self *Container) extractInstances(name string, components map[string][]*Component) ([]reflect.Value, []*Component) {

	list, ok := components[name]
	if !ok {
		return []reflect.Value{}, []*Component{}
	}

	instances := make([]reflect.Value, 0, len(list))
	producers := make([]*Component, 0, len(list))

	for _, component := range list {
		instance := self.getComponentInstance(component)
		if !instance.IsNil() {
			instances = append(instances, instance)
			producers = append(producers, component)
		}
	}

	return instances, producers

}

func (self *Container) getComponentInstances(name string) ([]reflect.Value, []*Component) {

	defer func() {
		if err, ok := recover().(error); ok {
			panic(errors.Wrapf(err, "Error during instanciations of '%s'", name))
		}
	}()

	instances, producers := self.extractInstances(name, self.testComponents)
	if len(instances) > 0 {
		return instances, producers
	}

	return self.extractInstances(name, self.coreComponents)

}

// EXTERNAL RESOLUTION

func (self *Container) convertToInterface(values []reflect.Value) []interface{} {

	instances := make([]interface{}, 0, len(values))
	for _, value := range values {
		instances = append(instances, value.Interface())
	}

	return instances

}

func (self *Container) GetComponents(name string) []interface{} {
	values, _ := self.getComponentInstances(name)
	return self.convertToInterface(values)
}

func (self *Container) GetComponent(name string) interface{} {

	instances, producers := self.getComponentInstances(name)

	if len(instances) == 0 {
		panic(fmt.Errorf("No producer found for '%s'.", name))
	} else if len(instances) > 1 {
		panic(fmt.Errorf("Too many producers found for '%s': %v.", name, producers))
	}

	return self.convertToInterface(instances)[0]

}

// CLOSE

func (self *Container) Close() {

	for _, component := range self.sorted {

		component.Close(self, self.instances[component])

		delete(self.instances, component)

	}

	self.sorted = nil

}

func (self *Container) ClearTests() {

	self.Close()

	self.testComponents = make(map[string][]*Component)

}

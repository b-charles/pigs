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

func (self *Container) wrap(object interface{}) func() interface{} {
	return func() interface{} {
		return object
	}
}

func (self *Container) PutFactory(factory interface{}, name string, aliases ...string) error {
	return self.putIn(factory, name, aliases, self.coreComponents)
}

func (self *Container) Put(object interface{}, name string, aliases ...string) error {
	return self.PutFactory(self.wrap(object), name, aliases...)
}

func (self *Container) TestPutFactory(factory interface{}, name string, aliases ...string) error {
	return self.putIn(factory, name, aliases, self.testComponents)
}

func (self *Container) TestPut(object interface{}, name string, aliases ...string) error {
	return self.TestPutFactory(self.wrap(object), name, aliases...)
}

// INTERNAL RESOLUTION

func (self *Container) getComponentInstance(component *Component) (reflect.Value, error) {

	if instance, ok := self.instances[component]; ok {
		return instance, nil
	}

	instance, err := component.Instanciate(self)
	if err != nil {
		return reflect.Value{}, err
	}

	self.instances[component] = instance

	if err := component.Initialize(self, instance); err != nil {
		return reflect.Value{}, err
	}
	if err := component.PostInit(self, instance); err != nil {
		return reflect.Value{}, err
	}

	self.sorted = append(self.sorted, component)

	return instance, nil

}

func (self *Container) extractInstances(name string, components map[string][]*Component) ([]reflect.Value, []*Component, error) {

	list, ok := components[name]
	if !ok {
		return []reflect.Value{}, []*Component{}, nil
	}

	instances := make([]reflect.Value, 0, len(list))
	producers := make([]*Component, 0, len(list))

	for _, component := range list {

		instance, err := self.getComponentInstance(component)

		if err != nil {
			return instances, producers, err
		}

		if !instance.IsNil() {
			instances = append(instances, instance)
			producers = append(producers, component)
		}

	}

	return instances, producers, nil

}

func (self *Container) getComponentInstances(name string) ([]reflect.Value, []*Component, error) {

	wrap := func(err error) error {
		if err != nil {
			return errors.Wrapf(err, "Error during instanciations of '%s'", name)
		} else {
			return nil
		}
	}

	if instances, producers, err := self.extractInstances(name, self.testComponents); err != nil || len(instances) > 0 {
		return instances, producers, wrap(err)
	}

	instances, producers, err := self.extractInstances(name, self.coreComponents)
	return instances, producers, wrap(err)

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
	if values, _, err := self.getComponentInstances(name); err != nil {
		panic(err)
	} else {
		return self.convertToInterface(values)
	}
}

func (self *Container) GetComponent(name string) interface{} {

	if instances, producers, err := self.getComponentInstances(name); err != nil {
		panic(err)
	} else if len(instances) == 0 {
		panic(fmt.Errorf("No producer found for '%s'.", name))
	} else if len(instances) > 1 {
		panic(fmt.Errorf("Too many producers found for '%s': %v.", name, producers))
	} else {
		return self.convertToInterface(instances)[0]
	}

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

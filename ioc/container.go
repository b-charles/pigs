// Package ioc provides an ioc framework.
package ioc

import (
	"fmt"
	"reflect"
	"time"
)

type context uint

const (
	core context = iota
	test
)

// A Container is a set of components. It manages the lifecycle of each
// component and take in charge the injection process.
type Container struct {
	coreComponents map[reflect.Type][]*component
	testComponents map[reflect.Type][]*component
	creationTime   time.Time
	instances      map[*component]*instance
	closables      []*instance
}

// NewContainer creates a new Container.
func NewContainer() *Container {

	container := &Container{
		coreComponents: make(map[reflect.Type][]*component),
		testComponents: make(map[reflect.Type][]*component),
		creationTime:   time.Now(),
		instances:      map[*component]*instance{},
		closables:      []*instance{}}

	container.Put(container)

	return container

}

// CreationTime returns the creation time of the container.
func (self *Container) CreationTime() time.Time {
	return self.creationTime
}

// REGISTRATION

// putIn creates a new component by its factory, signature functions
// (sigWrappers) and in the given context.
func (self *Container) putIn(factory reflect.Value, sigWrappers []any, context context) error {

	components := self.getComponentMap(context)

	comp, err := newComponent(self, factory, sigWrappers)
	if err != nil {
		return err
	}

	for _, sign := range comp.signatures {

		list, ok := components[sign]
		if !ok {
			list = make([]*component, 0)
		}

		components[sign] = append(list, comp)

	}

	return nil

}

// wrap converts a component into a factory which returns this component.
func wrap(component any) reflect.Value {

	value := reflect.ValueOf(component)

	out := []reflect.Type{reflect.TypeOf(component), error_type}

	wTyp := reflect.FuncOf([]reflect.Type{}, out, false)
	wVal := reflect.MakeFunc(wTyp, func(args []reflect.Value) []reflect.Value {
		return []reflect.Value{value, reflect.Zero(error_type)}
	})

	return wVal

}

// PutFactory records a component defined by a factory and optional
// signatures.
func (self *Container) PutFactory(factory any, signFuncs ...any) error {
	return self.putIn(reflect.ValueOf(factory), signFuncs, core)
}

// Put records directly a component with optional signatures.
func (self *Container) Put(component any, signFuncs ...any) error {
	return self.putIn(wrap(component), signFuncs, core)
}

// TestPutFactory records a test component defined by a factory and
// optional signatures.
func (self *Container) TestPutFactory(factory any, signFuncs ...any) error {
	return self.putIn(reflect.ValueOf(factory), signFuncs, test)
}

// TestPut records directly a test component with optional signatures.
func (self *Container) TestPut(component any, signFuncs ...any) error {
	return self.putIn(wrap(component), signFuncs, test)
}

// INTERNAL RESOLUTION

// instanciates gets the instance of the given component. If the instance is
// already created, the instance is returned. If not, the instance is created,
// recorded, initialized, post-initialized and returned.
func (self *Container) instanciates(component *component) (*instance, error) {

	if instance, ok := self.instances[component]; ok {
		return instance, nil
	}

	instance, err := component.instanciate()
	if err != nil {
		return nil, err
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

// getComponentMap gets the map corresponding to the given context.
func (self *Container) getComponentMap(context context) map[reflect.Type][]*component {
	if context == core {
		return self.coreComponents
	} else if context == test {
		return self.testComponents
	} else {
		panic(fmt.Errorf("Unknown context %v", context))
	}
}

// getInstancesOfContext returns the instances for a target typ and contained
// in a given context.
func (self *Container) getInstancesOfContext(typ reflect.Type, context context) ([]*instance, error) {

	components := self.getComponentMap(context)

	if list, ok := components[typ]; !ok {

		return []*instance{}, nil

	} else {

		instances := make([]*instance, 0, len(list))

		for _, component := range list {

			if instance, err := self.instanciates(component); err != nil {
				return nil, err
			} else {
				instances = append(instances, instance)
			}

		}

		return instances, nil

	}

}

// getInstances get all instances for a target typ. It searchs in the test and
// (if nothing is found) core context.
func (self *Container) getInstances(typ reflect.Type) ([]*instance, error) {

	if instances, err := self.getInstancesOfContext(typ, test); err != nil {
		return nil, err
	} else if len(instances) > 0 {
		return instances, nil
	}

	if instances, err := self.getInstancesOfContext(typ, core); err != nil {
		return nil, err
	} else {
		return instances, nil
	}

}

// getValue returns an injectable value for the given target.
func (self *Container) getValue(target reflect.Type) (reflect.Value, error) {

	instances, err := self.getInstances(target)
	if err != nil {
		return reflect.Zero(target), err
	}

	if len(instances) == 1 && instances[0].value.CanConvert(target) {

		return instances[0].value.Convert(target), nil

	} else if target.Kind() == reflect.Slice {

		elemTarget := target.Elem()

		instances, err = self.getInstances(elemTarget)

		slice := reflect.MakeSlice(target, 0, len(instances))
		if len(instances) == 0 {
			return slice, nil
		}

		for _, instance := range instances {

			v := instance.value
			if !v.CanConvert(elemTarget) {
				return reflect.Zero(target), fmt.Errorf("Can not convert '%v' to %v.", instance, elemTarget)
			}

			slice = reflect.Append(slice, v.Convert(elemTarget))

		}

		return slice, nil

	} else if len(instances) > 1 {

		producers := make([]*component, 0)
		for _, instance := range instances {
			producers = append(producers, instance.producer)
		}

		return reflect.Zero(target), fmt.Errorf("Too many components found: %v.", producers)

	} else if len(instances) == 0 {

		return reflect.Zero(target), fmt.Errorf("No component found for type '%v'.", target)

	} else {

		return reflect.Zero(target), fmt.Errorf("Can not convert '%v' to %v.", instances[0], target)

	}

}

// INJECTION

// getArguments returns initialized and injected arguments to call the given method.
func (self *Container) getArguments(method reflect.Value) ([]reflect.Value, error) {

	if method.Kind() != reflect.Func {
		return nil, fmt.Errorf("Can not use '%v' as a function.", method)
	}

	methodType := method.Type()

	nin := methodType.NumIn()
	args := make([]reflect.Value, nin, nin)

	for i := 0; i < nin; i++ {

		argType := methodType.In(i)

		if argValue, err := self.getValue(argType); err != nil {
			return nil, fmt.Errorf("Can not inject parameter #%d (type %v): %w", i, argType, err)
		} else {
			args[i] = argValue
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

// CallInjected call the given method, injecting its arguments.
func (self *Container) CallInjected(method any) error {

	// input checks

	methodValue := reflect.ValueOf(method)

	typ := methodValue.Type()
	if typ.Kind() != reflect.Func {
		return fmt.Errorf("The input should a function, not a %v.", typ)
	}
	nout := typ.NumOut()
	if nout > 1 {
		return fmt.Errorf("The method should return none or one output, not %d.", nout)
	} else if nout == 1 {
		if outType := typ.Out(0); !outType.AssignableTo(error_type) {
			return fmt.Errorf("The output of the method should be an error, not a '%v'.", outType)
		}
	}

	// get arguments

	args, err := self.getArguments(methodValue)

	defer func() {
		for c := len(self.closables) - 1; c >= 0; c-- {
			self.closables[c].Close()
		}
		self.closables = []*instance{}
	}()

	if err != nil {
		return err
	}

	// release unused instances

	self.instances = make(map[*component]*instance)

	if len(self.testComponents) == 0 {
		self.coreComponents = map[reflect.Type][]*component{}
	} else {
		self.testComponents = map[reflect.Type][]*component{}
	}

	// calling

	outs := methodValue.Call(args)

	// output

	if nout == 1 {
		return outs[0].Interface().(error)
	} else {
		return nil
	}

}

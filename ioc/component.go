package ioc

import (
	"fmt"
	"reflect"
)

type component struct {
	container *Container
	name      string
	aliases   []string
	factory   reflect.Value
}

func newComponent(container *Container, factory any, name string, aliases []string) (*component, error) {

	// check aliases uniqueness

	uniqueAliases := make(map[string]bool)
	uniqueAliases[name] = true
	for _, alias := range aliases {
		if _, ok := uniqueAliases[alias]; ok {
			return nil, fmt.Errorf("Alias specified more than once: '%s'.", alias)
		}
		uniqueAliases[alias] = true
	}

	// check factory signature

	factoryValue := reflect.ValueOf(factory)
	factoryType := reflect.TypeOf(factory)

	if factoryType.Kind() != reflect.Func {
		return nil, fmt.Errorf("The factory should be a function, but was %v.", factoryType)
	}

	if nout := factoryType.NumOut(); nout == 0 {
		return nil, fmt.Errorf("The factory should return at least one output.")
	} else if nout == 2 && !factoryType.Out(1).AssignableTo(error_type) {
		return nil, fmt.Errorf("The second output of the factory should be an %v, not '%v'.", error_type, factoryType.Out(1))
	} else if nout > 2 {
		return nil, fmt.Errorf("The factory should return one or two output, not %d.", nout)
	}

	// return
	return &component{container, name, aliases, factoryValue}, nil

}

func (self *component) instanciate() (*instance, error) {

	if self == nil {
		return nil, nil
	}

	outs, err := self.container.callInjected(self.factory)
	if err != nil {
		return nil, fmt.Errorf("Error during call of factory of '%s': %w", self.name, err)
	} else if len(outs) == 2 {
		if err := outs[1].Interface().(error); err != nil {
			return nil, fmt.Errorf("Error during instanciation of '%s': %w", self.name, err)
		}
	}

	obj := outs[0]
	if obj.Kind() == reflect.Interface {
		obj = obj.Elem()
	}

	return &instance{self, obj}, nil

}

func (self *component) String() string {
	if self == nil {
		return "nil"
	} else {
		return self.name
	}
}

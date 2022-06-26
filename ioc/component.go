package ioc

import (
	"fmt"
	"reflect"
)

type Component struct {
	name    string
	aliases []string
	factory reflect.Value
}

func NewComponent(factory interface{}, name string, aliases []string) (*Component, error) {

	wrap := func(err error) error {
		return fmt.Errorf("Error during registration of '%s': %w", name, err)
	}

	// get unique aliases

	uniqueAliases := make(map[string]bool)
	uniqueAliases[name] = true
	for _, alias := range aliases {
		if _, ok := uniqueAliases[alias]; ok {
			return nil, wrap(fmt.Errorf("Alias specified more than once: '%s'.", alias))
		}
		uniqueAliases[alias] = true
	}

	allAliases := make([]string, 0, len(uniqueAliases))
	for alias := range uniqueAliases {
		allAliases = append(allAliases, alias)
	}

	// return

	return &Component{name, allAliases, reflect.ValueOf(factory)}, nil

}

func (self *Component) instanciate(container *Container) (*Instance, error) {

	wrap := func(err error) (*Instance, error) {
		return voidInstance(self), fmt.Errorf("Error during instanciation of '%s': %w", self.name, err)
	}

	out, err := container.callInjected(self.factory)
	if err != nil {
		return wrap(err)
	}

	if len(out) == 0 {
		return wrap(fmt.Errorf("The factory should return at least one output."))
	}

	if len(out) == 2 {
		if err, ok := out[1].Interface().(error); !ok {
			return wrap(fmt.Errorf("The second output of the factory should be an error, not a '%v'.", out[1].Type()))
		} else if err != nil {
			return wrap(err)
		}
	}

	if len(out) > 2 {
		return wrap(fmt.Errorf("The factory should return one or two output, not %d.", len(out)))
	}

	obj := out[0]
	if obj.Kind() == reflect.Interface {
		obj = obj.Elem()
	}

	return newInstance(obj, self), nil

}

func (self *Component) String() string {
	return self.name
}

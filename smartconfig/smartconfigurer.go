package smartconfig

import (
	"fmt"
	"reflect"

	"github.com/b-charles/pigs/ioc"
	"github.com/b-charles/pigs/json"
	"github.com/b-charles/pigs/memfun"
)

type configurer struct {
	target reflect.Type
	setter func(NavConfig, reflect.Value) error
}

type SmartConfigurer struct {
	config      NavConfig
	configurers memfun.MemFun[reflect.Type, *configurer]
}

func findConfigurer(target reflect.Type, recfun func(reflect.Type) (*configurer, error)) (*configurer, error) {

	if target.Kind() == reflect.Pointer {
		return newPointerConfigurer(target, recfun)
	}

	if target.Kind() == reflect.Struct {
		return newStructConfigurer(target, recfun)
	}

	if target.Kind() == reflect.Slice {
		return newSliceConfigurer(target, recfun)
	}

	if target.Kind() == reflect.Map && target.Key() == string_type {
		return newMapConfigurer(target, recfun)
	}

	return nil, fmt.Errorf("No configurer found for type %v.", target)

}

func newSmartConfigurer(config NavConfig, parsers []Parser, inspectors []Inspector) (*SmartConfigurer, error) {

	configurers := memfun.NewMemFun(findConfigurer)

	for _, parser := range parsers {

		if configurer, err := parserConfigurer(parser); err != nil {
			return nil, fmt.Errorf("Error during recording the parser '%v': %w", parser, err)
		} else {

			if old, pres, err := configurers.Lookup(configurer.target); err != nil {
				return nil, err
			} else if pres {
				return nil, fmt.Errorf(
					"Two configurers '%v' and '%v' are defined to the same target %v.",
					old, parser, configurer.target)
			}

			configurers.Store(configurer.target, configurer)

		}

	}

	for _, inspector := range inspectors {

		if configurer, err := inspectorConfigurer(inspector); err != nil {
			return nil, fmt.Errorf("Error during recording the inspector '%v': %w", inspector, err)
		} else {

			if old, pres, err := configurers.Lookup(configurer.target); err != nil {
				return nil, err
			} else if pres {
				return nil, fmt.Errorf("Two configurers '%v' and '%v' are defined to the same target %v.", old, inspector, configurer.target)
			}

			configurers.Store(configurer.target, configurer)

		}

	}

	return &SmartConfigurer{
		config:      config,
		configurers: configurers,
	}, nil

}

func (self *SmartConfigurer) Configure(root string, configurable any) error {

	value := reflect.ValueOf(configurable)
	if value.Kind() == reflect.Pointer {
		value = value.Elem()
	}
	if !value.CanSet() {
		return fmt.Errorf("The value '%v' is not settable.", configurable)
	}

	if configurer, err := self.configurers.Get(value.Type()); err != nil {
		return err
	} else if err := configurer.setter(self.config.Get(root), value); err != nil {
		return err
	}

	return nil

}

func (self *SmartConfigurer) Json() json.JsonNode {
	types := self.configurers.Keys()
	return json.ReflectTypeSliceToJson(types)
}

func (self *SmartConfigurer) String() string {
	return self.Json().String()
}

func init() {
	ioc.PutNamedFactory("Smart Configuration", newSmartConfigurer)
}

func createConfig(root string, configurable any) any {

	factoryType := reflect.FuncOf(
		[]reflect.Type{smartConfigurer_type},
		[]reflect.Type{reflect.TypeOf(configurable), error_type},
		false)

	factory := reflect.MakeFunc(factoryType, func(args []reflect.Value) []reflect.Value {

		smartConfigurer := args[0].Interface().(*SmartConfigurer)
		if err := smartConfigurer.Configure(root, configurable); err != nil {
			return []reflect.Value{reflect.ValueOf(configurable), reflect.ValueOf(err)}
		} else {
			return []reflect.Value{reflect.ValueOf(configurable), reflect.Zero(error_type)}
		}

	})

	return factory.Interface()

}

func DefaultConfigure(root string, configurable any) {
	ioc.DefaultPutFactory(createConfig(root, configurable))
}

func DefaultConfigureNamed(name string, root string, configurable any) {
	ioc.DefaultPutNamedFactory(name, createConfig(root, configurable))
}

func Configure(root string, configurable any) {
	ioc.PutFactory(createConfig(root, configurable))
}

func ConfigureNamed(name string, root string, configurable any) {
	ioc.PutNamedFactory(name, createConfig(root, configurable))
}

func TestConfigure(root string, configurable any) {
	ioc.TestPutFactory(createConfig(root, configurable))
}

func TestConfigureNamed(name string, root string, configurable any) {
	ioc.TestPutFactory(name, createConfig(root, configurable))
}

package smartconfig

import (
	"fmt"
	"reflect"

	"github.com/b-charles/pigs/ioc"
)

type SmartConfigurer struct {
	config        NavConfig
	configurerMap map[reflect.Type]Configurer
}

func newSmartConfigurer(config NavConfig, parsers []Parser, configurers []Configurer) (*SmartConfigurer, error) {

	configurerMap := map[reflect.Type]Configurer{}

	for _, parser := range parsers {

		smart, err := smartify(parser)
		if err != nil {
			return nil, fmt.Errorf("Error during smartify the parser '%v': %w", parser, err)
		}

		target := smart.Target()

		if old, pres := configurerMap[target]; pres {
			return nil, fmt.Errorf("Two configurers '%v' and '%v' are defined to the same target %v.", old, parser, target)
		}

		configurerMap[target] = smart

	}

	for _, configurer := range configurers {

		target := configurer.Target()

		if old, pres := configurerMap[target]; pres {
			return nil, fmt.Errorf("Two configurers '%v' and '%v' are defined to the same target %v.", old, configurer, target)
		}

		configurerMap[target] = configurer

	}

	return &SmartConfigurer{
		config:        config,
		configurerMap: configurerMap,
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

	if configurer, err := self.FindConfigurer(value.Type()); err != nil {
		return err
	} else if err := configurer.Configure(self.config.Get(root), value); err != nil {
		return err
	}

	return nil

}

func (self *SmartConfigurer) FindConfigurer(target reflect.Type) (Configurer, error) {

	if configurer, present := self.configurerMap[target]; present {
		return configurer, nil
	}

	if target.Kind() == reflect.Pointer {
		configurer := newPointerConfigurer(target)
		self.configurerMap[target] = configurer
		configurer.analyze(self)
		return configurer, nil
	}

	if target.Kind() == reflect.Struct {
		configurer := newStructConfigurer(target)
		self.configurerMap[target] = configurer
		configurer.analyze(self)
		return configurer, nil
	}

	if target.Kind() == reflect.Slice {
		configurer := newSliceConfigurer(target)
		self.configurerMap[target] = configurer
		configurer.analyze(self)
		return configurer, nil
	}

	if target.Kind() == reflect.Map && target.Key() == string_type {
		configurer := newMapConfigurer(target)
		self.configurerMap[target] = configurer
		configurer.analyze(self)
		return configurer, nil
	}

	return nil, fmt.Errorf("No configurer found for type %v.", target)

}

func init() {
	ioc.PutFactory(newSmartConfigurer)
}

var smartConfigurer_type = reflect.TypeOf(func(*SmartConfigurer) {}).In(0)

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

func Configure(root string, configurable any) {
	ioc.PutFactory(createConfig(root, configurable))
}

func TestConfigure(root string, configurable any) {
	ioc.TestPutFactory(createConfig(root, configurable))
}

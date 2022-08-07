package smartconf

import (
	"fmt"
	"reflect"

	"github.com/b-charles/pigs/ioc"
)

type SmartConfigurer struct {
	config        NavConfig
	configurerMap map[reflect.Type]Configurer
}

func NewSmartConfigurer(config NavConfig, parsers []Parser, configurers []Configurer) (*SmartConfigurer, error) {

	configurerMap := map[reflect.Type]Configurer{}

	for _, parser := range parsers {

		smart, err := smartify(parser)
		if err != nil {
			return nil, fmt.Errorf("Error during smartify the parser '%v': %w", parser, err)
		}

		target := smart.Target()

		if old, pres := configurerMap[target]; pres {
			return nil, fmt.Errorf("Two configurers '%v' and '%v' are defined to the same target %v.",
				old, parser, target)
		}

		configurerMap[target] = smart

	}

	for _, configurer := range configurers {

		target := configurer.Target()

		if old, pres := configurerMap[target]; pres {
			return nil, fmt.Errorf("Two configurers '%v' and '%v' are defined to the same target %v.",
				old, configurer, target)
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

func pointerConfigurer(configurer Configurer) Configurer {

	typ := configurer.Target()
	target := reflect.PointerTo(typ)

	conf := func(config NavConfig, value reflect.Value) error {
		return configurer.Configure(config, value.Elem())
	}

	return &simpleConfigurer{target, conf}

}

func (self *SmartConfigurer) FindConfigurer(target reflect.Type) (Configurer, error) {

	if configurer, present := self.configurerMap[target]; present {
		return configurer, nil
	}

	elem := target
	if target.Kind() == reflect.Pointer {
		elem = target.Elem()
	}

	if elem != target {
		if configurer, present := self.configurerMap[elem]; present {
			wrapped := pointerConfigurer(configurer)
			self.configurerMap[target] = wrapped
			return wrapped, nil
		}
	}

	if target.Kind() == reflect.Struct {
		configurer := newStructConfigurer(target)
		self.configurerMap[target] = configurer
		configurer.analyze(self)
		return configurer, nil
	} else if elem != target && elem.Kind() == reflect.Struct {
		configurer := newStructConfigurer(elem)
		self.configurerMap[target] = pointerConfigurer(configurer)
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

	for typ, configurer := range self.configurerMap {
		if typ.AssignableTo(target) {
			self.configurerMap[target] = configurer
			return configurer, nil
		}
		if elem != target && typ.AssignableTo(elem) {
			wrapped := pointerConfigurer(configurer)
			self.configurerMap[target] = wrapped
			return wrapped, nil
		}
	}

	return nil, fmt.Errorf("No configurer found for type %v.", target)

}

func init() {
	ioc.PutFactory(NewSmartConfigurer)
}

var smartConfigurer_type = reflect.TypeOf(func(*SmartConfigurer) {}).In(0)

func createConfig(root string, configurable any) any {

	factoryType := reflect.FuncOf(
		[]reflect.Type{smartConfigurer_type},
		[]reflect.Type{reflect.TypeOf(configurable)},
		false)

	factory := reflect.MakeFunc(factoryType, func(args []reflect.Value) []reflect.Value {

		smartConfigurer := args[0].Interface().(*SmartConfigurer)
		smartConfigurer.Configure(root, configurable)

		return []reflect.Value{reflect.ValueOf(configurable)}

	})

	return factory.Interface()

}

func Configure(root string, configurable any) {
	ioc.PutFactory(createConfig(root, configurable))
}

func TestConfigure(root string, configurable any) {
	ioc.TestPutFactory(createConfig(root, configurable))
}

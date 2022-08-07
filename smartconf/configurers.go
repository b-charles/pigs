package smartconf

import (
	"fmt"
	"reflect"
	"strings"
)

// Struct configurer

type structConfigurer struct {
	target      reflect.Type
	configurers []fieldConfigurer
}

type fieldConfigurer struct {
	name       string
	key        string
	configurer Configurer
}

func newStructConfigurer(target reflect.Type) *structConfigurer {
	return &structConfigurer{
		target: target,
	}
}

func (self *structConfigurer) analyze(smartConfigurer *SmartConfigurer) error {

	if self.target.Kind() != reflect.Struct {
		return fmt.Errorf("The target '%v' is not a struct.", self.target)
	}

	nfields := self.target.NumField()
	self.configurers = make([]fieldConfigurer, nfields, nfields)

	for f := 0; f < nfields; f++ {

		field := self.target.Field(f)

		if !field.IsExported() {
			return fmt.Errorf("The field '%s' of %v is not exported.", field.Name, self.target)
		}

		key := field.Tag.Get("config")
		if key == "" {
			key = strings.ToLower(field.Name)
		}

		target := field.Type

		if configurer, err := smartConfigurer.FindConfigurer(target); err != nil {
			return fmt.Errorf("Can not configure field '%s' of %v: %w", field.Name, self.target, err)
		} else {
			self.configurers[f] = fieldConfigurer{field.Name, key, configurer}
		}

	}

	return nil

}

func (self *structConfigurer) Target() reflect.Type {
	return self.target
}

func (self *structConfigurer) Configure(config NavConfig, receiver reflect.Value) error {

	typ := receiver.Type()
	if typ != self.target {
		return fmt.Errorf("Unexpected value '%v' (%v): not a %v.", receiver, typ, self.target)
	}

	for f, fieldConfigurer := range self.configurers {

		field := receiver.Field(f)
		sub := config.Get(fieldConfigurer.key)

		if err := fieldConfigurer.configurer.Configure(sub, field); err != nil {
			return fmt.Errorf("Error during configuration of '%s' of %v (path: %s): %w",
				fieldConfigurer.name, self.target, sub.Path(), err)
		}

	}

	return nil

}

// Slice configurer

type sliceConfigurer struct {
	target     reflect.Type
	configurer Configurer
}

func newSliceConfigurer(target reflect.Type) *sliceConfigurer {
	return &sliceConfigurer{
		target: target,
	}
}

func (self *sliceConfigurer) analyze(smartConfigurer *SmartConfigurer) error {

	if self.target.Kind() != reflect.Slice {
		return fmt.Errorf("The target '%v' is not a slice.", self.target)
	}

	typ := self.target.Elem()

	if configurer, err := smartConfigurer.FindConfigurer(typ); err != nil {
		return fmt.Errorf("Can not configure slice of %v: %w", typ, err)
	} else {
		self.configurer = configurer
	}

	return nil

}

func (self *sliceConfigurer) Target() reflect.Type {
	return self.target
}

func (self *sliceConfigurer) Configure(config NavConfig, receiver reflect.Value) error {

	typ := receiver.Type()
	if typ != self.target {
		return fmt.Errorf("Unexpected value '%v' (%v): not a %v.", receiver, typ, self.target)
	}

	keys := config.Keys()
	value := reflect.MakeSlice(self.target, len(keys), len(keys))

	for k, key := range keys {

		v := value.Index(k)
		sub := config.Child(key)

		if err := self.configurer.Configure(sub, v); err != nil {
			return fmt.Errorf("Error during configuring element #%d of %v (config path: '%s'): %w",
				k, typ, sub.Path(), err)
		}

	}

	receiver.Set(value)

	return nil

}

// Map configurer

type mapConfigurer struct {
	target     reflect.Type
	elem       reflect.Type
	configurer Configurer
}

func newMapConfigurer(target reflect.Type) *mapConfigurer {
	return &mapConfigurer{
		target: target,
	}
}

var string_type = reflect.TypeOf("")

func (self *mapConfigurer) analyze(smartConfigurer *SmartConfigurer) error {

	if self.target.Kind() != reflect.Map {
		return fmt.Errorf("The target '%v' is not a map.", self.target)
	}
	if self.target.Key() != string_type {
		return fmt.Errorf("The key type of '%v' is not 'string'.", self.target)
	}

	self.elem = self.target.Elem()

	if configurer, err := smartConfigurer.FindConfigurer(self.elem); err != nil {
		return fmt.Errorf("Can not configure map of %v: %w", self.elem, err)
	} else {
		self.configurer = configurer
	}

	return nil

}

func (self *mapConfigurer) Target() reflect.Type {
	return self.target
}

func (self *mapConfigurer) Configure(config NavConfig, receiver reflect.Value) error {

	typ := receiver.Type()
	if typ != self.target {
		return fmt.Errorf("Unexpected value '%v' (%v): not a %v.", receiver, typ, self.target)
	}

	value := reflect.MakeMap(self.target)

	for _, key := range config.Keys() {

		sub := config.Child(key)

		v := reflect.New(self.elem)
		if err := self.configurer.Configure(sub, v.Elem()); err != nil {
			return fmt.Errorf("Error during configuring element '%s' of %v (config path: '%s'): %w",
				key, typ, sub.Path(), err)
		}

		value.SetMapIndex(reflect.ValueOf(key), v.Elem())

	}

	receiver.Set(value)

	return nil

}

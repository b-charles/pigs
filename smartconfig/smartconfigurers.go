package smartconfig

import (
	"fmt"
	"reflect"
	"strings"
)

// Pointer configurer

func newPointerConfigurer(smartConfigurer *SmartConfigurer, target reflect.Type) (*configurer, error) {

	if target.Kind() != reflect.Pointer {
		return nil, fmt.Errorf("The target '%v' is not a pointer.", target)
	}

	if subConfigurer, err := smartConfigurer.findConfigurer(target.Elem()); err != nil {
		return nil, fmt.Errorf("Can not configure to %v: %w", target, err)
	} else {

		return &configurer{
			source: fmt.Sprintf("<internal %v configurer>", target),
			target: target,
			setter: func(config NavConfig, receiver reflect.Value) error {

				typ := receiver.Type()
				if typ != target {
					return fmt.Errorf("Unexpected value '%v' (%v): not a %v.", receiver, typ, target)
				}

				value := reflect.New(target.Elem())
				if err := subConfigurer.setter(config, value.Elem()); err != nil {
					return err
				} else {
					receiver.Set(value)
					return nil
				}

			},
		}, nil

	}

}

// Struct configurer

type fieldConfigurer struct {
	name       string
	key        string
	configurer *configurer
}

func newStructConfigurer(smartConfigurer *SmartConfigurer, target reflect.Type) (*configurer, error) {

	if target.Kind() != reflect.Struct {
		return nil, fmt.Errorf("The target '%v' is not a struct.", target)
	}

	nfields := target.NumField()
	configurers := make([]fieldConfigurer, nfields, nfields)

	for f := 0; f < nfields; f++ {

		field := target.Field(f)

		if !field.IsExported() {
			return nil, fmt.Errorf("The field '%s' of %v is not exported.", field.Name, target)
		}

		key := field.Tag.Get("config")
		if key == "" {
			key = strings.ToLower(field.Name)
		}

		target := field.Type

		if configurer, err := smartConfigurer.findConfigurer(target); err != nil {
			return nil, fmt.Errorf("Can not configure field '%s' of %v: %w", field.Name, target, err)
		} else {
			configurers[f] = fieldConfigurer{field.Name, key, configurer}
		}

	}

	return &configurer{
		source: fmt.Sprintf("<internal %v configurer>", target),
		target: target,
		setter: func(config NavConfig, receiver reflect.Value) error {

			typ := receiver.Type()
			if typ != target {
				return fmt.Errorf("Unexpected value '%v' (%v): not a %v.", receiver, typ, target)
			}

			for f, fieldConfigurer := range configurers {

				field := receiver.Field(f)
				sub := config.Get(fieldConfigurer.key)

				if err := fieldConfigurer.configurer.setter(sub, field); err != nil {
					return fmt.Errorf("Error during configuration of '%s' of %v (path: %s): %w",
						fieldConfigurer.name, target, sub.Path(), err)
				}

			}

			return nil

		},
	}, nil

}

// Slice configurer

func newSliceConfigurer(smartConfigurer *SmartConfigurer, target reflect.Type) (*configurer, error) {

	if target.Kind() != reflect.Slice {
		return nil, fmt.Errorf("The target '%v' is not a slice.", target)
	}

	typ := target.Elem()

	if subConfigurer, err := smartConfigurer.findConfigurer(typ); err != nil {
		return nil, fmt.Errorf("Can not configure slice of %v: %w", typ, err)
	} else {

		return &configurer{
			source: fmt.Sprintf("<internal %v configurer>", target),
			target: target,
			setter: func(config NavConfig, receiver reflect.Value) error {

				typ := receiver.Type()
				if typ != target {
					return fmt.Errorf("Unexpected value '%v' (%v): not a %v.", receiver, typ, target)
				}

				keys := config.Keys()
				value := reflect.MakeSlice(target, len(keys), len(keys))

				for k, key := range keys {

					v := value.Index(k)
					sub := config.Child(key)

					if err := subConfigurer.setter(sub, v); err != nil {
						return fmt.Errorf("Error during configuring element #%d of %v (config path: '%s'): %w",
							k, typ, sub.Path(), err)
					}

				}

				receiver.Set(value)

				return nil

			},
		}, nil

	}

}

// Map configurer

func newMapConfigurer(smartConfigurer *SmartConfigurer, target reflect.Type) (*configurer, error) {

	if target.Kind() != reflect.Map {
		return nil, fmt.Errorf("The target '%v' is not a map.", target)
	}
	if target.Key() != string_type {
		return nil, fmt.Errorf("The key type of '%v' is not 'string'.", target)
	}

	elem := target.Elem()

	if subConfigurer, err := smartConfigurer.findConfigurer(elem); err != nil {
		return nil, fmt.Errorf("Can not configure map of %v: %w", elem, err)
	} else {

		return &configurer{
			source: fmt.Sprintf("<internal %v configurer>", target),
			target: target,
			setter: func(config NavConfig, receiver reflect.Value) error {

				typ := receiver.Type()
				if typ != target {
					return fmt.Errorf("Unexpected value '%v' (%v): not a %v.", receiver, typ, target)
				}

				value := reflect.MakeMap(target)

				for _, key := range config.Keys() {

					sub := config.Child(key)

					v := reflect.New(elem)
					if err := subConfigurer.setter(sub, v.Elem()); err != nil {
						return fmt.Errorf("Error during configuring element '%s' of %v (config path: '%s'): %w",
							key, typ, sub.Path(), err)
					}

					value.SetMapIndex(reflect.ValueOf(key), v.Elem())

				}

				receiver.Set(value)

				return nil

			},
		}, nil

	}

}

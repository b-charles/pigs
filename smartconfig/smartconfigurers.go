package smartconfig

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
)

// Pointer configurer

func newPointerConfigurer(target reflect.Type, recfun func(reflect.Type) (*configurer, error)) (*configurer, error) {

	if target.Kind() != reflect.Pointer {
		return nil, fmt.Errorf("The target '%v' is not a pointer.", target)
	}

	elementType := target.Elem()

	var (
		once          sync.Once
		subConfigurer *configurer
		subError      error
	)

	return &configurer{
		target: target,
		setter: func(config NavConfig, receiver reflect.Value) error {

			typ := receiver.Type()
			if typ != target {
				return fmt.Errorf("Unexpected value '%v' (%v): not a %v.", receiver, typ, target)
			}

			once.Do(func() {
				if configurer, err := recfun(elementType); err != nil {
					subError = fmt.Errorf("Can not configure to %v: %w", target, err)
				} else {
					subConfigurer = configurer
				}
			})
			if subError != nil {
				return subError
			}

			value := reflect.New(elementType)
			if err := subConfigurer.setter(config, value.Elem()); err != nil {
				return err
			} else {
				receiver.Set(value)
				return nil
			}

		},
	}, nil

}

// Struct configurer

func newStructConfigurer(target reflect.Type, recfun func(reflect.Type) (*configurer, error)) (*configurer, error) {

	if target.Kind() != reflect.Struct {
		return nil, fmt.Errorf("The target '%v' is not a struct.", target)
	}

	var (
		once     sync.Once
		subError error
	)

	nfields := target.NumField()
	configurers := make([]func(config NavConfig, receiver reflect.Value) error, nfields, nfields)

	return &configurer{
		target: target,
		setter: func(config NavConfig, receiver reflect.Value) error {

			once.Do(func() {

				for fieldNum := 0; fieldNum < nfields; fieldNum++ {
					f := fieldNum

					field := target.Field(f)

					if !field.IsExported() {
						subError = fmt.Errorf("The field '%s' of %v is not exported.", field.Name, target)
						return
					}

					key := field.Tag.Get("config")
					if key == "" {
						key = strings.ToLower(field.Name)
					}

					if configurer, e := recfun(field.Type); e != nil {

						subError = fmt.Errorf("Can not configure field '%s' of %v: %w",
							field.Name, target, e)
						return

					} else {

						configurers[f] = func(config NavConfig, receiver reflect.Value) error {

							receiverField := receiver.Field(f)
							sub := config.Get(key)

							if err := configurer.setter(sub, receiverField); err != nil {
								return fmt.Errorf("Error during configuration of '%s' of %v (path: %s): %w",
									field.Name, target, sub.Path(), err)
							} else {
								return nil
							}

						}

					}

				}

			})
			if subError != nil {
				return subError
			}

			typ := receiver.Type()
			if typ != target {
				return fmt.Errorf("Unexpected value '%v' (%v): not a %v.", receiver, typ, target)
			}

			for _, funconfig := range configurers {
				if err := funconfig(config, receiver); err != nil {
					return err
				}
			}

			return nil

		},
	}, nil

}

// Slice configurer

func newSliceConfigurer(target reflect.Type, recfun func(reflect.Type) (*configurer, error)) (*configurer, error) {

	if target.Kind() != reflect.Slice {
		return nil, fmt.Errorf("The target '%v' is not a slice.", target)
	}

	elementType := target.Elem()

	var (
		once          sync.Once
		subConfigurer *configurer
		subError      error
	)

	return &configurer{
		target: target,
		setter: func(config NavConfig, receiver reflect.Value) error {

			typ := receiver.Type()
			if typ != target {
				return fmt.Errorf("Unexpected value '%v' (%v): not a %v.", receiver, typ, target)
			}

			once.Do(func() {
				if configurer, err := recfun(elementType); err != nil {
					subError = fmt.Errorf("Can not configure slice of %v: %w", elementType, err)
				} else {
					subConfigurer = configurer
				}
			})
			if subError != nil {
				return subError
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

// Map configurer

func newMapConfigurer(target reflect.Type, recfun func(reflect.Type) (*configurer, error)) (*configurer, error) {

	if target.Kind() != reflect.Map {
		return nil, fmt.Errorf("The target '%v' is not a map.", target)
	}
	if target.Key() != string_type {
		return nil, fmt.Errorf("The key type of '%v' is not 'string'.", target)
	}

	elementType := target.Elem()

	var (
		once          sync.Once
		subConfigurer *configurer
		subError      error
	)

	return &configurer{
		target: target,
		setter: func(config NavConfig, receiver reflect.Value) error {

			typ := receiver.Type()
			if typ != target {
				return fmt.Errorf("Unexpected value '%v' (%v): not a %v.", receiver, typ, target)
			}

			once.Do(func() {
				if configurer, err := recfun(elementType); err != nil {
					subError = fmt.Errorf("Can not configure map of %v: %w", elementType, err)
				} else {
					subConfigurer = configurer
				}
			})
			if subError != nil {
				return subError
			}

			value := reflect.MakeMap(target)

			for _, key := range config.Keys() {

				sub := config.Child(key)

				v := reflect.New(elementType)
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

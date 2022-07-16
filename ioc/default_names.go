package ioc

import (
	"fmt"
	"reflect"
)

// DefaultComponentName returns the default name of a component type.
func DefaultComponentName(component any) string {
	return defaultComponentName(reflect.TypeOf(component))
}

// defaultComponentName returns the default name of a component type.
func defaultComponentName(typ reflect.Type) string {

	if typ.Kind() == reflect.Pointer || typ.Kind() == reflect.Slice ||
		(typ.Kind() == reflect.Map && typ.Key() == string_type) {
		typ = typ.Elem()
	}

	name := typ.Name()
	if name == "" {
		return typ.String()
	}

	pkg := typ.PkgPath()
	if pkg == "" {
		return name
	}

	return fmt.Sprintf("%s/%s", pkg, name)

}

// DefaultFactoryName returns the default name of a component type.
func DefaultFactoryName(factory any) string {
	if name, err := defaultFactoryName(reflect.TypeOf(factory)); err != nil {
		panic(err)
	} else {
		return name
	}
}

// defaultFactoryName returns the default name of a component type.
func defaultFactoryName(typ reflect.Type) (string, error) {

	if typ.Kind() != reflect.Func {
		return "", fmt.Errorf("The type %v should be a function.", typ)
	}

	o := typ.NumOut()
	if o == 0 {
		return "", fmt.Errorf("The function should return at least one value.")
	}

	return defaultComponentName(typ.Out(0)), nil

}

// defaultAlias returns the aliases of the given argument.
func defaultAlias(alias any) ([]string, error) {

	value := reflect.ValueOf(alias)
	typ := value.Type()

	if typ == string_type {
		return []string{value.String()}, nil
	}

	if typ.Kind() == reflect.Func {
		list := make([]string, 0)
		for i := 0; i < typ.NumIn(); i++ {
			list = append(list, defaultComponentName(typ.In(i)))
		}
		return list, nil
	}

	return nil, fmt.Errorf("Can not guess aliases of '%v'", alias)

}

// unsafeDefaultAlias returns the first alias, and panics in case of error.
func unsafeDefaultAlias(alias any) string {
	aliases, err := defaultAlias(alias)
	if err != nil {
		panic(err)
	}
	return aliases[0]
}

// DefaultAliases returns the aliases of given arguments.
func DefaultAliases(aliases ...any) []string {
	if list, err := defaultAliases(aliases...); err != nil {
		panic(err)
	} else {
		return list
	}
}

// defaultAliases returns the aliases of given arguments.
func defaultAliases(aliases ...any) ([]string, error) {

	list := make([]string, 0)

	for _, elt := range aliases {

		l, err := defaultAlias(elt)
		if err != nil {
			return nil, err
		}

		list = append(list, l...)

	}

	return list, nil

}

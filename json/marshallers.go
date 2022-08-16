package json

import (
	"fmt"
	"reflect"

	"github.com/b-charles/pigs/ioc"
	"github.com/b-charles/pigs/json/core"
)

// Defaults marshallers

func init() {

	ioc.Put(func(v string) (core.JsonNode, error) {
		return core.JsonString(v), nil
	}, func(JsonMarshaller) {})

	ioc.Put(func(v float64) (core.JsonNode, error) {
		return core.JsonFloat(v), nil
	}, func(JsonMarshaller) {})

	ioc.Put(func(v int) (core.JsonNode, error) {
		return core.JsonInt(v), nil
	}, func(JsonMarshaller) {})

	ioc.Put(func(v bool) (core.JsonNode, error) {
		return core.JsonBool(v), nil
	}, func(JsonMarshaller) {})

}

// Pointer marshaller

func newPointerMarshaller(mapper *JsonMapper, target reflect.Type, valueMarshaller *jsonValueMarshaller) error {

	if target.Kind() != reflect.Pointer {
		return fmt.Errorf("The target %v is not a pointer.", target)
	}

	if marshaller, err := mapper.getMarshaller(target.Elem()); err != nil {
		return fmt.Errorf("Can not marshal %v to json: %w", target, err)
	} else {
		valueMarshaller.f = func(v reflect.Value) (core.JsonNode, error) {
			if v.IsZero() {
				return core.JSON_NULL, nil
			} else {
				return marshaller.f(v.Elem())
			}
		}
		return nil
	}

}

// Struct marshaller

func newStructMarshaller(mapper *JsonMapper, target reflect.Type, valueMarshaller *jsonValueMarshaller) error {

	if target.Kind() != reflect.Struct {
		return fmt.Errorf("The target %v is not a struct.", target)
	}

	fieldMarshallers := make([]func(reflect.Value, *core.JsonBuilder) error, 0, target.NumField())
	for f := 0; f < target.NumField(); f++ {

		field := target.Field(f)
		if !field.IsExported() {
			continue
		}

		key := field.Tag.Get("json")
		if key == "" {
			key = field.Name
		}

		if marshaller, err := mapper.getMarshaller(field.Type); err != nil {
			return fmt.Errorf("Can not marshal field %v of %v to json: %w", field.Name, target, err)
		} else {
			numField := f
			fieldMarshallers = append(fieldMarshallers, func(v reflect.Value, b *core.JsonBuilder) error {
				if json, err := marshaller.f(v.Field(numField)); err != nil {
					return err
				} else {
					b.Set(key, json)
					return nil
				}
			})
		}

	}

	valueMarshaller.f = func(v reflect.Value) (core.JsonNode, error) {

		if v.IsZero() {
			return core.JSON_NULL, nil
		}

		b := core.NewJsonBuilder()

		for _, marshaller := range fieldMarshallers {
			if err := marshaller(v, b); err != nil {
				return core.JSON_EMPTY_OBJECT, err
			}
		}

		return b.Build(), nil

	}
	return nil

}

// Slice marshaller

func newSliceMarshaller(mapper *JsonMapper, target reflect.Type, valueMarshaller *jsonValueMarshaller) error {

	if target.Kind() != reflect.Slice {
		return fmt.Errorf("The target %v is not a slice.", target)
	}

	if marshaller, err := mapper.getMarshaller(target.Elem()); err != nil {
		return fmt.Errorf("Can not marshal %v to json: %w", target, err)
	} else {
		valueMarshaller.f = func(v reflect.Value) (core.JsonNode, error) {
			if v.IsZero() {
				return core.JSON_NULL, nil
			}
			elts := make([]core.JsonNode, v.Len())
			for i := 0; i < v.Len(); i++ {
				if elts[i], err = marshaller.f(v.Index(i)); err != nil {
					return core.JSON_EMPTY_ARRAY, err
				}
			}
			return core.NewJsonArray(elts), nil
		}
		return nil
	}

}

// Map marshaller

func newMapMarshaller(mapper *JsonMapper, target reflect.Type, valueMarshaller *jsonValueMarshaller) error {

	if target.Kind() != reflect.Map {
		return fmt.Errorf("The target %v is not a map.", target)
	} else if target.Key() != stringType {
		return fmt.Errorf("The target %v key type is not string.", target)
	}

	if marshaller, err := mapper.getMarshaller(target.Elem()); err != nil {
		return fmt.Errorf("Can not marshal %v to json: %w", target, err)
	} else {
		valueMarshaller.f = func(v reflect.Value) (core.JsonNode, error) {
			if v.IsZero() {
				return core.JSON_NULL, nil
			}
			elts := make(map[string]core.JsonNode)
			iter := v.MapRange()
			for iter.Next() {
				if e, err := marshaller.f(iter.Value()); err != nil {
					return core.JSON_EMPTY_OBJECT, err
				} else {
					elts[iter.Key().Interface().(string)] = e
				}
			}
			return core.NewJsonObject(elts), nil
		}
		return nil
	}

}

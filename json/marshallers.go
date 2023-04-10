package json

import (
	"fmt"
	"reflect"

	"github.com/b-charles/pigs/ioc"
)

// Defaults marshallers

type StringMarshaller func(v string) (JsonNode, error)
type Float64Marshaller func(v float64) (JsonNode, error)
type IntMarshaller func(v int) (JsonNode, error)
type BoolMarshaller func(v bool) (JsonNode, error)
type ErrorMarshaller func(v error) (JsonNode, error)

func init() {

	ioc.DefaultPut(func(v string) (JsonNode, error) {
		return JsonString(v), nil
	}, func(StringMarshaller) {})
	ioc.PutFactory(func(m StringMarshaller) (JsonMarshaller, error) { return m, nil })

	ioc.DefaultPut(func(v float64) (JsonNode, error) {
		return JsonFloat(v), nil
	}, func(Float64Marshaller) {})
	ioc.PutFactory(func(m Float64Marshaller) (JsonMarshaller, error) { return m, nil })

	ioc.DefaultPut(func(v int) (JsonNode, error) {
		return JsonInt(v), nil
	}, func(IntMarshaller) {})
	ioc.PutFactory(func(m IntMarshaller) (JsonMarshaller, error) { return m, nil })

	ioc.DefaultPut(func(v bool) (JsonNode, error) {
		return JsonBool(v), nil
	}, func(BoolMarshaller) {})
	ioc.PutFactory(func(m BoolMarshaller) (JsonMarshaller, error) { return m, nil })

	ioc.DefaultPut(func(v error) (JsonNode, error) {
		return JsonString(v.Error()), nil
	}, func(ErrorMarshaller) {})
	ioc.PutFactory(func(m ErrorMarshaller) (JsonMarshaller, error) { return m, nil })

}

// Marshallers of JsonNode implementations

func init() {

	ioc.Put(func(v JsonString) (JsonNode, error) {
		return v, nil
	}, func(JsonMarshaller) {})

	ioc.Put(func(v JsonFloat) (JsonNode, error) {
		return v, nil
	}, func(JsonMarshaller) {})

	ioc.Put(func(v JsonInt) (JsonNode, error) {
		return v, nil
	}, func(JsonMarshaller) {})

	ioc.Put(func(v JsonBool) (JsonNode, error) {
		return v, nil
	}, func(JsonMarshaller) {})

	ioc.Put(func(v *JsonObject) (JsonNode, error) {
		return v, nil
	}, func(JsonMarshaller) {})

	ioc.Put(func(v *JsonArray) (JsonNode, error) {
		return v, nil
	}, func(JsonMarshaller) {})

	ioc.Put(func(v JsonNull) (JsonNode, error) {
		return v, nil
	}, func(JsonMarshaller) {})

}

// Pointer marshaller

func newPointerMarshaller(mapper *JsonMapper, target reflect.Type, valueMarshaller *wrappedMarshaller) error {

	if target.Kind() != reflect.Pointer {
		return fmt.Errorf("The target %v is not a pointer.", target)
	}

	if marshaller, err := mapper.getMarshaller(target.Elem()); err != nil {
		return fmt.Errorf("Can not marshal %v to json: %w", target, err)
	} else {
		valueMarshaller.f = func(v reflect.Value) (JsonNode, error) {
			if v.IsZero() {
				return JSON_NULL, nil
			} else {
				return marshaller.f(v.Elem())
			}
		}
		return nil
	}

}

// Struct marshaller

func newStructMarshaller(mapper *JsonMapper, target reflect.Type, valueMarshaller *wrappedMarshaller) error {

	if target.Kind() != reflect.Struct {
		return fmt.Errorf("The target %v is not a struct.", target)
	}

	fieldMarshallers := make([]func(reflect.Value, *JsonBuilder) error, 0, target.NumField())
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
			fieldMarshallers = append(fieldMarshallers, func(v reflect.Value, b *JsonBuilder) error {
				if json, err := marshaller.f(v.Field(numField)); err != nil {
					return err
				} else {
					b.Set(key, json)
					return nil
				}
			})
		}

	}

	valueMarshaller.f = func(v reflect.Value) (JsonNode, error) {

		if v.IsZero() {
			return JSON_NULL, nil
		}

		b := NewJsonBuilder()

		for _, marshaller := range fieldMarshallers {
			if err := marshaller(v, b); err != nil {
				return JSON_EMPTY_OBJECT, err
			}
		}

		return b.Build(), nil

	}
	return nil

}

// Slice marshaller

func newSliceMarshaller(mapper *JsonMapper, target reflect.Type, valueMarshaller *wrappedMarshaller) error {

	if target.Kind() != reflect.Slice {
		return fmt.Errorf("The target %v is not a slice.", target)
	}

	if marshaller, err := mapper.getMarshaller(target.Elem()); err != nil {
		return fmt.Errorf("Can not marshal %v to json: %w", target, err)
	} else {
		valueMarshaller.f = func(v reflect.Value) (JsonNode, error) {
			if v.IsZero() {
				return JSON_NULL, nil
			}
			elts := make([]JsonNode, v.Len())
			for i := 0; i < v.Len(); i++ {
				if elts[i], err = marshaller.f(v.Index(i)); err != nil {
					return JSON_EMPTY_ARRAY, err
				}
			}
			return NewJsonArray(elts), nil
		}
		return nil
	}

}

// Map marshaller

func newMapMarshaller(mapper *JsonMapper, target reflect.Type, valueMarshaller *wrappedMarshaller) error {

	if target.Kind() != reflect.Map {
		return fmt.Errorf("The target %v is not a map.", target)
	} else if target.Key() != stringType {
		return fmt.Errorf("The target %v key type is not string.", target)
	}

	if marshaller, err := mapper.getMarshaller(target.Elem()); err != nil {
		return fmt.Errorf("Can not marshal %v to json: %w", target, err)
	} else {
		valueMarshaller.f = func(v reflect.Value) (JsonNode, error) {
			if v.IsZero() {
				return JSON_NULL, nil
			}
			elts := make(map[string]JsonNode)
			iter := v.MapRange()
			for iter.Next() {
				if e, err := marshaller.f(iter.Value()); err != nil {
					return JSON_EMPTY_OBJECT, err
				} else {
					elts[iter.Key().Interface().(string)] = e
				}
			}
			return NewJsonObject(elts), nil
		}
		return nil
	}

}

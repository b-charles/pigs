package json

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/b-charles/pigs/ioc"
)

// Defaults marshallers

type StringMarshaller func(v string) (JsonNode, error)
type Float64Marshaller func(v float64) (JsonNode, error)
type IntMarshaller func(v int) (JsonNode, error)
type BoolMarshaller func(v bool) (JsonNode, error)
type ErrorMarshaller func(v error) (JsonNode, error)

func init() {

	ioc.DefaultPutNamed("String Json marshaller (default)",
		func(v string) (JsonNode, error) {
			return JsonString(v), nil
		}, func(StringMarshaller) {})
	ioc.PutNamedFactory("String Json marshaller (promoter)",
		func(m StringMarshaller) (JsonMarshaller, error) { return m, nil })

	ioc.DefaultPutNamed("Float64 Json marshaller (default)",
		func(v float64) (JsonNode, error) {
			return JsonFloat(v), nil
		}, func(Float64Marshaller) {})
	ioc.PutNamedFactory("Float64 Json marshaller (promoter)",
		func(m Float64Marshaller) (JsonMarshaller, error) { return m, nil })

	ioc.DefaultPutNamed("Int Json marshaller (default)",
		func(v int) (JsonNode, error) {
			return JsonInt(v), nil
		}, func(IntMarshaller) {})
	ioc.PutNamedFactory("Int Json marshaller (promoter)",
		func(m IntMarshaller) (JsonMarshaller, error) { return m, nil })

	ioc.DefaultPutNamed("Bool Json marshaller (default)",
		func(v bool) (JsonNode, error) {
			return JsonBool(v), nil
		}, func(BoolMarshaller) {})
	ioc.PutNamedFactory("Bool Json marshaller (promoter)",
		func(m BoolMarshaller) (JsonMarshaller, error) { return m, nil })

	ioc.DefaultPutNamed("Error Json marshaller (default)",
		func(v error) (JsonNode, error) {
			return JsonString(v.Error()), nil
		}, func(ErrorMarshaller) {})
	ioc.PutNamedFactory("Error Json marshaller (promoter)",
		func(m ErrorMarshaller) (JsonMarshaller, error) { return m, nil })

}

// Marshallers of JsonNode implementations

func init() {

	ioc.PutNamed("JsonString marshaller",
		func(v JsonString) (JsonNode, error) {
			return v, nil
		}, func(JsonMarshaller) {})

	ioc.PutNamed("JsonFloat marshaller",
		func(v JsonFloat) (JsonNode, error) {
			return v, nil
		}, func(JsonMarshaller) {})

	ioc.PutNamed("JsonInt marshaller",
		func(v JsonInt) (JsonNode, error) {
			return v, nil
		}, func(JsonMarshaller) {})

	ioc.PutNamed("JsonBool marshaller",
		func(v JsonBool) (JsonNode, error) {
			return v, nil
		}, func(JsonMarshaller) {})

	ioc.PutNamed("JsonObject marshaller",
		func(v *JsonObject) (JsonNode, error) {
			return v, nil
		}, func(JsonMarshaller) {})

	ioc.PutNamed("JsonArray marshaller",
		func(v *JsonArray) (JsonNode, error) {
			return v, nil
		}, func(JsonMarshaller) {})

	ioc.PutNamed("JsonNull marshaller",
		func(v JsonNull) (JsonNode, error) {
			return v, nil
		}, func(JsonMarshaller) {})

}

// Jsoner marshaller

type Jsoner interface {
	Json() JsonNode
}

func newJsonerMarshaller(target reflect.Type) (*wrappedMarshaller, error) {

	if !target.Implements(jsonerType) {
		return nil, fmt.Errorf("The target %v doesn't implements the Jsoner interface.", target)
	}

	return &wrappedMarshaller{
		t: target,
		f: func(v reflect.Value) (JsonNode, error) {
			return v.Interface().(Jsoner).Json(), nil
		},
	}, nil

}

// Pointer marshaller

func newPointerMarshaller(
	target reflect.Type,
	recfun func(reflect.Type) (*wrappedMarshaller, error)) (*wrappedMarshaller, error) {

	if target.Kind() != reflect.Pointer {
		return nil, fmt.Errorf("The target %v is not a pointer.", target)
	}

	var (
		once           sync.Once
		marshaller     *wrappedMarshaller
		marshaller_err error
	)

	return &wrappedMarshaller{
		t: target,
		f: func(v reflect.Value) (JsonNode, error) {

			once.Do(func() {
				marshaller, marshaller_err = recfun(target.Elem())
			})
			if marshaller_err != nil {
				return nil, fmt.Errorf("Can not marshal %v to json: %w", target, marshaller_err)
			}

			if v.IsZero() {
				return JSON_NULL, nil
			} else {
				return marshaller.f(v.Elem())
			}

		},
	}, nil

}

// Struct marshaller

func newStructMarshaller(
	target reflect.Type,
	recfun func(reflect.Type) (*wrappedMarshaller, error)) (*wrappedMarshaller, error) {

	if target.Kind() != reflect.Struct {
		return nil, fmt.Errorf("The target %v is not a struct.", target)
	}

	var (
		once           sync.Once
		marshaller_err error
	)

	nfields := target.NumField()
	fieldMarshallers := make([]func(reflect.Value, *JsonBuilder) error, 0, nfields)

	return &wrappedMarshaller{
		t: target,
		f: func(v reflect.Value) (JsonNode, error) {

			once.Do(func() {

				for fieldNum := 0; fieldNum < nfields; fieldNum++ {
          f := fieldNum

					field := target.Field(f)
					if !field.IsExported() {
						continue
					}

					key := field.Tag.Get("json")
					if key == "" {
						key = field.Name
					}

					if marshaller, e := recfun(field.Type); e != nil {

						marshaller_err = fmt.Errorf("Can not marshal field %v of %v to json: %w", field.Name, target, e)
						return

					} else {

						fieldMarshallers = append(fieldMarshallers, func(v reflect.Value, b *JsonBuilder) error {

							if json, err := marshaller.f(v.Field(f)); err != nil {
								return err
							} else {
								b.Set(key, json)
								return nil
							}

						})

					}

				}

			})
			if marshaller_err != nil {
				return nil, fmt.Errorf("Can not marshal %v to json: %w", target, marshaller_err)
			}

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

		},
	}, nil

}

// Slice marshaller

func newSliceMarshaller(
	target reflect.Type,
	recfun func(reflect.Type) (*wrappedMarshaller, error)) (*wrappedMarshaller, error) {

	if target.Kind() != reflect.Slice {
		return nil, fmt.Errorf("The target %v is not a slice.", target)
	}

	var (
		once           sync.Once
		marshaller     *wrappedMarshaller
		marshaller_err error
	)

	return &wrappedMarshaller{
		t: target,
		f: func(v reflect.Value) (JsonNode, error) {

			once.Do(func() {
				marshaller, marshaller_err = recfun(target.Elem())
			})
			if marshaller_err != nil {
				return nil, fmt.Errorf("Can not marshal %v to json: %w", target, marshaller_err)
			}

			if v.IsZero() {
				return JSON_NULL, nil
			}
			elts := make([]JsonNode, v.Len())
			var err error
			for i := 0; i < v.Len(); i++ {
				if elts[i], err = marshaller.f(v.Index(i)); err != nil {
					return JSON_EMPTY_ARRAY, err
				}
			}
			return NewJsonArray(elts), nil

		},
	}, nil

}

// Map marshaller

func newMapMarshaller(
	target reflect.Type,
	recfun func(reflect.Type) (*wrappedMarshaller, error)) (*wrappedMarshaller, error) {

	if target.Kind() != reflect.Map {
		return nil, fmt.Errorf("The target %v is not a map.", target)
	} else if target.Key() != stringType {
		return nil, fmt.Errorf("The target %v key type is not string.", target)
	}

	var (
		once           sync.Once
		marshaller     *wrappedMarshaller
		marshaller_err error
	)

	return &wrappedMarshaller{
		t: target,
		f: func(v reflect.Value) (JsonNode, error) {

			once.Do(func() {
				marshaller, marshaller_err = recfun(target.Elem())
			})
			if marshaller_err != nil {
				return nil, fmt.Errorf("Can not marshal %v to json: %w", target, marshaller_err)
			}

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

		},
	}, nil

}

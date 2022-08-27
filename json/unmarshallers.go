package json

import (
	"fmt"
	"reflect"

	"github.com/b-charles/pigs/ioc"
)

// Defaults unmarshallers

func init() {

	ioc.Put(func(json JsonNode) (string, error) {
		if json.IsNull() {
			return "", nil
		} else if !json.IsString() {
			return "", fmt.Errorf("Can not parse json %v as a string.", json)
		} else {
			return json.AsString(), nil
		}
	}, func(JsonUnmarshaller) {})

	ioc.Put(func(json JsonNode) (float64, error) {
		if json.IsNull() {
			return 0, nil
		} else if !json.IsFloat() {
			return 0, fmt.Errorf("Can not parse json %v as a float.", json)
		} else {
			return json.AsFloat(), nil
		}
	}, func(JsonUnmarshaller) {})

	ioc.Put(func(json JsonNode) (int, error) {
		if json.IsNull() {
			return 0, nil
		} else if !json.IsInt() {
			return 0, fmt.Errorf("Can not parse json %v as an integer.", json)
		} else {
			return json.AsInt(), nil
		}
	}, func(JsonUnmarshaller) {})

	ioc.Put(func(json JsonNode) (bool, error) {
		if json.IsNull() {
			return false, nil
		} else if !json.IsBool() {
			return false, fmt.Errorf("Can not parse json %v as a boolean.", json)
		} else {
			return json.AsBool(), nil
		}
	}, func(JsonUnmarshaller) {})

}

// Pointer unmarshaller

func newPointerUnmarshaller(mapper *JsonMapper, target reflect.Type, valueUnmarshaller *jsonValueUnmarshaller) error {

	if target.Kind() != reflect.Pointer {
		return fmt.Errorf("The target %v is not a pointer.", target)
	}

	if unmarshaller, err := mapper.getUnmarshaller(target.Elem()); err != nil {
		return fmt.Errorf("Can not unmarshal %v to json: %w", target, err)
	} else {
		valueUnmarshaller.f = func(json JsonNode) (reflect.Value, error) {
			if json.IsNull() {
				return reflect.Zero(target), nil
			} else if v, err := unmarshaller.f(json); err != nil {
				return reflect.Value{}, err
			} else {
				return v.Addr(), nil
			}
		}
		return nil
	}

}

// Struct unmarshaller

func newStructUnmarshaller(mapper *JsonMapper, target reflect.Type, valueUnmarshaller *jsonValueUnmarshaller) error {

	if target.Kind() != reflect.Struct {
		return fmt.Errorf("The target %v is not a struct.", target)
	}

	fieldsUnmarshallers := make([]func(JsonNode, reflect.Value) error, 0, target.NumField())
	for f := 0; f < target.NumField(); f++ {

		field := target.Field(f)
		if !field.IsExported() {
			continue
		}

		key := field.Tag.Get("json")
		if key == "" {
			key = field.Name
		}

		if unmarshaller, err := mapper.getUnmarshaller(field.Type); err != nil {
			return fmt.Errorf("Can not unmarshal field %v of %v to json: %w", field.Name, target, err)
		} else {
			numField := f
			fieldsUnmarshallers = append(fieldsUnmarshallers, func(json JsonNode, value reflect.Value) error {
				if v, err := unmarshaller.f(json.GetMember(key)); err != nil {
					return fmt.Errorf("Can't unmarshall field %v: %w", key, err)
				} else {
					value.Field(numField).Set(v)
					return nil
				}
			})
		}

	}

	valueUnmarshaller.f = func(json JsonNode) (reflect.Value, error) {

		value := reflect.New(target).Elem()

		if json.IsNull() {
			return value, nil
		} else if !json.IsObject() {
			return reflect.Value{}, fmt.Errorf("Can not parse json %v as an object.", json)
		} else {
			for _, unmarshaller := range fieldsUnmarshallers {
				if err := unmarshaller(json, value); err != nil {
					return reflect.Value{}, err
				}
			}
			return value, nil
		}

	}
	return nil

}

// Slice unmarshaller

func newSliceUnmarshaller(mapper *JsonMapper, target reflect.Type, valueUnmarshaller *jsonValueUnmarshaller) error {

	if target.Kind() != reflect.Slice {
		return fmt.Errorf("The target %v is not a slice.", target)
	}

	if unmarshaller, err := mapper.getUnmarshaller(target.Elem()); err != nil {
		return fmt.Errorf("Can not unmarshal %v to json: %w", target, err)
	} else {

		valueUnmarshaller.f = func(json JsonNode) (reflect.Value, error) {

			if json.IsNull() {
				return reflect.Zero(target), nil
			} else if !json.IsArray() {
				return reflect.Value{}, fmt.Errorf("Can not parse json %v as a slice.", json)
			} else {
				l := json.GetLen()
				elts := reflect.MakeSlice(target, l, l)
				for i := 0; i < l; i++ {
					if e, err := unmarshaller.f(json.GetElement(i)); err != nil {
						return reflect.Value{}, fmt.Errorf("Can't unmarshall element #%v: %w", i, err)
					} else {
						elts.Index(i).Set(e)
					}
				}
				return elts, nil
			}

		}

		return nil

	}

}

// Map unmarshaller

func newMapUnMarshaller(mapper *JsonMapper, target reflect.Type, valueUnmarshaller *jsonValueUnmarshaller) error {

	if target.Kind() != reflect.Map {
		return fmt.Errorf("The target %v is not a map.", target)
	} else if target.Key() != stringType {
		return fmt.Errorf("The target %v key type is not string.", target)
	}

	if unmarshaller, err := mapper.getUnmarshaller(target.Elem()); err != nil {
		return fmt.Errorf("Can not unmarshal %v to json: %w", target, err)
	} else {

		valueUnmarshaller.f = func(json JsonNode) (reflect.Value, error) {

			if json.IsNull() {
				return reflect.Zero(target), nil
			} else if !json.IsObject() {
				return reflect.Value{}, fmt.Errorf("Can not parse json %v as a slice.", json)
			} else {
				members := reflect.MakeMap(target)
				for _, k := range json.GetKeys() {
					if e, err := unmarshaller.f(json.GetMember(k)); err != nil {
						return reflect.Value{}, fmt.Errorf("Can't unmarshall member '%v': %w", k, err)
					} else {
						members.SetMapIndex(reflect.ValueOf(k), e)
					}
				}
				return members, nil
			}

		}

		return nil

	}

}

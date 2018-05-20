package smartconf

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

// type Parser interface {
// 	   Parse(env map[string]string, root string) (interface{}, error)
// }

var string_type reflect.Type = reflect.TypeOf("")
var map_string_string_type reflect.Type = reflect.TypeOf(map[string]string{})
var error_type reflect.Type = reflect.TypeOf((*error)(nil)).Elem()

type SmartConf struct {
	env     map[string]string
	parsers map[reflect.Type]reflect.Value
}

func NewSmartConf(env map[string]string, parsers []interface{}) *SmartConf {

	parserMap := make(map[reflect.Type]reflect.Value)
	for _, p := range parsers {

		var method reflect.Value
		if pvalue := reflect.ValueOf(p); pvalue.Kind() == reflect.Func {
			method = pvalue
		} else {
			if method = pvalue.MethodByName("Parse"); !method.IsValid() {
				panic(fmt.Errorf("The given object doesn't implement a 'Parse' method: %v", p))
			}
		}

		methodType := method.Type()
		var numIn int

		if numIn = methodType.NumIn(); numIn != 1 && numIn != 2 {
			panic(fmt.Errorf("The 'Parse' method of %v should take 1 or 2 inputs.", p))
		}

		if numIn == 1 {
			if firstIn := methodType.In(0); firstIn != string_type {
				panic(fmt.Errorf("The argument of the 'Parse' method of %v should be a string, not a %v.", p, firstIn))
			}
		} else { // numIn == 2
			if firstIn := methodType.In(0); firstIn != map_string_string_type {
				panic(fmt.Errorf("The first argument of the 'Parse' method of %v should be a map[string]string, not a %v.", p, firstIn))
			}
			if secondIn := methodType.In(1); secondIn != string_type {
				panic(fmt.Errorf("The second argument of the 'Parse' method of %v should be a string, not a %v.", p, secondIn))
			}
		}

		if methodType.NumOut() != 2 {
			panic(fmt.Errorf("The 'Parse' method of %v should return 2 outputs.", p))
		}
		if secondOut := methodType.Out(1); secondOut != error_type {
			panic(fmt.Errorf("The second output of the 'Parse' method of %v should be an error, not a %v.", p, secondOut))
		}

		outType := methodType.Out(0)
		if _, ok := parserMap[outType]; ok {
			panic(fmt.Errorf("A parser for type %v is already registered.", outType))
		}

		// shortForm -> longForm
		if numIn == 1 {

			parserMap[outType] = reflect.ValueOf(func(env map[string]string, root string) (interface{}, error) {
				if value, ok := env[root]; ok {

					outs := method.Call([]reflect.Value{reflect.ValueOf(value)})

					var err error = nil
					if interr := outs[1]; !interr.IsNil() {
						err = interr.Interface().(error)
					}

					return outs[0].Interface(), err

				} else {
					return nil, fmt.Errorf("No value defined for '%s'.", root)
				}
			})

		} else {

			parserMap[outType] = method

		}

	}

	return &SmartConf{
		env:     env,
		parsers: parserMap,
	}

}

func (self *SmartConf) findParser(class reflect.Type) (reflect.Value, error) {

	var parser reflect.Value = reflect.Value{}

	for ret, p := range self.parsers {
		if ret.AssignableTo(class) {
			if parser.IsValid() {
				return parser, fmt.Errorf("More than one parser found for type %v", class)
			} else {
				parser = p
			}
		}
	}

	return parser, nil

}

func (self *SmartConf) callParser(parser reflect.Value, root string) (reflect.Value, error) {

	outs := parser.Call([]reflect.Value{reflect.ValueOf(self.env), reflect.ValueOf(root)})

	var value reflect.Value = outs[0]
	if value.Kind() == reflect.Interface {
		value = value.Elem()
	}

	var err error = nil
	if interr := outs[1]; !interr.IsNil() {
		err = interr.Interface().(error)
	}

	return value, err

}

func (self *SmartConf) mustFindParser(class reflect.Type) (reflect.Value, error) {

	if parser, err := self.findParser(class); !parser.IsValid() && err == nil {
		return parser, fmt.Errorf("No parser found for %v.", class)
	} else {
		return parser, err
	}

}

func (self *SmartConf) getAndConvert(name string, converter func(string) (reflect.Value, error)) (reflect.Value, error) {

	if str, ok := self.env[name]; ok {
		return converter(str)
	} else {
		return reflect.Value{}, fmt.Errorf("The value '%s' is undefined.", name)
	}

}

func (self *SmartConf) getValue(name string, class reflect.Type) (reflect.Value, error) {

	switch class.Kind() {

	case reflect.Bool:
		return self.getAndConvert(name, func(str string) (reflect.Value, error) {
			value, err := strconv.ParseBool(str)
			return reflect.ValueOf(bool(value)), err
		})

	case reflect.Int:
		return self.getAndConvert(name, func(str string) (reflect.Value, error) {
			value, err := strconv.ParseInt(str, 0, 0)
			return reflect.ValueOf(int(value)), err
		})

	case reflect.Int8:
		return self.getAndConvert(name, func(str string) (reflect.Value, error) {
			value, err := strconv.ParseInt(str, 0, 8)
			return reflect.ValueOf(int8(value)), err
		})

	case reflect.Int16:
		return self.getAndConvert(name, func(str string) (reflect.Value, error) {
			value, err := strconv.ParseInt(str, 0, 16)
			return reflect.ValueOf(int16(value)), err
		})

	case reflect.Int32:
		return self.getAndConvert(name, func(str string) (reflect.Value, error) {
			value, err := strconv.ParseInt(str, 0, 32)
			return reflect.ValueOf(int32(value)), err
		})

	case reflect.Int64:
		return self.getAndConvert(name, func(str string) (reflect.Value, error) {
			value, err := strconv.ParseInt(str, 0, 64)
			return reflect.ValueOf(int64(value)), err
		})

	case reflect.Uint:
		return self.getAndConvert(name, func(str string) (reflect.Value, error) {
			value, err := strconv.ParseUint(str, 0, 0)
			return reflect.ValueOf(uint(value)), err
		})

	case reflect.Uint8:
		return self.getAndConvert(name, func(str string) (reflect.Value, error) {
			value, err := strconv.ParseUint(str, 0, 8)
			return reflect.ValueOf(uint8(value)), err
		})

	case reflect.Uint16:
		return self.getAndConvert(name, func(str string) (reflect.Value, error) {
			value, err := strconv.ParseUint(str, 0, 16)
			return reflect.ValueOf(uint16(value)), err
		})

	case reflect.Uint32:
		return self.getAndConvert(name, func(str string) (reflect.Value, error) {
			value, err := strconv.ParseUint(str, 0, 32)
			return reflect.ValueOf(uint32(value)), err
		})

	case reflect.Uint64:
		return self.getAndConvert(name, func(str string) (reflect.Value, error) {
			value, err := strconv.ParseUint(str, 0, 64)
			return reflect.ValueOf(uint64(value)), err
		})

	case reflect.Float32:
		return self.getAndConvert(name, func(str string) (reflect.Value, error) {
			value, err := strconv.ParseFloat(str, 32)
			return reflect.ValueOf(float32(value)), err
		})

	case reflect.Float64:
		return self.getAndConvert(name, func(str string) (reflect.Value, error) {
			value, err := strconv.ParseFloat(str, 64)
			return reflect.ValueOf(float64(value)), err
		})

	case reflect.Interface:

		parser, err := self.mustFindParser(class)
		if err != nil {
			return reflect.Value{}, err
		}

		return self.callParser(parser, name)

	case reflect.Map:

		if keyClass := class.Key(); keyClass != string_type {
			return reflect.Value{}, fmt.Errorf("Only map with key of type string are supported (found: %v).", keyClass)
		}

		re := regexp.MustCompile(fmt.Sprintf("^%s\\.([^.]+)", regexp.QuoteMeta(name)))

		keys := make(map[string]bool)
		for key, _ := range self.env {
			if match := re.FindStringSubmatch(key); match != nil {
				keys[match[1]] = true
			}
		}

		mapValue := reflect.MakeMapWithSize(class, len(keys))

		elemType := class.Elem()
		for key, _ := range keys {
			if value, err := self.getValue(fmt.Sprintf("%s.%s", name, key), elemType); err == nil {
				mapValue.SetMapIndex(reflect.ValueOf(key), value)
			} else {
				return reflect.Value{}, errors.Wrapf(err, "Error during parse value of map index '%s'", key)
			}
		}

		return mapValue, nil

	case reflect.Ptr:

		if parser, err := self.findParser(class); err != nil {
			return reflect.Value{}, err
		} else if parser.IsValid() {
			return self.callParser(parser, name)
		}

		if value, err := self.getValue(name, class.Elem()); err == nil {
			return value.Addr(), nil
		} else {
			return reflect.Value{}, err
		}

	case reflect.Slice:

		sliceSize := 0
		for mapHasKeyStartingWith(self.env, fmt.Sprintf("%s[%d]", name, sliceSize)) {
			sliceSize++
		}

		sliceValue := reflect.MakeSlice(class, 0, sliceSize)

		elemType := class.Elem()
		for i := 0; i < sliceSize; i++ {
			if value, err := self.getValue(fmt.Sprintf("%s[%d]", name, i), elemType); err == nil {
				sliceValue = reflect.Append(sliceValue, value)
			} else {
				return reflect.Value{}, errors.Wrapf(err, "Error during parse value of slice at index %d", i)
			}
		}

		return sliceValue, nil

	case reflect.String:
		return self.getAndConvert(name, func(str string) (reflect.Value, error) {
			return reflect.ValueOf(str), nil
		})

	case reflect.Struct:

		if parser, err := self.findParser(class); err != nil {
			return reflect.Value{}, err
		} else if parser.IsValid() {
			return self.callParser(parser, name)
		}

		structValue := reflect.New(class).Elem()

		for i := 0; i < structValue.NumField(); i++ {

			field := class.Field(i)
			fieldValue := structValue.Field(i)

			if !fieldValue.CanSet() {
				return reflect.Value{}, fmt.Errorf("Field '%s' of %v is not settable.", field.Name, class)
			}

			if value, err := self.getValue(name+lowerize(field.Name), field.Type); err == nil {
				structValue.Field(i).Set(value)
			} else {
				return reflect.Value{}, errors.Wrapf(err, "Error during parse of field '%s' of %v", field.Name, class)
			}

		}

		return structValue, nil

	default:
		return reflect.Value{}, fmt.Errorf("Unexpected type: %v.", class)

	}

}

func (self *SmartConf) Configure(root string, obj interface{}) (interface{}, error) {

	if value, err := self.getValue(root, reflect.TypeOf(obj)); err == nil {
		return value.Interface(), nil
	} else {
		return nil, err
	}

}

func (self *SmartConf) MustConfigure(root string, obj interface{}) interface{} {

	if value, err := self.Configure(root, obj); err != nil {
		panic(err)
	} else {
		return value
	}

}

func mapHasKeyStartingWith(env map[string]string, start string) bool {

	for key, _ := range env {
		if strings.HasPrefix(key, start) {
			return true
		}
	}

	return false

}

var upper_chars *regexp.Regexp = regexp.MustCompile("[A-Z]")

func lowerize(value string) string {

	return upper_chars.ReplaceAllStringFunc(value, func(str string) string {
		return "." + strings.ToLower(str)
	})

}

package smartconfig

import "reflect"

var (
	string_type          = reflect.TypeOf("")
	navconfig_type       = reflect.TypeOf(func(NavConfig) {}).In(0)
	error_type           = reflect.TypeOf(func(error) {}).In(0)
	smartConfigurer_type = reflect.TypeOf(func(*SmartConfigurer) {}).In(0)
)

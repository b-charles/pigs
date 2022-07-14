package ioc

import "sync"

var instance *Container
var once sync.Once

func ContainerInstance() *Container {
	once.Do(func() {
		instance = NewContainer()
	})
	return instance
}

func ErroneousPutNamedFactory(factory any, name string, aliases ...any) error {
	return ContainerInstance().PutNamedFactory(factory, name, aliases...)
}

func ErroneousPutNamed(object any, name string, aliases ...any) error {
	return ContainerInstance().PutNamed(object, name, aliases...)
}

func ErroneousPutFactory(factory any, aliases ...any) error {
	return ContainerInstance().PutFactory(factory, aliases...)
}

func ErroneousPut(object any, aliases ...any) error {
	return ContainerInstance().Put(object, aliases...)
}

func ErroneousTestPutNamedFactory(factory any, name string, aliases ...any) error {
	return ContainerInstance().TestPutNamedFactory(factory, name, aliases...)
}

func ErroneousTestPutNamed(object any, name string, aliases ...any) error {
	return ContainerInstance().TestPutNamed(object, name, aliases...)
}

func ErroneousTestPutFactory(factory any, aliases ...any) error {
	return ContainerInstance().TestPutFactory(factory, aliases...)
}

func ErroneousTestPut(object any, aliases ...any) error {
	return ContainerInstance().TestPut(object, aliases...)
}

func ErroneousCallInjected(method any) error {
	return ContainerInstance().CallInjected(method)
}

func PutNamedFactory(factory any, name string, aliases ...any) {
	err := ErroneousPutNamed(factory, name, aliases...)
	if err != nil {
		panic(err)
	}
}

func PutNamed(object any, name string, aliases ...any) {
	err := ErroneousPutNamed(object, name, aliases...)
	if err != nil {
		panic(err)
	}
}

func PutFactory(factory any, aliases ...any) {
	err := ErroneousPutFactory(factory, aliases...)
	if err != nil {
		panic(err)
	}
}

func Put(object any, aliases ...any) {
	err := ErroneousPut(object, aliases...)
	if err != nil {
		panic(err)
	}
}

func TestPutNamedFactory(factory any, name string, aliases ...any) {
	err := ErroneousTestPutNamedFactory(factory, name, aliases...)
	if err != nil {
		panic(err)
	}
}

func TestPutNamed(object any, name string, aliases ...any) {
	err := ErroneousTestPutNamed(object, name, aliases...)
	if err != nil {
		panic(err)
	}
}

func TestPutFactory(factory any, aliases ...any) {
	err := ErroneousTestPutFactory(factory, aliases...)
	if err != nil {
		panic(err)
	}
}

func TestPut(object any, aliases ...any) {
	err := ErroneousTestPut(object, aliases...)
	if err != nil {
		panic(err)
	}
}

func CallInjected(method any) {
	err := ErroneousCallInjected(method)
	if err != nil {
		panic(err)
	}
}

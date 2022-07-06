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

func PutNamedFactory(factory interface{}, name string, aliases ...string) error {
	return ContainerInstance().PutNamedFactory(factory, name, aliases...)
}

func PutNamed(object interface{}, name string, aliases ...string) error {
	return ContainerInstance().PutNamed(object, name, aliases...)
}

func PutFactory(factory interface{}, aliases ...string) error {
	return ContainerInstance().PutFactory(factory, aliases...)
}

func Put(object interface{}, aliases ...string) error {
	return ContainerInstance().Put(object, aliases...)
}

func TestPutNamedFactory(factory interface{}, name string, aliases ...string) error {
	return ContainerInstance().TestPutNamedFactory(factory, name, aliases...)
}

func TestPutNamed(object interface{}, name string, aliases ...string) error {
	return ContainerInstance().TestPutNamed(object, name, aliases...)
}

func TestPutFactory(factory interface{}, aliases ...string) error {
	return ContainerInstance().TestPutFactory(factory, aliases...)
}

func TestPut(object interface{}, aliases ...string) error {
	return ContainerInstance().TestPut(object, aliases...)
}

func CallInjected(method interface{}) error {
	return ContainerInstance().CallInjected(method)
}

func Close() {
	ContainerInstance().Close()
}

func ClearTests() {
	ContainerInstance().ClearTests()
}

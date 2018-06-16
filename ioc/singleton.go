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

func PutFactory(factory interface{}, name string, aliases ...string) error {
	return ContainerInstance().PutFactory(factory, name, aliases...)
}

func Put(object interface{}, name string, aliases ...string) error {
	return ContainerInstance().Put(object, name, aliases...)
}

func TestPutFactory(factory interface{}, name string, aliases ...string) error {
	return ContainerInstance().TestPutFactory(factory, name, aliases...)
}

func TestPut(object interface{}, name string, aliases ...string) error {
	return ContainerInstance().TestPut(object, name, aliases...)
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

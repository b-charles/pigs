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

func PutFactory(factory interface{}, inputs []string, name string, aliases ...string) bool {
	ContainerInstance().PutFactory(factory, inputs, name, aliases...)
	return true
}

func Put(object interface{}, name string, aliases ...string) bool {
	ContainerInstance().Put(object, name, aliases...)
	return true
}

func TestPutFactory(factory interface{}, inputs []string, name string, aliases ...string) bool {
	ContainerInstance().TestPutFactory(factory, inputs, name, aliases...)
	return true
}

func TestPut(object interface{}, name string, aliases ...string) bool {
	ContainerInstance().TestPut(object, name, aliases...)
	return true
}

func GetComponents(name string) []interface{} {
	return ContainerInstance().GetComponents(name)
}

func GetComponent(name string) interface{} {
	return ContainerInstance().GetComponent(name)
}

func Close() {
	ContainerInstance().Close()
}

func ClearTests() {
	ContainerInstance().ClearTests()
}

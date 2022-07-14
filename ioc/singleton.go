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

func PutNamedFactory(factory any, name string, aliases ...any) error {
	return ContainerInstance().PutNamedFactory(factory, name, aliases...)
}

func PutNamed(object any, name string, aliases ...any) error {
	return ContainerInstance().PutNamed(object, name, aliases...)
}

func PutFactory(factory any, aliases ...any) error {
	return ContainerInstance().PutFactory(factory, aliases...)
}

func Put(object any, aliases ...any) error {
	return ContainerInstance().Put(object, aliases...)
}

func TestPutNamedFactory(factory any, name string, aliases ...any) error {
	return ContainerInstance().TestPutNamedFactory(factory, name, aliases...)
}

func TestPutNamed(object any, name string, aliases ...any) error {
	return ContainerInstance().TestPutNamed(object, name, aliases...)
}

func TestPutFactory(factory any, aliases ...any) error {
	return ContainerInstance().TestPutFactory(factory, aliases...)
}

func TestPut(object any, aliases ...any) error {
	return ContainerInstance().TestPut(object, aliases...)
}

func CallInjected(method any) {
	err := ContainerInstance().CallInjected(method)
	if err != nil {
		panic(err)
	}
}

func UnsureCallInjected(method any) error {
	return ContainerInstance().CallInjected(method)
}

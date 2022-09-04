package ioc

import (
	"sync"
)

var (
	containerInstance *Container
	once              sync.Once
)

// ContainerInstance returns a unique Container instance (singleton).
func ContainerInstance() *Container {
	once.Do(func() {
		containerInstance = NewContainer()
	})
	return containerInstance
}

// ErroneousDefaultPutFactory records a default component defined by a factory and
// optional signatures and returns an error if something wrong happened.
func ErroneousDefaultPutFactory(factory any, signFuncs ...any) error {
	return ContainerInstance().DefaultPutFactory(factory, signFuncs...)
}

// ErroneousDefaultPut records directly a default component with optional signatures
// and returns an error if something wrong happened.
func ErroneousDefaultPut(object any, signFuncs ...any) error {
	return ContainerInstance().DefaultPut(object, signFuncs...)
}

// ErroneousPutFactory records a component defined by a factory and optional
// signatures and returns an error if something wrong happened.
func ErroneousPutFactory(factory any, signFuncs ...any) error {
	return ContainerInstance().PutFactory(factory, signFuncs...)
}

// ErroneousPut records directly a component with optional signatures and
// returns an error if something wrong happened.
func ErroneousPut(object any, signFuncs ...any) error {
	return ContainerInstance().Put(object, signFuncs...)
}

// ErroneousTestPutFactory records a test component defined by a factory and
// optional signatures and returns an error if something wrong happened.
func ErroneousTestPutFactory(factory any, signFuncs ...any) error {
	return ContainerInstance().TestPutFactory(factory, signFuncs...)
}

// ErroneousTestPut records directly a test component with optional signatures
// and returns an error if something wrong happened.
func ErroneousTestPut(object any, signFuncs ...any) error {
	return ContainerInstance().TestPut(object, signFuncs...)
}

// ErroneousCallInjected call the given method, injecting its arguments and
// returns an error if something wrong happened.
func ErroneousCallInjected(method any) error {
	return ContainerInstance().CallInjected(method)
}

// DefaultPutFactory records a default component defined by a factory and optional
// signatures. Panics if something wrong happened.
func DefaultPutFactory(factory any, signFuncs ...any) {
	err := ErroneousDefaultPutFactory(factory, signFuncs...)
	if err != nil {
		panic(err)
	}
}

// DefaultPut records directly a default component with optional signatures. Panics
// if something wrong happened.
func DefaultPut(object any, signFuncs ...any) {
	err := ErroneousDefaultPut(object, signFuncs...)
	if err != nil {
		panic(err)
	}
}

// PutFactory records a component defined by a factory and optional
// signatures. Panics if something wrong happened.
func PutFactory(factory any, signFuncs ...any) {
	err := ErroneousPutFactory(factory, signFuncs...)
	if err != nil {
		panic(err)
	}
}

// Put records directly a component with optional signatures. Panics if
// something wrong happened.
func Put(object any, signFuncs ...any) {
	err := ErroneousPut(object, signFuncs...)
	if err != nil {
		panic(err)
	}
}

// TestPutFactory records a test component defined by a factory and optional
// signatures. Panics if something wrong happened.
func TestPutFactory(factory any, signFuncs ...any) {
	err := ErroneousTestPutFactory(factory, signFuncs...)
	if err != nil {
		panic(err)
	}
}

// TestPut records directly a test component with optional signatures. Panics
// if something wrong happened.
func TestPut(object any, signFuncs ...any) {
	err := ErroneousTestPut(object, signFuncs...)
	if err != nil {
		panic(err)
	}
}

// CallInjected call the given method, injecting its arguments. Panics if
// something wrong happened.
func CallInjected(method any) {
	err := ErroneousCallInjected(method)
	if err != nil {
		panic(err)
	}
}

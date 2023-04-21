package ioc

import "sync"

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

// ErroneousDefaultPutNamedFactory records a default component defined by its
// name, its factory and optional signatures and returns an error if something
// wrong happened.
func ErroneousDefaultPutNamedFactory(name string, factory any, signFuncs ...any) error {
	return ContainerInstance().RegisterFactory(Def, name, factory, signFuncs...)
}

// ErroneousDefaultPutFactory records an anonymous default component defined by
// its factory and optional signatures and returns an error if something wrong
// happened.
func ErroneousDefaultPutFactory(factory any, signFuncs ...any) error {
	return ContainerInstance().RegisterFactory(Def, "", factory, signFuncs...)
}

// ErroneousDefaultPutNamed records a default component defined by its name,
// its value and optional signatures and returns an error if something wrong
// happened.
func ErroneousDefaultPutNamed(name string, object any, signFuncs ...any) error {
	return ContainerInstance().RegisterComponent(Def, name, object, signFuncs...)
}

// ErroneousDefaultPut records an anonymous default component defined by its
// value and optional signatures and returns an error if something wrong
// happened.
func ErroneousDefaultPut(object any, signFuncs ...any) error {
	return ContainerInstance().RegisterComponent(Def, "", object, signFuncs...)
}

// ErroneousPutNamedFactory records a core component defined by its name, its
// factory and optional signatures and returns an error if something wrong
// happened.
func ErroneousPutNamedFactory(name string, factory any, signFuncs ...any) error {
	return ContainerInstance().RegisterFactory(Core, name, factory, signFuncs...)
}

// ErroneousPutFactory records an anonymous core component defined by its
// factory and optional signatures and returns an error if something wrong
// happened.
func ErroneousPutFactory(factory any, signFuncs ...any) error {
	return ContainerInstance().RegisterFactory(Core, "", factory, signFuncs...)
}

// ErroneousPutNamed records a core component defined by its name, its value
// and optional signatures and returns an error if something wrong happened.
func ErroneousPutNamed(name string, object any, signFuncs ...any) error {
	return ContainerInstance().RegisterComponent(Core, name, object, signFuncs...)
}

// ErroneousPut records directly an anonymous core component defined by its
// value and optional signatures and returns an error if something wrong
// happened.
func ErroneousPut(object any, signFuncs ...any) error {
	return ContainerInstance().RegisterComponent(Core, "", object, signFuncs...)
}

// ErroneousTestPutNamedFactory records a test component defined by its name,
// its factory and optional signatures and returns an error if something wrong
// happened.
func ErroneousTestPutNamedFactory(name string, factory any, signFuncs ...any) error {
	return ContainerInstance().RegisterFactory(Test, name, factory, signFuncs...)
}

// ErroneousTestPutFactory records an anonymous test component defined by its
// factory and optional signatures and returns an error if something wrong
// happened.
func ErroneousTestPutFactory(factory any, signFuncs ...any) error {
	return ContainerInstance().RegisterFactory(Test, "", factory, signFuncs...)
}

// ErroneousTestPut records a test component defined by its name, its value and
// optional signatures and returns an error if something wrong happened.
func ErroneousTestPutNamed(name string, object any, signFuncs ...any) error {
	return ContainerInstance().RegisterComponent(Test, name, object, signFuncs...)
}

// ErroneousTestPut records an anonymous test component defined by its value
// and optional signatures and returns an error if something wrong happened.
func ErroneousTestPut(object any, signFuncs ...any) error {
	return ContainerInstance().RegisterComponent(Test, "", object, signFuncs...)
}

// ErroneousCallInjected call the given method, injecting its arguments and
// returns an error if something wrong happened.
func ErroneousCallInjected(method any) error {
	return ContainerInstance().CallInjected(method)
}

// DefaultPutNamedFactory records a default component defined by its name, its
// factory and optional signatures. Panics if something wrong happened.
func DefaultPutNamedFactory(name string, factory any, signFuncs ...any) {
	if err := ErroneousDefaultPutNamedFactory(name, factory, signFuncs...); err != nil {
		panic(err)
	}
}

// DefaultPutFactory records an anonymous default component defined by its
// factory and optional signatures. Panics if something wrong happened.
func DefaultPutFactory(factory any, signFuncs ...any) {
	if err := ErroneousDefaultPutFactory(factory, signFuncs...); err != nil {
		panic(err)
	}
}

// DefaultPutNamed records a default component defined by its name, its value
// and optional signatures. Panics if something wrong happened.
func DefaultPutNamed(name string, object any, signFuncs ...any) {
	if err := ErroneousDefaultPutNamed(name, object, signFuncs...); err != nil {
		panic(err)
	}
}

// DefaultPut records an anonymous default component defined by its value and
// optional signatures. Panics if something wrong happened.
func DefaultPut(object any, signFuncs ...any) {
	if err := ErroneousDefaultPut(object, signFuncs...); err != nil {
		panic(err)
	}
}

// PutNamedFactory records a core component defined by its name, its factory
// and optional signatures. Panics if something wrong happened.
func PutNamedFactory(name string, factory any, signFuncs ...any) {
	if err := ErroneousPutNamedFactory(name, factory, signFuncs...); err != nil {
		panic(err)
	}
}

// PutFactory records an anonymous core component defined by its factory and
// optional signatures. Panics if something wrong happened.
func PutFactory(factory any, signFuncs ...any) {
	if err := ErroneousPutFactory(factory, signFuncs...); err != nil {
		panic(err)
	}
}

// PutNamed records a core component defined by its name, its value and
// optional signatures. Panics if something wrong happened.
func PutNamed(name string, object any, signFuncs ...any) {
	if err := ErroneousPutNamed(name, object, signFuncs...); err != nil {
		panic(err)
	}
}

// Put records a core component defined by its value and optional signatures.
// Panics if something wrong happened.
func Put(object any, signFuncs ...any) {
	if err := ErroneousPut(object, signFuncs...); err != nil {
		panic(err)
	}
}

// TestPutNamedFactory records a test component defined by its name, its
// factory and optional signatures. Panics if something wrong happened.
func TestPutNamedFactory(name string, factory any, signFuncs ...any) {
	if err := ErroneousTestPutNamedFactory(name, factory, signFuncs...); err != nil {
		panic(err)
	}
}

// TestPutFactory records an anonymous test component defined by its factory
// and optional signatures. Panics if something wrong happened.
func TestPutFactory(factory any, signFuncs ...any) {
	if err := ErroneousTestPutFactory(factory, signFuncs...); err != nil {
		panic(err)
	}
}

// TestPutNamed records a test component defined by its name, its value and
// optional signatures. Panics if something wrong happened.
func TestPutNamed(name string, object any, signFuncs ...any) {
	if err := ErroneousTestPutNamed(name, object, signFuncs...); err != nil {
		panic(err)
	}
}

// TestPut records an anonymous test component defined by its value and
// optional signatures. Panics if something wrong happened.
func TestPut(object any, signFuncs ...any) {
	if err := ErroneousTestPut(object, signFuncs...); err != nil {
		panic(err)
	}
}

// CallInjected call the given method, injecting its arguments. Panics if
// something wrong happened.
func CallInjected(method any) {
	if err := ErroneousCallInjected(method); err != nil {
		panic(err)
	}
}

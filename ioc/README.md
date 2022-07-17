# Pigs - IOC

Yet another IOC framework.

## What's IOC?

If you don't know what is IOC or DI, [this link](http://lmgtfy.com/?s=d&q=what%27s+inversion+of+control%3F) will change your life.

## Goals

This framework has been written with some goals in mind:
* No requirement for the components: third libraries can be injected quickly and nicely.
* A strong typed support: you should never have to cast anything. Never.
* A solution for cyclic dependencies: a component A which needs a component B which needs also the component A, that's not a problem of design, it's just the ruthless real world.
* A support for tests. That means:
    * the possibility to mock up a different environment for each unit test,
    * the lazy loading of components, to load only what you want to test, for each unit test.
* Auto-discovery, also known as voodoo origins.

And some goals are deliberated ignored:
* No notion of scope: it's classical for IOC frameworks to include the notion of scope, which defines visibility and lifecycle of components (singleton, prototype...). But not here.
* No integration with any web framework, but web framework can be nicely integrated in it.
* No AOP support: voodoo is ok, but not satanism.

## How it work

### Container and component

A component is something: a struct, a pointer, a map, a chan, an integer... It's a variable that should be unique (singleton) and which will be injected in all other components that need it. Each component is defined by its type and optionally some _signatures_: some additional interfaces which are implemented by the component and declared when registering the component[^duck]. Multiple components can share the same type and signatures to enable auto-discovery, but that concept can be confusing and error-prone: it can be considered good practice to register each component with a specific and unique type, and using signatures only when auto-discovery is needed.

A container is a set of components: it manages the lifecycle of each component and take in charge the injection process. The container is a managed component itself and can be injected in any component.

[^duck]: The concept of signature is not very in line with the spirit of the duck typing of Go, but it's a need for performance (what should be injected, without checking each component) and flexibility (what should *not* be injected, even if the interface matches the implementation).

### Injections

The framework is based on two types of injection: injection of (pointers to) structures and injection of functions.

#### Injection of structures

The framework can inject a structure, or rather a pointer to a structure.

When the framework inject a structure, each field is observed, and only the fields tagged with `inject:""` are injected. The type of the field defines he component to inject. So if this variable `injected` is injected:
```go
type MyInjectedStruct struct {
    NotInjected *NotInjected
    FirstInject *Something `inject:""`
    SecondInject SomeInterface `inject:""`
}

var injected *MyInjectedStruct = &MyInjectedStruct{}
```
the first field `NotInjected` is not injected. The second and third fields `FirstInject` and `SecondInject` are injected with components with a type or a signature respectively `*Something` and `SomeInterface`. Please note that each injected field should be settable (so exported, with a name beginning with an upper case).

#### Injection of functions

At other points of the framework, the container use some user-defined functions, where each arguments need to be injected. Like for structures, the type of each argument is used to select the component to inject:
```go
func InjectedFunc(first mypkg.SomeInterface, second *anotherlib.RandomComponent) { ... }
```
If the container calls the function, the first argument will be injected with a component with a type or signature `mypkg.SomeInterface`  and the second argument with a component with a type or signature `*anotherlib.RandomComponent`.

Injected functions can take any number of input arguments, none included.

#### Auto-discovery injection

Whether it is an injection in a structure or in a function, if no component or too many components are defined, the container returns an error. But you can inject all components sharing the same type or signature: that's the foundation of the auto-discovery. If the injected value is of type slice, the container will start by trying to resolve the injection normally (a component can be defined as a slice). If no single component match the definition (i.e. if there is no or several components found for inject the value), the container will create a slice, add in it all the components with the correct type or signature, and inject that slice:
```go
type MyInjectedStruct struct {
    AllComponents []SomeInterface `inject:""`
}

var injected *MyInjectedStruct = &MyInjectedStruct{}
```
Here, if `injected` is injected by the container, the field `AllComponents` will contains every components defined with the type or signature `SomeInterface`. If no component are found, the injected slice is empty.

#### Testing

Like it's said in the main goals of the framework, there is no concept of scope, or at least not extendable scope. In fact, the container is divided in two set of components: one for the core components, and one for tests. These two sets are comparable to scopes: every time an injection is proceeded, the container start by searching any component in the test set. If nothing is found, then the container use the core set. In case of auto-discovery injection (slice), if at least one matching component is found in the test set, that components are injected and the components in the core set are not used (which means you can not run a test with an auto-discovered slice empty when it is not outside of testing).

The API is defined to record a component in the core set or in the test set. With that mechanisms, it's easy to define clean unit test that mock some components and limit each test to a part of the application (not the entire application, but not only one component neither).

### Container and components lifecycles

At least. After all these theoretical concepts, we will finally see the concrete part of the framework.

The lifecycle of the container is going through several steps:
* definition,
* exploitation:
  * for each necessary component:
      * instantiation,
      * injection,
      * post-initialization,
  * run main function
  * close all closable component
* redefinition (for tests)

Let's see each steps in details.

### Definition

To be able to use a component, you should start by recording it, or by defining how it will be created. That's the definition step.

You can call the functions `ioc.Put` to record directly the component:
```go
type DemoComponent struct { ... }

func init() {
    ioc.Put(&DemoComponent{})
}
```
Recording a component doesn't inject it (not yet). It is only recorded as a not yet initialized component. Signatures can be defined with one or several little weird functions. For example, `func(FirstInterface, SecondInterface) {}` defines two signatures `FirstInterface` and  `SecondInterface`. Signatures can be defined in the `ioc.Put` function:
```go
type DemoComponent struct { ... }

func init() {
    ioc.Put(&DemoComponent{}, func(FirstInterface, SecondInterface) {})
}
```
When a component is registered with signatures, the framework checks that the component can be casted to each signature type and panics if this is not the case.

Instead of `ioc.Put`, you can also use `ioc.PutFactory` to define a factory, a function which will create the component:
```go
package mypkg

type DemoComponent struct {}

func ComponentFactory() *DemoComponent {
    return &DemoComponent{}
}

func init() {
    ioc.PutFactory(ComponentFactory, func(FirstInterface, SecondInterface) {});
}
```
The function `PutFactory` use the type of the first output to get the type of the component. Like `ioc.Put`, if signatures are defined, `ioc.PutFactory` will check that each signature is implemented by the returned type of the factory. Recording a factory doesn't call it, but it will be used to instantiate the component in a next step. At that time, the arguments of the factory will be injected like it is described in the section [Injection of functions](#injection-of-functions). The factory should at least return the created component, and can also return an error:
```go
type DemoComponent struct {}

func ComponentFactory() (*DemoComponent, error) {
    ...
    if somethingFishy {
        return nil, errors.New("Something's fishy")
    }
    ...
    return &DemoComponent{}, nil
}

func init() {
    ioc.PutFactory(ComponentFactory)
}
```

Factories can not define cyclic dependencies (i.e. a factory produces a component `A` which is needed to another factory to create a component `B` which should be injected in the factory of `A`). To resolve the problem, you have to break the cycle, or wait another step in the component's lifecycle (like [Injection](#injection) or [Post-Initialization](#post-initialization)) to inject the required component.

Like explain in the section [Testing](#testing), the functions `Put` and `PutFactory` define the component in the core set. The functions `TestPut` and `TestPutFactory` can be used in the same way to define a component in the test set.

Finaly, all this method have their `Erroneous*` prefixed version (e.g. `ErroneousTestPut`) which doesn't panic but returns an error if something wrong happened.

### Exploitation

Only defining components doesn't create anything. You have to specify what are the main components, the components you need to instantiate and get to starting your application (or doing some tests). The container will instantiate these components, and also their dependencies. The API defines the function `CallInjected` that can be used to retrieve that main components from the container.
```go
package mypkg

type MyMainComponent struct { ... }

func (self *MyMainComponent) start() { ... }

func init() {
    ioc.Put(&MyMainComponent{})
}

func main() {

    ioc.CallInjected(func(component *MyMainComponent) {
        component.start()
    })

}
```
As you should have guessed, the function `CallInjected` take as sole argument a function that will be injected (see [Injection of functions](#injection-of-functions)). The given function can return nothing or an error. The method `CallInjected` panics if an error occurs, but you can also use `ErroneousCallInjected` which instead returns an error if something had happened.

#### Initialization

During the call of `CallInjected`, the instance of some components will be required. This is how components are created and initialized:

##### Instantiation

Of course, the first step is to create an instance. If the component has been defined with the functions `Put` or `TestPut`, the instance is directly used. If the component has been defined with a factory, with the function `PutFactory` or `TestPutFactory`, the factory is called with its argument injected.

If a factory return a not-null error, the container stops all the process as soon as possible and return the error wrapped in some context messages.

##### Injection

Regardless of the registration mode (instance or factory), if the component is a pointer to a structure, it is injected. The process is defined in the section [Injection of structures](#injection-of-structures).

At this point, cyclic dependencies can be resolved.

##### Post-Initialization

If the component have a method `PostInit`, the container calls it. The method is injected, like described in [Injection of functions](#injection-of-functions). The method should return nothing or an error. Like factories, if one `PostInit` returns a not-null error, the container stops everything and returns the error wrapped in context messages.

#### Run main function

When each necessary components are fully initialized, the framework call the given function of `CallInjected`. That's what we wanted from the start and where your business begins. Be carefully if you are working with multi-threads: the end of this function will trigger the next phase of the container and so the closing of components.

#### Close

After the main function executed, the container will close automatically every component implementing the interface `io.Closer`. If a component panics or returns a non-null error at the call of its method `Close`, the error is silently discarded. The `Close` methods are called in the reverse order of component instantiation (main components first, components without dependencies last).

### Redefinition

It should be only one call of `CallInjected` during the run of the application or at each unit test: after a call of `CallInjected`, every instances are released along the core component definitions if no test component are found. Otherwise, if at least one component is defined in the test scope, instances and only test component definitions are released.

So, if you are running unit tests, you have to defined some fixture to redefined each time all the test components you want to use, and for each test, every component is re-instanced. Be sure to define at least one component in the test scope (even if it is not used and not instanciated) before each call of `CallInjected` or you will loosing the definitions of all the core components.

If no test component is defined, the framework considers that you are really running your application. In this case, in order to consume the least RAM as possible, all unused component instances and all component definitions will be released (forgotten by the framework) and can be deleted by the garbage collector if no other component reference them.

## Usage expample

To a better understanding of how the framework can be used, here an extract of unit tests with `Ginkgo`:

[pkg/main.go]
```go
package pkg

import (
  "fmt"

  "github.com/b-charles/pigs/ioc"
)

type Doer interface {
  do()
}

// Yeller

type Yeller struct {}

func (self *Yeller) do() {
  fmt.Printf("I'm doing!\n")
}

func init() {
  ioc.Put(&Yeller{}, func(Doer){})
}

// TenTimer

type TenTimer struct {
  Doer Doer `inject:""`
}

func (self *TenTimer) doTenTimes() {
  for i := 0; i < 10; i++ {
    self.Doer.do()
  }
}

func init() {
  ioc.Put(&TenTimer{})
}

func main() {
  ioc.CallInjected(func(t *TenTimer) {
    t.doTenTimes()
  })
}
```

[pkg/main_test.go]
```go
package pkg_test

import (
  "testing"

  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"
  "github.com/b-charles/pigs/ioc"
  . "pkg"
)

func TestIoc(t *testing.T) {
  RegisterFailHandler(Fail)
  RunSpecs(t, "App test suite")
}

type Mock struct {
  Called int
}

func (self *Mock) do()  {
  self.Called++
}

var _ = Describe("App tests", func() {

  var mock *Mock

  BeforeEach(func() {
    mock = &Mock{}
    ioc.TestPut(mock, func(Doer){})
  }

  It("should do 10 times with the mock", func() {

    ioc.CallInjected(func(t *TenTimer) {
      t.doTenTimes()
    })

    Expect(mock.Called).To(Equal(10))

  })

})
```


# Pigs - IOC

Yet another IOC framework.

## What's IOC?

If you don't know what is IOC or DI, [this link](https://letmegooglethat.com/?q=what+is+inversion+of+control) will change your life.

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
 * No notion of extendable scope: it's classical for IOC frameworks to include the notion of scope, which defines visibility and life cycle of components (singleton, prototype...). But not here. The framework defines a notion of scopes for a different goal (component definition precedence) and the user can not extend the pre-defined scopes.
 * No integration with any web framework, but web framework can be nicely integrated in it.
 * No AOP support: voodoo is ok, but not satanism.

## How it work

### Container and components

A component is something: a struct, a pointer, a map, a chan, an integer... It's a variable that should be unique (singleton) and which will be injected in all other components that need it. Each component is defined by:
 * an optional name: the name doesn't play any role in the injection process but can be defined for debugging purposes.
 * its value or its factory: you can directly define the component value or produce a function which will create the component value. That kind of function is called a factory.
 * its main type, computed by reflection on the given value or the factory.
 * and optionally some _signature_ interfaces which are some additional interfaces implemented by the component[^duck].

Multiple components can share the same main type and signatures to enable auto-discovery.
Any component with a value `nil` will be silently discarded and never be injected in another component. This behaviour, together with auto-discovery and scopes features, helps implement conditional component creation.

A container is a set of components: it manages the life cycle of each component and take in charge the injection process. The framework defines one instance of this container and redirect all API calls to that container. The container defines by itself two special components of type `ioc.ContainerStatus` and `ioc.ContainerInfo` (see [Special components](#special-components)), and the framework defines also some components for classical integration (see [Default integration](#default-integration)).

[^duck]: The concept of signature is not very in line with the spirit of the duck typing of Go, but it's a need for performance (what should be injected, without checking each component) and flexibility (what should *not* be injected, even if the interface matches the implementation).

### Injections

The framework is based on two types of injection: injection of (pointers to) structures and injection of functions.

#### Injection of structures

The framework can inject a structure, or rather a pointer to a structure.

When the framework inject a structure, only the fields tagged with `inject:""` are injected. The type of the field defines the component to inject. So if this variable `injected` is injected:
```go
type MyInjectedStruct struct {
    NotInjected *NotInjected
    FirstInject *Something `inject:""`
    SecondInject SomeInterface `inject:""`
}

var injected *MyInjectedStruct = &MyInjectedStruct{}
```
the first field `NotInjected` is not injected. The second field `FirstInject` is injected with a component defined with a main type at `*Something` and the third field `SecondInject` is injected with a component with a main type or a signature `SomeInterface`. Please note that each injected field should be settable (so exported, with a name beginning with an upper case).

#### Injection of functions

At other points of the framework, the container use some user-defined functions, where each arguments need to be injected. Like for structures, the type of each argument is used to select the component to inject:
```go
func InjectedFunc(first mypkg.SomeInterface, second *anotherlib.RandomComponent) { ... }
```
If the container calls the function, the first argument will be injected with a component with a main type or signature `mypkg.SomeInterface` and the second argument with a component with a main type `*anotherlib.RandomComponent`.

Injected functions can take any number of input arguments, none included.

#### Auto-discovery injection

Whether it is an injection in a structure or in a function, if no component or too many components are defined, the container returns an error. But you can inject all components sharing the same type or signature: that's the foundation of the auto-discovery. If the injected value is of type slice, the container will start by trying to resolve the injection normally (a component can be defined as a slice). If no component match the definition, the container will create a slice, add in it all the components with the correct type or signature, and inject that slice:
```go
type MyInjectedStruct struct {
    AllComponents []SomeInterface `inject:""`
}

var injected *MyInjectedStruct = &MyInjectedStruct{}
```
Here, if `injected` is injected by the container, the field `AllComponents` will contains every components defined with the type or signature `SomeInterface`. The injected slice will never contain a `nil` value: `nil` components are discarded. If no component is found, the injected slice is empty. This feature can be useful to implement an optional injection.

#### Scopes: Default, Core and Testing

Like it's said in the main goals of the framework, there is no concept of scope, or at least not extendable scope. In fact, the container is divided in three set of components: one for the default components, one for the core components, and one for tests. Every time an injection is proceeded, the container start by searching any component in the test set. If nothing is found or if the components are `nil`, the container use the core set. If still nothing is found (or only `nil` components), then the container use the default set. In case of auto-discovery injection (slice), if at least one matching component is found in the test set, that components are injected and the components in the core set are not used (which also means if you have configured some components to be auto-discovered in the default or core scope, you can not run a test with the auto-discovered slice empty), and in the same spirit, if the core set is used and at least one component is found, the default set is ignored.

The API is defined to record a component in the default set, the core set or in the test set. With that mechanisms, it's easy to:
 * write some libraries with default functionalities which can be overloaded by the main application,
 * define clean unit test that mock some components and limit each test to a part of the application (not the entire application, but not only one component neither).

You can also promote a component to a "superior" scope (a default component to a core or a test component, or a core component as a test component). If during the instantiation of a test component by a factory (a function which returns the component value), a component of the same type produced by the factory is required, the container doesn't throw a cyclic dependency error (as it should) but search if a core or a default component can be used. In the same way, if a core component factory requires a component of the same type, the container search if a default component can be used. This mechanism can be useful to define quickly one component which will be injected in a slice by auto-discovery (even if another more convenient strategy is possible, see the next section), or to configure a component for a test environment.

#### Overloadable components in auto-discovery injection

This section is about an advanced trick used in several place in Pigs libraries and required that you understand the ioc API explained in the next sections. If you read this documentation for the first time, you can skip this part and come back latter.

Sometimes, a problem occurs: how to define a default overloadable component, but in the same time make it available to auto-discovery with other core registered components?

For example, let's take a library which defines a `Burger` component which gathers all components implementing an `Ingredient` interface:
```go
type Ingredient interface {
    ...
}

type Burger struct {
    Ingredients []Ingredient `inject:""`
}

func init() {
    ioc.Put(&Burger{})
}
```
The library was developed to be extensible and must accept `Ingredient` components defined anywhere in the code, to allow some applications to demonstrate gastronomic creativity. The library also proposes a `type Bun struct {...}` which implements `Ingredient`. That `Bun` is interesting enough to be registered as a component, but it should exists an option to overload it: some people may want use bagels instead.

A possible strategy is to only define default implementations without ioc integration, and let the application's main code register the chosen ingredients as core components. The library can also register that default implementations as default components, but if a component is registered in the core scope, the wanted default components have still to be promoted in core scope individually.

A better approach consists of this three steps:
 * An additional specific interface is defined for this default implementation. This interface should at least be assignable to the injected interface, and can be simply defined as a type copy:
   ```go
   type Bread Ingredient
   ```
 * The default implementation is recorded in the default scope, with a signature set to the specific interface (not the injected):
   ```go
   func init() {
      ioc.DefaultPut(&Bum{}, func(Bread){})
   }
   ```
 * A factory is defined to promote any component implementing the specific interface to a core component of the injected type:
   ```go
   func init() {
      ioc.PutFactory(func(bread Bread) Ingredient { return bread })
   }
   ```

By default, the `Bum` component which implements `Bread` will be promoted as an core `Ingredient` by the factory, and will be injected in the `Burger` with other core `Ingredient` components. But it can also be discarded and overloaded by simply register a component in the core or test scope with the `Bread` signature:
```go
func init() {
  ioc.Put(&Bagel{}, func(Bread){})
}
```

### Container and components life cycles

At least. After all these theoretical concepts, we will finally see the concrete part of the framework.

The life cycle of the container is going through several steps:
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

You can call the function `ioc.Put(component any, signatures ...any)` to record directly the component:
```go
type DemoComponent struct { ... }

func init() {
    ioc.Put(&DemoComponent{})
}
```
Calling the function with `nil` as input will do nothing. Recording a component doesn't inject it (not yet). It is only recorded as a not yet initialized component. Signature interfaces can be defined as input argument types of one or several empty functions. For example, `func(FirstInterface, SecondInterface) {}` defines two signatures `FirstInterface` and  `SecondInterface`. Signatures can be provided in the `ioc.Put` function:
```go
type DemoComponent struct { ... }

func init() {
    ioc.Put(&DemoComponent{}, func(FirstInterface, SecondInterface) {})
}
```
When a component is registered with signatures, the framework checks that the component can be casted to each signature type and panics if it is not the case.

You can also give a name or a label to the component, using the function `ioc.PutNamed(name string, component any, signatures ...any)`:
```go
type DemoComponent struct { ... }

func init() {
    ioc.PutNamed("My demo component", &DemoComponent{})
}
```
The name is only for debugging, can be any string, and doesn't modify the behaviour of the container regarding the component.

Instead of `ioc.Put` and `ioc.PutNamed`, you can also use `ioc.PutFactory(factory any, signatures ...any)` and `ioc.PutNamedFactory(name string, factory any, signatures ...any)` to define a factory, a function which will create the component:
```go
package mypkg

type DemoComponent struct {}

func ComponentFactory() *DemoComponent {
    return &DemoComponent{}
}

func init() {
    ioc.PutNamedFactory("My demo component", ComponentFactory, func(FirstInterface, SecondInterface) {});
}
```
The functions `ioc.PutFactory` and `ioc.PutNamedFactory` use the type of the first output to get the type of the component. Like `ioc.Put`, registering a `nil` value will do nothing, and if signatures are defined, `ioc.PutFactory` `ioc.PutNamedFactory` will check that each signature is implemented by the returned type of the factory. Recording a factory doesn't call it, but it will be used to instantiate the component in a next step. At that time, the arguments of the factory will be injected like it is described in the section [Injection of functions](#injection-of-functions). The factory should at least return the created component or `nil`: in that case, the component will be discarded and not used in the future injections. The factory can also return an error:
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
    ioc.PutNamedFactory("My demo component", ComponentFactory)
}
```

Factories can not define cyclic dependencies (i.e. a factory produces a component `A` which is needed to another factory to create a component `B` which should be injected in the factory of `A`). To resolve the problem, you have to break the cycle, or wait another step in the component's life cycle (like [Injection](#injection) or [Post-Initialization](#post-initialization)) to inject the required component.

Like explained in the section [Scopes: Default, Core and Testing](#scopes-default-core-and-testing), the functions `Put`, `PutNamed`, `PutFactory` and `PutNamedFactory` define the component in the core set. The functions `DefaultPut`, `DefaultPutNamed`, `DefaultPutFactory` and `DefaultPutNamedFactory` can be used in the same way to define a component in the default set, and the functions `TestPut`, `TestPutNamed`, `TestPutFactory` and `TestPutNamedFactory` for the test set.

Finally, all this methods have their `Erroneous*` prefixed version (e.g. `ErroneousTestPutNamed(name string, component any, signatures ...any) error`) which doesn't panic but returns an error if something wrong happened.

### Exploitation

Only defining components doesn't create anything. You have to specify what are the main components, the components you need to instantiate and get to starting your application (or doing some tests). The container will instantiate these components, and also their dependencies. The API defines the function `CallInjected(injected any)` that can be used to retrieve that main components from the container.
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
As you should have guessed, the function `CallInjected` take as sole argument a function that will be injected (see [Injection of functions](#injection-of-functions)). The given function can return nothing or an error. The method `CallInjected` panics if an error occurs, but you can also use `ErroneousCallInjected(injected any) error` which instead returns an error if something had happened.

#### Initialization

During the call of `CallInjected`, the instance of some components will be required. This is how components are created and initialized:

##### Instantiation

Of course, the first step is to create an instance. If the component has been directly defined by its value, the instance is used. If the component has been defined with a factory, the factory is called with its argument injected. If the factory returns a not-null error, the container stops all the process as soon as possible and returns the error wrapped in some context message.

##### Injection

Regardless of the registration mode (value or factory), if the component is a pointer to a structure, it is injected. The process is defined in the section [Injection of structures](#injection-of-structures).

At this point, cyclic dependencies are not an issue anymore.

##### Post-Initialization

If the component have a method `PostInit`, the container calls it. The method is injected, like described in [Injection of functions](#injection-of-functions). The method should return nothing or an error. Like factories, if one `PostInit` returns a not-null error, the container stops everything and returns the error wrapped in context messages.

#### Run main function

When each necessary components are fully initialized, the framework call the given function of `CallInjected` (or `ErroneousCallInjected`). That's what we wanted from the start and where your business begins. Be careful if you are working with multi-threads (goroutines): the end of this function will trigger the next phase of the container and so the closing of components.

#### Close

After the main function executed, the container will close automatically every component implementing the interface `io.Closer`. If a component panics or returns a non-null error at the call of its method `Close() error`, the error is silently discarded. The `Close` methods are called in the reverse order of component instantiation (main components first, components without dependencies last).

### Redefinition

It should be only one call of `CallInjected` during the run of the application or at each unit test: after a call of `CallInjected`, every instances are released along the component definitions (default and core) if no test components are found. Otherwise, if at least one component is defined in the test scope (even if it is not instantiated or injected in another component), instances of all scopes and component definitions of test scope only are released.

So, if you are running unit tests, you have to defined some fixture to redefined each time all the test components you want to use, and for each test, every component is re-instanced. Be sure to define at least one component in the test scope, even if it is not used, before each call of `CallInjected` or you will loosing the definitions of all the core components.

If no test component is defined, the framework considers that you are really running your application. In this case, in order to consume the least RAM as possible, all unused component instances and all component definitions will be released (forgotten by the framework) and can be deleted by the garbage collector if no other component reference them.

## Usage example

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

// Yeller, implementing Doer

type Yeller struct {}

func (self *Yeller) do() {
  fmt.Printf("I'm doing!\n")
}

func init() {
  ioc.Put(&Yeller{}, func(Doer){})
}

// TenTimer, repeating 10 times a Doer

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
    t.doTenTimes() // expect displaying `I'm doing!` ten times
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

// Mock, test component implementing Doer

type Mock struct {
  Called int
}

func (self *Mock) do()  {
  self.Called++
}

var _ = Describe("App tests", func() {

  var mock *Mock

  BeforeEach(func() {
    // Register a Mock instance in test scope as a Doer
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

## Special components

### Container status

A special component of type `ioc.ContainerStatus` is defined by the container itself and can be used to analyse and display the internal status of the container. The usage of this component should be limited to a temporary debugging, since injecting this component will maintain a reference to all other components and instances, preventing the garbage collector to release anything recorded in the container.

The `String() string` and `Print()` methods produce gorgeous outputs which should be easy enough to understand without a detailed documentation.

### Container info

A special component  of type `ioc.ContainerInfo` is also defined by the container. It can be injected to get some information about the usage of the container:
 * `TestMode() bool` returns `true` if the container is used in a test environment, i.e. if at least one test component is recorded.
 * `CreationTime() time.Time` returns the time of the container creation.
 * `StartingTime() time.Time` returns the time of the call of `CallInjected` function.
 * `ClosingTime() time.Time` returns the time of the end of the `CallInjected` function. The returned time is the zero value until the `CallInjected` ends, and can be used in the `Close` methods of the components.

Times are generated by using the standard `time` package, and ignore [the Clock integration](#clock). During tests, the container is created once, so the creation time will be constant for all tests, and the starting and closing times are updated at each unit test.

## Default integration

### Afero

The library defines a component in the default scope for [Afero](https://github.com/spf13/afero), which can be used to access to the filesystem:

```go
DefaultPutFactory(afero.NewOsFs, func(afero.Fs) {})
```

### Clock

The library defines also a component in the default scope for [Clock](https://github.com/benbjohnson/clock), which wraps the standard library `time`:

```go
DefaultPutFactory(clock.New, func(clock.Clock) {})
```

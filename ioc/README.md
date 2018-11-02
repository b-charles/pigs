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
* No notion of scope: it's classical for IOC frameworks to include the notion of scope, which defines visibility and lifecycle of components. But not here.
* No integration with any web framework, but web framework can be nicely integrated in it.
* No AOP support: voodoo is ok, but not satanism.

## How it work

### Container and component

A component is something: a struct, a pointer, a map, a chan, an integer... It's a variable that should be unique (singleton) and which will be injected in all other components that need it. Each component is defined with a main name and optionally some aliases. The main name can be any string but `nil` or `""` (empty), and have to be unique in all the application. Aliases can be any string (with the same exceptions) and are not required to be unique (several components can share a same alias).

A container is a set of components: it manages the lifecycle of every component and take in charge the injection process. To know what component should be injected, the container look after names: each injection defines the name (main name or alias) of the component that should be injected. If there is a problem of type (the component injected doesn't correspond to the expected type), the container detects it and returns an error. It's the responsibility of the developer to choose wisely the names of the components to avoid collisions, and to define components and injections that are compatible.

The container is a managed component itself and can be accessed by the name `ApplicationContainer`.

### Injections

The framework is based on two types of injection: injection of structures and injection of functions.

#### Injection of structures

At different points in the framework, it appears the need to inject a structure, or a pointer to a structure. Two flavors are possible: "tagged only" and "all injected".

For the mode "tagged only", only the tagged `inject` fields are injected. The value of the tag defines the name (or the alias) of the component to inject. If the value of the tag is empty, the name of the field is used in place of. So if this variable `injected` is injected:
```go
type MyInjectedStruct struct {
    NotInjected *NotInjected
    FirstInject SomeInterface `inject:"Something"`
    SecondInject AnotherInterface `inject:""`
}

var injected *MyInjectedStruct = &MyInjectedStruct{}
```
the first field `NotInjected` is not injected, the second field `FirstInject` is injected with a component named `Something`, and the third field is injected with a component named `SecondInject`. Please note that each injected field should be settable (so exported, with a name beginning with an upper case).

For the mode "all injected", every field should be injected, and the tag `inject` can be used only to define the correct name of the component to inject.

#### Injection of functions

At other points of the framework, the container use some user-defined functions, where the arguments need to be injected. Two prototypes of function can be injected, and indifferently used: with each argument injected, or with an injected structure.

The first form is when the function is defined with none, one or several arguments, each expected to be injected. In that case, the name of the type of each argument is used as the name (or alias) of the component to inject:
```go
func InjectedFunc(first mypkg.SomeInterface, second *RandomComponent) { ... }
```
If the container calls the function, the first argument will be injected with a component named `SomeInterface` and the second argument with a component named `RandomComponent`. Note that package prefix and pointer star `*` are ignored. Of course, the function can be defined without any argument.

This first form is pretty clear but not very flexible. The other is little more messy, but more useful: the function is defined with only one argument, a structure or a pointer to a structure, which will be injected as defined in the section [Injection of structures](#injection-of-structures), with the flavor "all injected". The type of this sole argument can be properly defined or defined directly in the function:
```go
func injectedFunc(injected struct {
    SomeComponent SomeInterface
    Another RandomInterface `inject:"RandomComponent"`
}) { ... }
```
Here, the argument `injected` will be injected with a component named `SomeComponent` and another one named `RandomComponent`.

#### Auto-discovery injection

Whether it is an injection in a structure or in a function, if no component or too many components are defined, the container returns an error. But you can inject all components sharing an alias: that's the foundation of the auto-discovery. If the injected value is of type slice, the container will start by trying to resolve the injection normally (a component can be defined as a slice). If no component match the definition, the container will create a slice, add in it all the components with an alias matching the name and inject that slice:
```go
type MyInjectedStruct struct {
    AllComponents []SomeInterface `inject:"MagicService"`
}

var injected *MyInjectedStruct = &MyInjectedStruct{}
```
Here, if `injected` is injected by the container, the field `AllComponents` will contains every components defined with the alias `"MagicService"` (at the condition that all respect the interface `SomeInterface`).

With the same spirit, you can also inject a map of components sharing a same alias, where the keys (of type `string`) are the main names of the associated components:
```go
type MyInjectedStruct struct {
    AllComponents map[string]SomeInterface `inject:"MagicService"`
}

var injected *MyInjectedStruct = &MyInjectedStruct{}
```

#### Testing

Like it's said in the main goals of the framework, there is no concept of scope, or at least not extendable scope. In fact, the container is divided in two set of components: one for the core components, and one for tests. These two sets are comparable to scopes: every time an injection is proceeded, the container start by searching any component in the test set. If nothing is found, then the container use the core set. In case of auto-discovery injection (slice or map), if at least one matching component is found in the test set, only that components are injected and the components in the core set are not used.

The API is defined to record a component in the core set or in the test set, to clean up all instances (not the definition) of components and delete the definitions of the test set components. With that mechanisms, it's easy to define clean unit test that mock some components and limit each test to a part of the application (not the entire application, but not only one component neither).

### Component lifecycle

At least. After all these theoretical concepts, we will finally see the concrete part of the framework.

Managed by a container, the lifecycle of each component is going through several steps:
* definition,
* initialization:
    * instantiation,
    * injection,
    * post-initialization,
* exploitation,
* and destruction, at the termination of the application.

Let's see each steps in details.

### Definition

To be able to use a component, you should start by recording it, or by defining how it will be created. That's the definition step.

You can call the functions `ioc.Put` to record directly the component, its main name and optionally its aliases:
```go
type DemoComponent struct { ... }

func init() {
    if err := ioc.Put(&DemoComponent{}, "MyDemoComponent", "DemoComponentAlias"); err != nil {
        panic(err)
    }
}
```
Recording a component doesn't inject it (not yet). It is only recorded as a not yet initialized component. As you can see, the function `Put` returns an error if something bad had happened (e.g. a component with the same main name is already registered).

Instead of `ioc.Put`, you can also use `ioc.PutFactory` to define a factory, a function which will create the component:
```go
type DemoComponent struct {}

func ComponentFactory() *DemoComponent {
    return &DemoComponent{}
}

func init() {
    if err := ioc.PutFactory(ComponentFactory, "MyDemoComponent", "DemoComponentAlias"); err != nil {
        panic(err)
    }
}
```
Recording a factory doesn't call it, but it will be used to instantiate the component in a next step. At that time, the arguments of the factory will be injected like it is described in the section [Injection of functions](#injection-of-functions). The factory should at least return the created component, and can also return an error:
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
    if err := ioc.PutFactory(ComponentFactory, "MyDemoComponent", "DemoComponentAlias"); err != nil {
        panic(err)
    }
}
```

Factories can not define cyclic dependencies (i.e. a factory produces a component `A` which is needed to another factory to create a component `B` which should be injected in the factory of `A`). To resolve the problem, you have break the cycle, and wait another step in the component's lifecycle (like [Injection](#injection) or [Post-Initialization](#post-initialization)) to inject the required component.

Like explain in the section [Testing](#testing), the functions `Put` and `PutFactory` define the component in the core set. The functions `TestPut` and `TestPutFactory` can be used in the same way to define a component in the test set.

### Initialization

At some point in the application, probably triggered by the exploitation step, the instance of the component will be required. It will not be completely ready before the end step of the step post-initialization. This is how it is created and initialized:

#### Instantiation

Of course, the first step is to create an instance. If the component has been defined with the functions `Put` or `TestPut`, the instance is directly used. If the component has been defined with a factory, with the function `PutFactory` or `TestPutFactory`, the factory is called with its argument injected.

If a factory return a not-null error, the container stops all the process as soon as possible and return the error wrapped in some context messages.

#### Injection

Regardless of the registration mode (instance or factory), if the component is a structure or a pointer to a structure, it is injected. The process is defined in the section [Injection of structures](#injection-of-structures), with the mode "tagged-only".

At this point, cyclic dependencies can be resolved, but that means that maybe the injected component is not fully initialized (only in case of a cycle: if there is no cycle dependency, the injected component should be ready to be used).

#### Post-Initialization

If the component have a method `PostInit`, the container calls it. The method is injected, like described in [Injection of functions](#injection-of-functions). The method should return nothing or an error. Like factories, if one `PostInit` returns a not-null error, the container stops everything and returns the error wrapped in context messages.

### Exploitation

Only defining components doesn't create anything. You have to specify what are the main components, the components you need to instantiate and get to starting your application (or doing some tests). The container will instantiate these components, and also their dependencies. The API defines the function `CallInjected` that can be used to retrieve that main components from the container.
```go
type MyMainComponent struct { ... }

func (self *MyMainComponent) start() { ... }

func init() {

    if err := ioc.Put(&MyMainComponent{}, "MyMainComponent"); err != nil {
        panic(err)
    }

}

func main() {

    err := ioc.CallInjected(func(component *MyMainComponent) {
        component.start()
    })

    if err != nil {
        panic(err)
    }

}
```
As you should have guessed, the function `CallInjected` take as sole argument a function that will be injected (see [Injection of functions](#injection-of-functions)). The given function can return nothing or an error. The method `CallInjected` return an error if something had happened.

### Destruction

The container can close automatically every component implementing the interface `io.Closer`. To cleanly close the container, you have to call the function `Close`:
```go
type MyMainComponent struct { ... }

func (self *MyMainComponent) start() { ... }

func (self *MyMainComponent) close() error { ... }

func init() {

    if err := ioc.Put(&MyMainComponent{}, "MyMainComponent"); err != nil {
        panic(err)
    }

}

func main() {

    defer ioc.Close()

    err := ioc.CallInjected(func(component *MyMainComponent) {
        component.start()
    })

    if err != nil {
        panic(err)
    }

}
```
If a component returns a non-null error at the call of its method `Close`, the error is silently discarded. The `Close` methods are called in the same order than the components are instantiated (component without dependencies first, main components last).

Another way to close every components, and in the same time delete all component registration in the test set, is the function `ClearTests`. The function is essentially useful for unit test. Here an extract of unit tests with `Gomega`:

[pkg/main.go]
```go
package pkg

type SmallStuffDoer interface {
    doSmallStuff() error
}

// Component A

type ComponentA struct {}

func (self *ComponentA) doSmallStuff() error { ... }

func init() {
    if err := ioc.Put(&ComponentA{}, "ComponentA", "SmallStuffDoer"); err != nil {
        panic(err)
    }
}

// Component B

type ComponentB struct {
    SmallStuffDoer SmallStuffDoer `inject:""`
}

func (self *ComponentB) doBigStuff() error { ... }

func init() {
    if err := ioc.Put(&ComponentB{}, "ComponentB"); err != nil {
        panic(err)
    }
}

```

[pkg/main_test.go]
```go
package pkg_test

import (
    "testing"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
    . "pkg"
)

func TestIoc(t *testing.T) {
    RegisterFailHandler(Fail)
    RunSpecs(t, "App test suite")
}

type SmallStuffDoerMock struct {}

func (self *SmallStuffDoerMock) doSmallStuff() error { ... }

var _ = Describe("App tests", func() {

    AfterEach(func() {
        ioc.ClearTests()
    }

    Describe("Component B should work", func() {

        BeforeEach(func() {
            ioc.TestPut(&SmallStuffDoerMock{}, "Mock", "SmallStuffDoer")
        }

        It("should work with the mock", func() {

            var b *ComponentB
            ioc.CallInjected(func(injected *ComponentB) {
                b = injected
            })

            Expect(b.doBigStuff()).Should(Succeed())

        })

    })

})
```


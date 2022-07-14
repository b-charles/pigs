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

A container is a set of components: it manages the lifecycle of each component and take in charge the injection process. To know what component should be injected, the container look after names: each injection defines the name (main name or alias) of the component that should be injected. If there is a problem of type (the component injected doesn't correspond to the expected type), the container detects it and returns an error. It's the responsibility of the developer to choose wisely the names of the components to avoid collisions (the framework will help), and to define components and injections that are compatible.

The container is a managed component itself and can be accessed by the name `github.com/b-charles/pigs/ioc/Container`.

### Default naming process

Defining each name and alias through a string can be combersome and hazardous in case of refactoring. It's always possible to define the main name and aliases with strings, but the framework offers other options.

For main names of component or default names of injections, the framework can do some reflection to deference the type (if necessary) and use [the methods `PkgName` and `Name` of `reflect.Type`](https://pkg.go.dev/reflect#Type) to compute a default name, with the format `<pkgName>/<name>`. That name should be unique for a type, which means that if several components with the same type should be defined, the developper - you - has to name "manually" each problematic component.

Alias definition during component registration are more tricky: aliases are about the interfaces implemented by the component[^duck], and Go doesn't provide a way to access to a type as a value without having a value of this type, like [the `.class` syntax in Java](https://docs.oracle.com/javase/tutorial/reflect/class/classNew.html) or direct access in [Python](https://docs.python.org/3/reference/compound_stmts.html#class-definitions) or [Ruby](https://ruby-doc.org/docs/ruby-doc-bundle/Manual/man-1.4/syntax.html#const). (I know it's not very fair to compare Go with interpreted langages, but it's more to show what we need.) So, the best way found (so far) is to use a function with anonymous inputs, without output and with an empty body:
```go
func(SomeInterface){}
```
If this function is given in place of an alias, the framework analyzes the inputs types and, like the main name, use the method `PkgName` and `Name` of `reflect.Type` to generate an alias. Several aliases can be given in one function, simply use more input arguments:
```go
func(SomeInterface, AnotherInterface){}
```

[^duck]: The concept of aliasing is not very in line with the spirit of the duck typing of Go, but it's a need for performance (what should be injected, without checking each component) and flexibility (what should *not* be injected, even if the interface matches the implentation).

### Injections

The framework is based on two types of injection: injection of structures and injection of functions.

#### Injection of structures

At different points in the framework, it appears the need to inject a structure, or a pointer to a structure. Two flavors are possible: "tagged only" and "all injected".

For the mode "tagged only", only the tagged `inject` fields are injected. The value of the tag defines the name (or the alias) of the component to inject. If the value of the tag is empty, a default name is computed, with the process explained at [Default naming process](#default-naming-process). If no component is defined by this name, the name of the field is used. So if this variable `injected` is injected:
```go
package mypkg

type MyInjectedStruct struct {
    NotInjected *NotInjected
    FirstInject SomeInterface `inject:"Something"`
    SecondInject AnotherInterface `inject:""`
}

var injected *MyInjectedStruct = &MyInjectedStruct{}
```
the first field `NotInjected` is not injected, the second field `FirstInject` is injected with a component named (or aliased) `Something`, and the third field is injected with a component named (or aliased) `<path of mypkg module>/mypkg/AnotherInterface` or, if no component is found, with a component named or aliased `SecondInject`. Please note that each injected field should be settable (so exported, with a name beginning with an upper case).

For the mode "all injected", every field should be injected, and the tag `inject` can be used only to define the correct name of the component to inject.

#### Injection of functions

At other points of the framework, the container use some user-defined functions, where the arguments need to be injected. Two prototypes of function can be injected, and indifferently used: with each argument injected, or with an injected structure.

The first form is when the function is defined with none, one or several arguments, each expected to be injected. In that case, the name of the type of each argument is used as the name (or alias) of the component to inject:
```go
func InjectedFunc(first mypkg.SomeInterface, second *anotherlib.RandomComponent) { ... }
```
If the container calls the function, the first argument will be injected with a component named `<path of mypkg module>/mypkg/SomeInterface` (see [Default naming process](#default-naming-process)) and the second argument with a component named `<path of anotherlib module>/anotherlib/RandomComponent`: on the contrary of structure injection, note that there is no fallback on the name of the argument. Of course, the function can be defined without any argument.

This first form is pretty clear and should be sufficient in most cases, but you can't define arbitrary injection names. The other is little more messy, but more powerful: the function is defined with only one argument, a structure or a pointer to a structure, which will be injected as defined in the section [Injection of structures](#injection-of-structures), with the flavor "all injected". The type of this sole argument can be properly defined or defined directly in the function:
```go
func injectedFunc(injected struct {
    SomeComponent mypkg.SomeInterface
    Another RandomInterface `inject:"RandomComponent"`
}) { ... }
```
Here, the argument `injected` will be injected with a component named `<path of mypkg module>/mypkg/SomeComponent` and another one named `RandomComponent`.

#### Auto-discovery injection

Whether it is an injection in a structure or in a function, if no component or too many components are defined, the container returns an error. But you can inject all components sharing an alias: that's the foundation of the auto-discovery. If the injected value is of type slice, the container will start by trying to resolve the injection normally (a component can be defined as a slice). If no component match the definition (i.e. if there is several components found with the requested name or the found component is not assignable to the injected value), the container will create a slice, add in it all the components with an alias matching the name and inject that slice:
```go
type MyInjectedStruct struct {
    AllComponents []SomeInterface `inject:"MagicService"`
}

var injected *MyInjectedStruct = &MyInjectedStruct{}
```
Here, if `injected` is injected by the container, the field `AllComponents` will contains every components defined with the alias `"MagicService"` (at the condition that all respect the interface `SomeInterface`). The injected slice can not be empty: at leat one component should be found, or the framework returns an error[^empty].

[^empty]: This behaviour, throwing an error if an "auto-discovered" slice is empty, may seem inconvenient, but some fallback mecanism (test/core or struct field type/struct field name) are more controllable and intuitively expected with that rule. Just make sure every auto-discovered slice has at least one component definition, implementing a no-op if necessary.

With the same spirit, you can also inject a map of components sharing a same alias, where the keys (of type `string`) are the main names of the associated components:
```go
type MyInjectedStruct struct {
    AllComponents map[string]SomeInterface `inject:"MagicService"`
}

var injected *MyInjectedStruct = &MyInjectedStruct{}
```
Again, the injected map can not be empty, and at least one component with that alias should be found.

#### Testing

Like it's said in the main goals of the framework, there is no concept of scope, or at least not extendable scope. In fact, the container is divided in two set of components: one for the core components, and one for tests. These two sets are comparable to scopes: every time an injection is proceeded, the container start by searching any component in the test set. If nothing is found, then the container use the core set. In case of auto-discovery injection (slice or map), if at least one matching component is found in the test set, only that components are injected and the components in the core set are not used.

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

You can call the functions `ioc.PutNamed` to record directly the component, its main name and optionally its aliases:
```go
type DemoComponent struct { ... }

func init() {
    ioc.PutNamed(&DemoComponent{}, "MyDemoComponent", "DemoComponentAlias")
}
```
Recording a component doesn't inject it (not yet). It is only recorded as a not yet initialized component. The function `Put` panics if something bad had happened (e.g. a component with the same main name is already registered). Aliases can be defined with strings, or with one or several weird functions, as explained in [Default naming process](#default-naming-process).

You can also use the function `ioc.Put` to record the component and let the framework to choose the main name (again, see [Default naming process](#default-naming-process)):
```go
type DemoComponent struct { ... }

func init() {
    ioc.Put(&DemoComponent{})
}
```

Instead of `ioc.Put` and `ioc.PutNamed`, you can also use `ioc.PutFactory` and `ioc.PutNamedFactory` to define a factory, a function which will create the component:
```go
package mypkg

type DemoComponent struct {}

func ComponentFactory() *DemoComponent {
    return &DemoComponent{}
}

func init() {
    ioc.PutNamedFactory(ComponentFactory, "MyDemoComponent", "DemoComponentAlias")
    // or
    ioc.PutFactory(ComponentFactory);
}
```
The function `PutFactory` use the type of the first output to compute a main name (here, it should be `mypkg.DemoComponent`). Recording a factory doesn't call it, but it will be used to instantiate the component in a next step. At that time, the arguments of the factory will be injected like it is described in the section [Injection of functions](#injection-of-functions). The factory should at least return the created component, and can also return an error:
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

Like explain in the section [Testing](#testing), the functions `Put`, `PutNamed`, `PutFactory` and `PutNamedFactory` define the component in the core set. The functions `TestPut`, `TestPutNamed`, `TestPutFactory` and `TestPutNamedFactory` can be used in the same way to define a component in the test set.

Finaly, all this method have their `Erroneous` prefixed version (e.g. `ErroneousTestPutNamed`) which doesn't panic but returns an error if something wrong happened.

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

During the call of `CallInjected`, the instance of some components will be required. This is how each component is created and initialized:

##### Instantiation

Of course, the first step is to create an instance. If the component has been defined with the functions `Put`, `PutNamed`, `TestPut` or `TestPutNamed`, the instance is directly used. If the component has been defined with a factory, with the function `PutFactory`, `PutNamedFactory`, `TestPutFactory or `TestPutNamedFactory`, the factory is called with its argument injected.

If a factory return a not-null error, the container stops all the process as soon as possible and return the error wrapped in some context messages.

##### Injection

Regardless of the registration mode (instance or factory), if the component is a pointer to a structure, it is injected. The process is defined in the section [Injection of structures](#injection-of-structures), with the mode "tagged-only".

At this point, cyclic dependencies can be resolved.

##### Post-Initialization

If the component have a method `PostInit`, the container calls it. The method is injected, like described in [Injection of functions](#injection-of-functions). The method should return nothing or an error. Like factories, if one `PostInit` returns a not-null error, the container stops everything and returns the error wrapped in context messages.

#### Run main function

When each necessary components are fully initialized, the framework call the given function of `CallInjected`. That's what we wanted from the start and where your business begins. Be carefully if you are working with multi-threads: the end of this function will trigger the next phase of the container and so the closing of the components.

#### Close

After the main function executed, the container will close automatically every component implementing the interface `io.Closer`. If a component panics or returns a non-null error at the call of its method `Close`, the error is silently discarded. The `Close` methods are called in the same order than the components are instantiated (component without dependencies first, main components last).

### Redefinition

It should be only one call of `CallInjected` during the run of the application or at each unit test: after a call of `CallInjected`, every instances are released along the core component definitions if no test component are found, otherwise instances and only test component definitions are released if at least one component is defined in the test scope.

So, if you are running unit tests, you have to defined some fixture to redefined each time all the test components you want to use, and for each test, every component is re-instanced. Be sure to define at least one component in the test scope (even if it is not used and not instanciated) before each call of `CallInjected` or you will loosing the definitions of all the core components.

If no test component is defined, the framework considers that you are really running your application. In this case, in order to consume the least RAM as possible, all unused component instances and all component definitions will be released (forgotten by the framework) and can be deleted by the garbage collector if no other component reference it.

### Container lifecycle awareness

You can define some special components which will be called during the exploitation phase (`CallInjected`). That components should implement the desired interface and be declared with the correct alias:
 * `github.com/b-charles/pigs/ioc/PreCallAwared`: Any component with that alias (and so implementing the `PreCallAwared` interface) will be call at the beginning of the `CallInjected` method.
 * `github.com/b-charles/pigs/ioc/PostInstAwared`: Components with that alias will be called in `CallInjected`, after the instanciation and initialisation of the arguments of the given method, but before actually calling it.
 * `github.com/b-charles/pigs/ioc/PreCloseAwared`: Components with that alias are called after the execution of the function given to `CallInjected`.
 * `github.com/b-charles/pigs/ioc/PostCloseAwared`: Components with that alias are called at the end of the `CallInjected` method.

## Usage expample

To a better understanding of how the framework can be used, here an extract of unit tests with `Ginkgo`:

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
    ioc.Put(&ComponentA{}, func(SmallStuffDoer){})
}

// Component B

type ComponentB struct {
    SmallStuffDoer SmallStuffDoer `inject:""`
}

func (self *ComponentB) doBigStuff() error { ... }

func init() {
    ioc.Put(&ComponentB{})
}

func main() {
    ioc.CallInjected(func(b *ComponentB) error {
        return b.doBigStuff()
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
    . "pkg"
)

func TestIoc(t *testing.T) {
    RegisterFailHandler(Fail)
    RunSpecs(t, "App test suite")
}

type SmallStuffDoerMock struct {}

func (self *SmallStuffDoerMock) doSmallStuff() error { ... }

var _ = Describe("App tests", func() {

    Describe("Component B should work", func() {

        BeforeEach(func() {
            ioc.TestPut(&SmallStuffDoerMock{}, func(SmallStuffDoer){})
        }

        It("should work with the mock", func() {

            ioc.CallInjected(func(b *ComponentB) {
                Expect(b.doBigStuff()).Should(Succeed())
            })

        })

    })

})
```


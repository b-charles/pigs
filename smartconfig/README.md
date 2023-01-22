# Pigs - Smart config

Retrieve the configuration, but in a smart way.

## Why?

You are using the [config](../config) module, but it's the twentieth times you've used `strconv.Atoi` to cast a string into your config struct? You are at the right place.

This module is a sleight of hand: you give a prefix and a blank struct to the `Configure` function, and the struct is filled with the values from the config module and you can use it from the ioc framework as an injectable component. The fields of the struct can be something other than strings, the module defines and uses parsers to convert strings into the right type.

## How it works?

### Parsers

First, before any configuration, the module get all `Parser`s components.

A `Parser` is a function taking a string in input and returning the parsed value and an error:
```go
type Parser func(string) (T, error)
```
The module use the return type to know when the parser should be use. The functions `strconv.Atoi` and `strconv.ParseBool` are already defined as `Parser`s.

### Configurers

Parsers are kind but not very powerful. Sometimes you need to explore what configuration keys are available to construct the final value. The module defines the more flexible interface `Configurer`:
```go
type Configurer interface {
  Target() reflect.Type
  Configure(NavConfig, reflect.Value) error
}
```

The method `Target` returns what type the configurer is able to parse to. The method `Configure` does the job: the interface `NavConfig` is a wrapper of `config.Configuration` with some methods to _navigate_ in the keys (see below). The given `reflect.Value` is the receipient of the result of the parsing and the configurer should returns its parsing by using [the method `Set`](https://pkg.go.dev/reflect#Value.Set) of this value.

The `NavConfig` interface represents the configuration as a tree. Each keys present in `config.Configuration` are split by the char `.` and each parts corresponds to a node in the tree. The interface defines methods to move around: get the current value, get the available sub keys, get the parent, a child node...
```go
type NavConfig interface {
  Root() NavConfig
  Parent() NavConfig
  Path() string
  Value() string
  Keys() []string
  Child(string) NavConfig
  Get(string) NavConfig
}
```
The root node corresponds to the key `""` (empty string). The method `Get` differs of the method `Child` by splitting the given key first and applies the method `Child` on each parts of the key. If the input key starts with a `.`, the process is applied from the root instead of the current node.

### Special configurers

#### Struct

The module handle structs and pointers to struct. Of course, the input of the method `Configure` should be a pointer so it has to be settable, and each field name should starts with an upper case for the same reason. Each field can be annoted with the tag `config` wich can defined the key to use (relative or absolute, see [the `Get` method of `NavConfig`](#configurers)). If no tag are found, the field name in lowercase is used as a relative sub key. The configurer search for each field which parser or configurer has to be used and can be called recursively.

#### Slices

The module also handles slices. Available keys are sorted, integer numbers in numerical order first, then other keys in lexicographical order, and the same parser or configurer is used for each sub key found.

#### Maps

Finally, maps with string keys (`map[string]T`) are also supported. Not much to say: it just works.

### The `Configure` function

The `Configure` function take two inputs:
```go
func Configure(root string, configurable any)
```

 * A root path: for each configured field, the corresponding configuration key will be prefixed by this root path (except if the field is annoted with the tag `config` and the value of the tag starts with an `.`). This root path can be `""` (empty string).
 * A pointer to a struct, which will be setted by the special configurer for struct. The function will register it as an ioc component.

Two other functions, `DefaultConfigure` and `TestConfigure` are also defined to register the configuration struct in the default scope and in test scope of the ioc framework.


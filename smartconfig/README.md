# Pigs - Smart config

Retrieve the configuration, but in a smart way.

## Why?

You are using the [config](../config) package, but it's the twentieth times you've used `strconv.Atoi` to cast a string into your config struct? You are at the right place.

This package is a sleight of hand: you give a prefix and a blank struct to the `Configure` function, and the struct is filled with the values from the config package and you can use it from the ioc framework as an injectable component. The fields of the struct can be something other than strings, the package defines and uses parsers to convert strings into the right type.

## How it works?

### Parsers

First, before any configuration, the package get all `Parser`s components.

A `Parser` is a function taking a string in input and returning the parsed value and an error:
```go
type Parser func(string) (T, error)
```
The package use the return type to know when the parser should be used.

Default parsers are defined following the [3 steps tricks to register overloadable components](../ioc/README.md#overloadable-components-in-auto-discovery-injection). So, to overwrite a default implementation, you have to register a component implementing the specific interface in the core or test scope. The default parsers are:
| type | signature | comment |
| --- | --- | --- |
| `string` | `type StringParser func(string) (string, error)` | identity function |
| `float64` | `type Float64Parser func(string) (float64, error)` | based on `strconv.ParseFloat(string, 64)` |
| `int` | `type IntParser func(string) (int, error)` | based on `strconv.Atoi` |
| `bool` | `type BoolParser func(string) (bool, error)` | based on `strconv.ParseBool` |

### Inspectors

Parsers are kind but not very powerful. Sometimes you need to explore what configuration keys are available to construct the final value. The package defines the more flexible type `Inspector`:
```go
type Inspector func(NavConfig) (T, error)
```
Like for `Parser`, the function returns what type it is able to parse to. The main difference is that the function takes a `NavConfig` as input: `NavConfig` is an interface representing the configuration as a tree. Each keys present in `config.Configuration` are split by the char `.` and each parts corresponds to a node in the tree. The interface defines methods to move around: get the current value, get the available sub keys, get the parent, a child node...
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

The package handle structs and pointers to struct. Of course, the input of the method `Configure` should be a pointer so it has to be settable, and each field name should starts with an upper case for the same reason. Each field can be annotated with the tag `config` which can defined the key to use (relative or absolute, see [the `Get` method of `NavConfig`](#inspectors)). If no tag are found, the field name in lowercase is used as a relative sub key. The configurer search for each field which parser or configurer has to be used and can be called recursively.

#### Slices

The package also handles slices. Available keys are sorted, integer numbers in numerical order first, then other keys in lexicographical order, and the same parser or configurer is used for each sub key found.

#### Maps

Finally, maps with string keys (`map[string]T`) are also supported. Not much to say: it just works.

### The `Configure` function

The `Configure` function take two inputs:
```go
func Configure(root string, configurable any)
```

 * A `root` path: for each configured field, the corresponding configuration key will be prefixed by this root path (except if the field is annotated with the tag `config` and the value of the tag starts with an `.`). This root path can be `""` (empty string).
 * A pointer to a struct, which will be setted by the special configurer for struct. The function will register it as an ioc component.

Two other functions, `DefaultConfigure` and `TestConfigure` are also defined to register the configuration struct in the default scope and in test scope of the ioc framework.


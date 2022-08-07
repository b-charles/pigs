# Pigs - Config

Configure the application.

## What's that?

Hard-coding the configuration in your code is not a good idea. Many people have tried and they are unanimous: don't do it.

There is different approach to configure your application: command line arguments, environment variables, configuration files in different format (properties, json, yaml, toml, ...) ... This module can be used to get all theses sources of configuration and merge them in an simple and [injectable](../ioc/README.md) component.

## What's supported?

For now, the framework supports command line arguments, environment variables and some utilities to define default values. But adding other sources of configuration can be easily done.

## How it works?

The module defines two major structures: `ConfigSource` to handle a source of configuration, and `Configuration` to store the final merge of all sources.

Each `ConfigSource` should be declared in the IOC framework with the alias `"ConfigSource"`, and the `Configuration` component can be retrieved by the name `"Configuration"`.

### The `ConfigSource`s

A `ConfigSource` is a component which implements this simple interface:

```go
type ConfigSource interface {
  GetPriority() int
  LoadEnv() map[string]string
}
```

The method `LoadEnv` returns the variable defined by the source. The priority returned by `GetPriority` is used to sort the sources: if two sources define two values for the same variable, the source with the greater priority will overload the other.

By simply defining a dependency on the config module will defines some default sources.

#### Default configuration source

This configuration source is used to defined programmatically default values. It's defined with the priority `-999`, in the module `github.com/b-charles/pigs/config/confsources/conf`.

To define a default value, simply use the method `Set` in an init function:
```go
package mypkg

import "github.com/b-charles/pigs/config"

func init() {

  config.Set("my.little.variable", "Greatest Value Ever!")
  config.Set("another.critical.variable", "Foo Bar")

}
```

The module can also be used to define some values for unit tests. When using the function `SetTest`, a special component is created using the given map as a configuration source and defined in the test scope. That component will be automatically deleted after the call of `ioc.CallInjected`, so the method should be called in a fixture (before tests) of unit tests.

#### Environment variables

This configuration source handles environment variables. It's defined with the priority `0`.

Before being returned for merging, each variable name is processed:
 * all characters are casted to lower case,
 * all characters `_` are replaced by a `.`.

#### Command line arguments

Command line arguments are also processed with a priority `100`.

Arguments can be in one of theses format (with `[name]` the name of the variable, and `[value]` its value):
* `--[name]=[value]` (e.g.: `--music=Jimmy`)
* `--[name]='[value]'` (e.g.: `--music='Summer Vibe'`)
* `--[name]="[value]"` (e.g.: `--music="Via con me"`)
* `--[name]` for boolean value, which will be associated with the value `true`
* `--no-[name]` for boolean value, which will be associated with the value `false`.

### The `Configuration`

The `Configuration` component manages the merging of all sources in this component, and resolve placeholders: for each value, each occurance of the pattern `${<myvalue>}` is replaced with the value of `<myvalue>`. So, if a config source defines a value `name` with `Batman` and another value `whoami` with `I'm ${name}`, the resolving process will convert `whoami` to `I'm Batman`. Placeholders can be chained and nested:
| name | value | resolved |
| --- | --- | --- |
| `ironman` | `Tony Stark` | `Tony Stark` |
| `super` | `${ironman}` | `Tony Stark` |
| `best` | `${super}` | `Tony Stark` |
| `what` | `iron` | `iron` |
| `who` | `${${what}man}` | `Tony Stark` |

## And now?

The usage of this module should be efficient but not convenient. The module [smartconfig](../smartconfig/README.md) can be useful to get a chunck of typed configuration values.


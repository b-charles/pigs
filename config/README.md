# Pigs - Config

Configure the application.

## What's that?

Hard-coding the configuration in your code is not a good idea. Many people have tried and they are unanimous: don't do it.

There is different approach to configure your application: command line arguments, environment variables, configuration files in different format (properties, json, yaml, toml, ...) ... This package can be used to get all theses sources of configuration and merge them in a simple and [injectable](../ioc/README.md) component.

## What's supported?

For now, the package supports command line arguments, environment variables, json files and some utilities to define default values. But adding other sources of configuration can be easily done.

## How it works?

The package defines two major structures: `ConfigSource` to handle a source of configuration, and `Configuration` to store the final merge of all sources.

### The `ConfigSource`s

`ConfigSource` is an interface, a component signature defined by:

```go
type ConfigSource interface {
  GetPriority() int
  LoadEnv(MutableConfig) error
}

type MutableConfig interface {
  HasKey(string) bool
  Keys() []string
  GetRaw(string) (string, bool)
  Lookup(string) (string, bool, error)
  Get(string) string
  Set(string, string)
}
```

The priority returned by `GetPriority` is used to sort the sources: sources with lower priority will be called first and so sources with greater priority will be able to override values. The method `LoadEnv` records the variable defined by the source in the `MutableConfig` by the method `Set`. The `MutableConfig` other methods can be used to have a partial configuration and use the other sources to configure this source. The methods `Lookup` and `Get` returns value with resolved placeholders (see the section [The `Configuration` component](#the-configuration-component)).

The package defines some default sources.

#### Default configuration source

The default configuration source is build-in and is not based on the `ConfigSource` interface. The default values are always loaded first.

To define a default value, simply use the method `Set(key, value string)` or `SetMap(values map[string]string)` in an init function:
```go
package mypkg

import "github.com/b-charles/pigs/config"

func init() {

  config.Set("my.little.variable", "Greatest Value Ever!")
  config.SetMap(map[string]string{
    "you.great.entry":           "500",
    "another.critical.variable": "Foo Bar",
  })

}
```

Default value should be defined once and can not be redefined.

#### Test configuration source

The package can also be used to define some values for unit tests. For each call of the functions `Test(key, value string)` or `TestMap(values map[string]string)`, a special component is created and defined in the test scope. Each time a component is created, its priority is increased, meaning you can override a configuration value:
```go

  config.TestMap(map[string]string{
    "my.little.variable":        "Any banal value",
    "another.critical.variable": "Foo Bar",
  })

  // override my.little.variable
  config.Test("my.little.variable", "Greatest Value Ever!")

```

Since they are registered in the test scope, this components will be automatically deleted after the call of `ioc.CallInjected`, so the configuration should be done in a fixture (before tests) of unit tests. This also means that using these functions will discard any other configuration source defined in the core scope (except the default configuration source since it's build-in).

#### Environment variables

The package offers a default `ConfigSource` which loads environment variables.

This default config source defines its priority at `0`. It modifies each environment variable name:
 * all characters are converted to lower case,
 * all characters `_` are replaced by `.`.

The registration of this default component follows the [classical trick to register overloadable components](../ioc/README.md#overloadable-components-in-auto-discovery-injection). So, if the default implementation doesn't suit you, you can defines a component with the signature `type EnvVarConfigSource ConfigSource` in the core or the test scope of the ioc framework.

#### Command line arguments

The package also offers a default `ConfigSource` to process command line arguments.

This default config source is defined with a priority `100`. Arguments can be in one of theses format (with `[name]` the name of the variable, and `[value]` its value):
* `--[name]=[value]` (e.g.: `--music=Jimmy`)
* `--[name]='[value]'` (e.g.: `--music='Summer Vibe'`)
* `--[name]="[value]"` (e.g.: `--music="Via con me"`)
* `--[name]` which will be associated with the value `true`
* `--no-[name]` which will be associated with the value `false`.
Arguments with only one starting dash (e.g.: `-music=Losers`) are also accepted.

Like for the environment variable config source, the default implementation is registered in the ioc framework following [classical schema to register overloadable components](../ioc/README.md#overloadable-components-in-auto-discovery-injection). So you can replace the default component by registering a component with the signature `type ArgsConfigSource ConfigSource` in the core or test scope.

#### Json files

The package defines a default `ConfigSource` to process Json files.

This default config source is defined with the priority `200` and gets any configuration key previously defined starting with `config.json` and corresponding to an existing file. Then each file will be loaded, parsed and integrated in the configuration. The path separator of the file path should always be `/` and the path can be absolute (starting with an `/`) or relative to the working directory.

Again, like for environment variables and command line arguments, the integration of this default config source is done with [the 3 steps to register overloadable components](../ioc/README.md#overloadable-components-in-auto-discovery-injection). To replace the default component, you can simply register a component with the signature `type JsonFilesConfigSource ConfigSource` in the core or the test scope.

### The `Configuration` component

The `Configuration` component manages the merging of all sources, and expose the result as an injectable component:

```go
type Configuration interface {
  HasKey(string) bool
  Keys() []string
  GetRaw(string) (string, bool)
  Lookup(string) (string, bool, error)
  Get(string) string
}
```

The default implementation of this interface is registered in the default scope and can be overridden.

In the default implementation, the methods `Lookup` and `Get` resolve placeholders: for each value, each occurrence of the pattern `${<myvalue>}` is replaced with the value of `<myvalue>`. So, if a config source defines a value `name` with `Batman` and another value `whoami` with `I'm ${name}`, the resolving process will convert `whoami` to `I'm Batman`. Placeholders can be chained and nested:
| name | value | resolved |
| --- | --- | --- |
| `ironman` | `Tony Stark` | `Tony Stark` |
| `super` | `${ironman}` | `Tony Stark` |
| `best` | `${super}` | `Tony Stark` |
| `what` | `iron` | `iron` |
| `who` | `${${what}man}` | `Tony Stark` |

## And now?

The usage of this package should be efficient but not convenient. The package [smartconfig](../smartconfig/README.md) can be useful to get a chunk of typed configuration values.


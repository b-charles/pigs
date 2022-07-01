# Pigs - Config

Configure the application.

## What's that?

Hard-coding the configuration in your code is not a good idea. Many people have tried and they are unanimous: don't do it.

There is different approach to configure your application: command line arguments, environment variables, configuration files in different format (properties, json, yaml, toml, ...) ... This module can be used to get all theses sources of configuration and merge them in an simple and [injectable](../ioc/README.md) `map[string]string`.

## What's supported?

[The Twelve-Factor App](https://12factor.net/config) manifesto is pretty specific about how it should be done: you should use environment variables. We choosed to add the support of command line arguments, and some utilities to define default values, but nothing else.

But adding other sources of configuration can be easily done.

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

To define a default value, simply use the method `SetDefault` in an init function:
```go
package mypkg

import "github.com/b-charles/pigs/config/confsources/conf"

func init() {

  conf.SetDefault("my.little.variable", "Greatest Value Ever!")
  conf.SetDefault("another.critical.variable", "Foo Bar")

}
```

The module can also be used to define some values for unit tests. When using the function `SetEnvForTests`, a special component is created using the given map as a configuration source and defined in the test scope. The method should be called in the fixture (before tests) of unit tests, and the method `ioc.ClearTests` should be called to clear the definitions.

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

The `Configuration` component is simply defined as a map of string to string:
```go
type Configuration map[string]string
```

The module manages the merging of all sources in this component, and application needing to be configured can retrieve it by its name `Configuration`.

Of course, modifying the component is not a good idea: define and use a dedicated `ConfigSource` instead.

## And now?

The usage of this module should be efficient but not convenient. The module [smartconf](../smartconf/README.md) adds the concept of `Parser` to convert a string to something else, and can be useful to get a chunck of typed configuration values.


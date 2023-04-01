
# Log

It's time to log.

## Why again?

There is a lot of logging librairies available in the Go ecosystem ([build-in](https://pkg.go.dev/log), [Zap](https://github.com/uber-go/zap), [Zerolog](https://github.com/rs/zerolog), [Logrus](https://github.com/sirupsen/logrus) ...), so, you can wonder, why oh why an other logger framework?

The primal interest of this logging solution is it's integrated with the other modules: [ioc](../ioc/README.md), [json](../json/README.md) and [smartconfig](../smartconfig/README.md). So, in one hand, you have to fully understand theses dependency modules to understand how to use this library, and in the other, you have access to a fully integrated solution, working on the same basis than the rest of your application.

A logging library is always very opinionated: what should be logged (a message, some structured data...), how (in the standard output, in a file, with or without rotation...), which features should be available (multiple loggers, additional context, hooks...)... This library is no exception to this rule, but tries to keep the API simple and fully extendable, with all the structural components in the default scope so you can overload and defines your own logic and features.

## Quick start

By default, a logger instance is already defined and can be used by any component:
```go

type MyComponent struct {
  Logger log.Logger `inject:""`
  ...
}

func (my *MyComponent) doSomething( input string ) error {

  ...

  my.Logger.Info()
    .Set("what", "I did something")
    .Set("with", input)
    .Log();

  return nil

}

func init() {
  ioc.Put(&MyComponent{})
}

```

## Longer start

To be as extendable as possible, this library introduces different concepts and interfaces, and each interface comes with an default implementation and a default instance in the ioc default scope (so which can be overwritten).

We will see in detail the `Logger` interface, the main interface of the library, but for now, it's enough to simply say that a loggers are immutable and defined with:
 * a name, a string which can be structured with dots (`.`) to reflect a hierarchy of the loggers: the logger `first` is the parent of `first.second`.
 * a minimum level, which can be choose from `Trace`, `Debug`, `Info`, `Warn`, `Error` and `Fatal`.

### The level configurer

The default implementation is build with the help of a `LevelConfigurer`, a component which is dedicated to provide the minimum level of a logger by its name:
```go
type LevelConfigurer interface {
  GetLevel(string) Level
}
```

A default implementation is provided and registered in the ioc default scope: the name is casted to lower case, then used as suffix of `log.level.` to search in the configuration (see [config](../config/README.md) and [smartconfig](../smartconfig/README.md)). If no level is defined, the parents are used successively, until a level is found or no parent is left. At the end of the process, if no level is found, the level `Info` is used.

### Appenders

The appenders are used in the default implementation to output the log somewhere. The log is provided as a Json node (see [json](../json/README.md)), and should not return any error:
```go
type Appender interface {
  Append(json.JsonNode)
}
```

A default implementation writing the node in the standard output is recorded in the ioc default scope.

### Contextualizers

The interface `Contextualizer` is defined to add some context, i.e. some values, in the in-progess log. It's definition is linked with the `Logger` and `LoggerBuilder` interfaces:
```go
type Contextualizer interface {
  AddContext(Logger, Level, LogBuilder)
}
```

A default contextualizer is recorded in the ioc default scope, with type `*DefaultContextualizer`: it's purpose is to add the date of the log creation and the level of the log.

A type `StaticContextualizer` is also defined to quickly convert a `map[string]any` map to a contextualizer.

### The logger factory

The `LoggerFactory` is an interface (with a provided default implementation) which can be used to create a `Logger` by it's name:
```go
type LoggerFactory interface {
  NewLogger(name string) Logger
}
```

The default implementation is injected with a `LevelConfigurer` and the available `Contextualizer`s and `Appender`s (see the auto-discovery feature of the ioc framework). Then, when a new logger is requested, the factory return a fully configured `Logger`.

### The logger

Finally, the main interface: the `Logger` interface can be used to log messages.

The usage of the logger interface sould be straight forward. To log a message, the interface returns a `LogBuilder` which can be used to populate each field of the message, in the same way than the Json builder (see the [json](../json/README.md) library). Don't forget to flush the message building by calling the `Log` function.

Loggers are immutable but a `Contextualizer` can be added by creating a new logger from an existing one.

A default logger is registered in the ioc default scope, by calling the `LoggerFactory` with the name `root`. A default level for this level is also defined by configuring (see the default configuration source of [config](../config/README.md)):
 * the key `log.level` is defined with `Info`.
 * the key `log.level.root` is defined with `${log.level}`.


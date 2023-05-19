
# Log

It's time to log.

## Why again?

There is a lot of logging libraries available in the Go ecosystem ([build-in](https://pkg.go.dev/log), [Zap](https://github.com/uber-go/zap), [Zerolog](https://github.com/rs/zerolog), [Logrus](https://github.com/sirupsen/logrus) ...), so, you can wonder, why oh why an other logger framework?

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
 * a minimum level.

### Levels

Levels can be choose from `Trace`, `Debug`, `Info`, `Warn`, `Error` and `Fatal`.

Integration with other libraries are proposed, following the [principle to register overloadable components](../ioc/README.md#overloadable-components-in-auto-discovery-injection):
 * A json marshaller is defined in the default scope, and can be overloaded by registering a component with the signature `type LevelMarshaller func(Level) (json.JsonNode, error)`.
 * A json unmarshaller is also defined in the default scope and can be overloaded by a component with the signature `type LevelUnmarshaller func(json.JsonNode) (Level, error)`.
 * Finally, a smart config parser is defined in the default scope and can be overloaded by a component with the signature `type LevelParser func(string) (Level, error)`.

### The level configurer

The default implementation is build with the help of a `LevelConfigurer`, a component which is dedicated to provide the minimum level of a logger by its name:
```go
type LevelConfigurer interface {
  GetLevel(string) Level
}
```

A default implementation is provided and registered in the ioc default scope: the name is changed to lower case, then used as suffix of `log.level.` to search in the configuration (see [config](../config/README.md) and [smartconfig](../smartconfig/README.md)). If no level is defined, the parents are used successively, until a level is found or no parent is left. The key `log.level` is by default defined with `Info`.

### Appenders

The appenders are used in the default implementation to output the log somewhere. The log is provided as a Json node (see [json](../json/README.md)), and should not return any error:
```go
type Appender interface {
  Append(json.JsonNode)
}
```

A default implementation writing the node in the standard output is recorded in the ioc default scope. The registration is [**not** done with the famous three steps shenanigan](../ioc/README.md#overloadable-components-in-auto-discovery-injection), so if you register any `Appender` component in the core or test scope, this default component will be ignored.

### Contextualizers

The interface `Contextualizer` is defined to add some context, i.e. some values, in the in-progress log. It's definition is linked with the `Logger` and `LoggerBuilder` interfaces:
```go
type Contextualizer interface {
  GetPriority() int
  AddContext(Logger, Level, LogBuilder)
}
```

During the log creation, the contextualizers are sorted by the integer returned by the `GetPriority` method and called before the log building.

A default contextualizer is recorded in the ioc default scope [by the method to register overloadable components](../ioc/README.md#overloadable-components-in-auto-discovery-injection), with signature `type DateLevelContextualizer Contextualizer`. With a priority `0`, it's purpose is to add the date of the log creation and the level of the log.

A type `StaticContextualizer` implementing `Contextualizer` is defined with to constructors `NewStaticContextualizer(key string, value any) StaticContextualizer` and `NewStaticContextualizerMap(m map[string]any) StaticContextualizer` to easily create contextualizers.

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

The usage of the logger interface should be straight forward. To log a message, the interface returns a `LogBuilder` which can be used to populate each field of the message, in the same way than the Json builder (see the [json](../json/README.md) library). Don't forget to flush the message building by calling the `Log()` function.

Loggers are immutable but a `Contextualizer` can be added by creating a new logger from an existing one with the methods `AddContextualizer(Contextualizer) Logger` and `AddContext(string, any) Logger`.

A default logger is registered in the ioc default scope, by calling the `LoggerFactory` with the name `root`. A default level for this level is also defined by configuring (see the default configuration source of [config](../config/README.md)):
 * the key `log.level` is defined with `Info`.
 * the key `log.level.root` is defined with `${log.level}`.


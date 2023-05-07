# Json

## Json

Json.

## Goals

This library has two purposes:
 * to be a Json encoding/decoding library, nice and efficient
 * integrate some Json marshallers and unmarshallers in the [ioc](../ioc) framework.

## Json core

The core of the package is dedicated to manipulate Json documents.

The Json document is abstracted behind the interface `JsonNode`. A `JsonNode` can be one of this types:
 * `JsonString` for a `string`.
 * `JsonInt` for an integer (`int`), which means a number without a decimal part.
 * `JsonFloat` for a float (`float64`), a number with a decimal part.
 * `JsonObject` representing a Json object. *The key order is saved and will be listed in the read or generation order.* The constant `JSON_EMPTY_OBJECT` represents an empty object.
 * `JsonArray` for a Json array. The constant `JSON_EMPTY_ARRAY` represents an empty array.
 * `JsonBool` for a boolean, with two constants `JSON_TRUE` and `JSON_FALSE` for `true` and `false`.
 * `JsonNull` for `null` with one defined instance `JSON_NULL`.

Each of these types is unmutable, and values can not be modified after been created. The Json serialisation can be obtained by calling the method `String`.

To create a JsonNode, you can:
 * Use the suitable constructor (`NewJsonObject` and `NewJsonArray`), cast the value to the expecting node type (for `JsonString`, `JsonInt` or `JsonFloat`) or use the defined constants (`JSON_TRUE`, `JSON_FALSE` or `JSON_NULL`).
 * Use a defined parser: `Parse`, `ParseAll`, `ParseString` or `ParseAllString`
 * Use a builder:
    ```go
    b := NewJsonBuilder()
    b.SetString("a", "hello") // Simple object field
    b.SetString("b.c\\.d", "world") // Nested object fields
    b.SetBool("e[0]", true) // Array
    b.SetInt("e[1].f", 42) // Mixing object and array
    node := b.Build()
    
    fmt.Print(node.String()) // expect `{"a":"hello","b":{"c.d":"world"},"e":[true,{"f":42}]}`
    ```
   The function `EscapePath` can be useful to escape special characters (`.` and `[`) in the path.

## Json lib

The library defines _marshallers_ and _unmarshallers_, which are ioc components used to convert a value to Json and vice versa. A marshaller should be defined as a `func(T) (JsonNode, error)` where `T` is the supported marshallable type. An unmarshaller should be defined as a `func(JsonNode) (T, error)` where `T` is again the supported unmarshallable type. Depending the usage, some type can be associate to only a marshaller or an unmarshaller, but several marshallers or unmarshallers can not be defined with the same type target.

The library will create automatically a marshaller for any type without an associated marshaller but implementing the `Jsoner` interface:
```go
type Jsoner interface {
	Json() JsonNode
}
```

All marshallers and unmarshallers are grouped and can be called from an interface `Jsons` with several methods:
 * The method `Marshal(any) (JsonNode, error)` use the correct marshaller to produce a `JsonNode`.
 * The method `MarshalToString(any) (string, error)` add a step to the `Marshal` method and convert the `JsonNode` to a `string`.
 * The method `Unmarshal(json JsonNode, callback func(T)) error` convert a node to a value and call the callback function with that value. The expected type is computed by looking at the callback argument.
 * The method `UnmarshalFromString(json string, callback func(T)) error` parse the given string before calling the method `Unmarshal`.

The json component will do his best to handle marshallers defined for interfaces and try to select the best matching interface to a given unkown instance. But this kind of marshallers will slow down the all process, and marshallers defined on concrete types (not interface) should be prefered. 

The lib defines some marshallers and unmarshallers. In order to be able to overwrite the default implementation but in the same time includes the implementation in the core scope, the default (un)marshallers are defined in 3 steps:
 * An interface is defined, with the same signature than a dedicated marshaller but without refencing it: e.g. for string marshalling:
    ```go
    type StringMarshaller func(v string) (JsonNode, error)
    ```
 * A default implementation of this interface is defined in the default scope. Again for string:
    ```go
	ioc.DefaultPutNamed("String Json marshaller (default)",
		func(v string) (JsonNode, error) {
			return JsonString(v), nil
		}, func(StringMarshaller) {})
    ```
 * Then, a factory is defined to promote the interface implementation as a valid (un)marshaller in the core scope:
    ```go
	ioc.PutNamedFactory("String Json marshaller (promoter)",
		func(m StringMarshaller) (JsonMarshaller, error) { return m, nil })
    ```
So, to overwrite an implementation, you have to register a component implementing the specific interface in the core or test scope. The default implementation and associated interfaces are:

| type | marshaller interface | unmarshaller interface |
| --- | --- | --- |
| `string` | `StringMarshaller` | `StringUnmarshaller` |
| `float64` | `Float64Marshaller` | `Float64Unmarshaller` |
| `int` | `IntMarshaller` | `IntUnmarshaller` |
| `bool` | `BoolMarshaller` | `BoolUnmarshaller` |
| `error` | `ErrorMarshaller` | `ErrorUnmarshaller` |

The marshaller and unmarshaller defined for `error` are very basic, only rely on the `Error() string` method and simply use a string node to represents the error (no supports for wrapped errors, nor joined errors).

The library also defines marshallers and unmarshallers for the different `JsonNode` implementations and also a marshaller for the `Jsons` implementation itself (for logging and debugging essentially). Theses component should not be overwritten.


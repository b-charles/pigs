# Json

## Json

[Json.](https://www.json.org)

## Goals

This library has two purposes:
 * to be a Json encoding/decoding library, nice and efficient,
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

Each of these types is unmutable: values can not be modified after been created. The Json serialisation can be obtained by calling the method `String() string`.

To create a `JsonNode`, you can:
 * Use the different core functions:
    * For objects, you can use the constructor `NewJsonObjectMappedSorted[T any](members map[string]T, mapper func(v T) JsonNode, less func(a, b string) bool) *JsonObject` which is the general method to convert a map to a Json object, mapping each value of the map to a `JsonNode` with the `mapper` argument and sorting the Json object members with the `less` methods. Other more convenient constructors are also defined, like `NewJsonObjectMapped[T any](members map[string]T, mapper func(v T) JsonNode) *JsonObject` which use a default string comparator, `NewJsonObjectSorted(members map[string]JsonNode, less func(a, b string) bool)` for maps with already Json nodes values, and `NewJsonObject(members map[string]JsonNode) *JsonObject`. Two additional constructors `NewJsonObjectStringsSorted(members map[string]string, less func(a, b string) bool) *JsonObject` and `NewJsonObjectStrings(members map[string]string) *JsonObject` are defined to quickly convert a `map[string]string` to a Json object. Finally, the constant `JSON_EMPTY_OBJECT` can also be used for immutable empty object.
    * For arrays, you can use the general constructor `NewJsonArrayMapped[T any](elements []T, mapper func(v T) JsonNode) *JsonArray` which convert each element of a slice to a `JsonNode` by using the `mapper` function. More convenient constructor are defined for classical slices: `NewJsonArray(elements []JsonNode) *JsonArray`, `NewJsonArrayStrings(elements []string) *JsonArray`, `NewJsonArrayFloats(elements []float64) *JsonArray`, `NewJsonArrayInts(elements []int) *JsonArray` and `NewJsonArrayBools(elements []bool) *JsonArray`. Finally, the constant `JSON_EMPTY_ARRAY` can also be used for immutable empty array.
    * You can directly cast a `string` to a `JsonString`, an `int` to a `JsonInt` or a `float64` to a `JsonFloat`.
    * Finally, for a `JsonBool`, you should use the constants `JSON_TRUE` and `JSON_FALSE`, and the only instance available for a `JsonNull` is the constant `JSON_NULL`.
 * Use a defined parser: `Parse(source io.RuneReader) (JsonNode, error)`, `ParseAll(source io.RuneReader) ([]JsonNode, error)`, `ParseString(source string) (JsonNode, error)` or `ParseAllString(source string) ([]JsonNode, error)`
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
   The function `EscapePath(path string) string` can be useful to escape special characters (`.` and `[`) in the path.

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

The json component will do his best to handle marshallers defined for interfaces and try to select the best matching (most specific) interface to a given unkown instance. But this kind of marshallers will slow down the all process, and marshallers defined on concrete types (not interface) should be prefered. 

The lib defines some marshallers and unmarshallers. In order to be able to overwrite the default implementation but in the same time includes the implementation in the core scope, the default (un)marshallers are defined following the [3 steps tricks to register overloadable components](../ioc/README.md#overloadable-components-in-auto-discovery-injection). To overwrite an implementation, you have to register a component implementing with the specific signature in the core or test scope:

| type | marshaller signature | unmarshaller signature |
| --- | --- | --- |
| `string` | `type StringMarshaller func(string) (JsonNode, error)` | `type StringUnmarshaller func(JsonNode) (string, error)` |
| `float64` | `type Float64Marshaller func(float64) (JsonNode, error)` | `type Float64Unmarshaller func(JsonNode) (float64, error)` |
| `int` | `type IntMarshaller func(int) (JsonNode, error)` | `type IntUnmarshaller func(JsonNode) (int, error)` |
| `bool` | `type BoolMarshaller func(bool) (JsonNode, error)` | `type BoolUnmarshaller func(JsonNode) (bool, error)` |
| `error` | `type ErrorMarshaller func(v error) (JsonNode, error)` | `type ErrorUnmarshaller func(JsonNode) (error, error)` |

The marshaller and unmarshaller defined for `error` are very basic, only rely on the `Error() string` method and simply use a string node to represents the error (no supports for wrapped errors, nor joined errors).

The library also defines marshallers and unmarshallers for the different `JsonNode` implementations and also a marshaller for the `Jsons` implementation itself (for logging and debugging essentially). Theses component should not be overwritten.


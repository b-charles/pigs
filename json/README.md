# Json

## Json

Json.

## Goals

This library has two purposes:
 * to be a Json encoding/decoding library, nice and efficient
 * integrate some Json marshallers and unmarshallers in the [ioc](../ioc) framework.

## Json core

Json core is a package to manipulate Json documents.

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
 * Use the dedicated constructor (`NewJsonObject` and `NewJsonArray`), cast the value to the expecting node type (for `JsonString`, `JsonInt` or `JsonFloat`) or use the defined constants (`JSON_TRUE`, `JSON_FALSE` or `JSON_NULL`).
 * Use a defined parser: `Parse`, `ParseAll`, `ParseString` or `ParseAllString`
 * Use a builder:
```go
b := NewJsonBuilder()
b.SetString("a", "hello") // Simple object field
b.SetString("b.c.d", "world") // Nested object fields
b.SetBool("e[0]", true) // Array
b.SetInt("e[1].f", 42) // Mixing object and array
node := b.Build()

fmt.Print(node.String()) // expect `{"a":"hello","b":{"c":{"d":"world"}},"e":[true,{"f":42}]}`
```

## Json lib

The main library use the core to define _marshallers_ and _unmarshallers_, which are ioc components used to convert a value to Json and vice versa. A marshaller should be defined as a `func(T) (core.JsonNode, error)` where `T` is the supported marshallable type. An unmarshaller should be defined as a `func(core.JsonNode) (T, error)` where `T` is again the supported unmarshallable type. Depending the usage, some type can be associate to only a marshaller or an unmarshaller, but several marshallers or unmarshallers can not be defined with the same type target.

All marshallers and unmarshallers are grouped and can be called from an interface `Json` with several methods:
 * The method `Marshal(any) (core.JsonNode, error)` use the correct marshaller to produce a `JsonNode`.
 * The method `MarshalToString(any) (string, error)` add a step to the `Marshal` method and convert the `JsonNode` to a `string`.
 * The method `Unmarshal(json core.JsonNode, callback func(T)) error` convert a node to a value and call the callback function with that value. The expected type is computed by looking at the callback argument.
 * The method `UnmarshalFromString(json string, callback any) error` parse the given string before calling the method `Unmarshal`.


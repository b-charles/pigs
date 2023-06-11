# MemFun

Fun with Mem(oization).

## What's that?

MemFun is a small library to handle synchronized function memoization. It wraps a potentially recursive function (one input, one output and an error) and add a multi-thread proof memoization mechanism. For each input value, the given function is only called once. The function can be executed in parallel, in different goroutines, if the keys are different.

MemFun **is not**:
* [A synchronized map](https://pkg.go.dev/sync#Map): The API of this Map is not satisfactory, since it only accepts precomputed values. MemFun encapsulates a function to memoize, which is called when it is necessary.
* [A cache](https://github.com/patrickmn/go-cache): MemFun doesn't define expiration rules in memory or in time.
* [A single flight](https://github.com/golang/sync/blob/v0.2.0/singleflight/singleflight.go): Unlike this lib, MemFun stores the results for latter calls. But it has been a great source of inspiration.

## The API

The library starts by defining a type `type Computer[K any, V any] func(key K, recfun func(key K) (V, error)) (V, error)`. This type represents the function to memoize. It's parametrized with the types, `K` for the input/key type, and `V` for the output/value type. The parameter `recfun` can be used during the execution of the function for recursive values, or can be stored in the result for a lazy loading.

To memoize this `Computer`, wrap it in an `MemFun` instance with the function `memfun.NewMemFun[K any, V any](fun Computer[K, V]) MemFun[K, V]`. The `MemFun` interface is defined with:
```go
type MemFun[K any, V any] interface {
	Store(key K, value V)
	Get(key K) (V, error)
	Lookup(key K) (V, bool, error)
	Keys() []K
	Delete(key K)
}
```
 * The `Store` method can be used to store directly a value, bypassing the `Computer` function.
 * The `Get` method can be used to retreive a value, stored by the `Store` method or computed by the `Computer` function.
 * The `Lookup` method returns the same value as `Get`, except the value is not computed if it is missing. The additional boolean output can be used to determine if the value is present (`true`) or not (`false`).
 * The `Keys` method returns the currently present (memoized) keys.
 * The `Delete` method can be use to delete an entry.

If the function returns an error for some key, every call with the same key will returns the same error. If the method panics, every call with the same key will panic.

The `Get` and `Lookup` methods detect cyclic calls. In that case, a `CyclicLoopError` is returned as an error and the stack of the called key can be retrieve in the `Stack` property of the error.


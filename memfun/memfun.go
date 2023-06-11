package memfun

import "sync"

type Computer[K any, V any] func(key K, recfun func(key K) (V, error)) (V, error)

type MemFun[K any, V any] interface {
	Store(key K, value V)
	Get(key K) (V, error)
	Lookup(key K) (V, bool, error)
	Keys() []K
	Delete(key K)
}

type memVal struct {
	wg    sync.WaitGroup
	val   any
	err   error
	panic any
}

func newMemVal(call bool) *memVal {

	v := &memVal{}
	if call {
		v.wg.Add(1)
	}

	return v

}

func (self *memVal) compute(fun func() (any, error)) (any, error) {

	func() {

		defer func() {

			if r := recover(); r != nil {
				self.panic = r
			}

			self.wg.Done()

		}()

		self.val, self.err = fun()

	}()

	if self.panic != nil {
		panic(self.panic)
	}

	return self.val, self.err

}

func (self *memVal) get() (any, error) {

	self.wg.Wait()

	if self.panic != nil {
		panic(self.panic)
	}

	return self.val, self.err

}

type memFunImpl[K any, V any] struct {
	fun Computer[K, V]
	mu  sync.RWMutex
	m   map[any]*memVal
}

func (self *memFunImpl[K, V]) Store(key K, value V) {

	v := newMemVal(false)
	v.val = value

	self.mu.Lock()
	self.m[key] = v
	self.mu.Unlock()

}

func (self *memFunImpl[K, V]) Get(key K) (V, error) {
	return self.get(key, make(map[any]bool))
}

func (self *memFunImpl[K, V]) get(key K, called map[any]bool) (V, error) {

	self.mu.RLock()
	if v, ok := self.m[key]; ok {
		self.mu.RUnlock()
		val, err := v.get()
		return val.(V), err
	}
	self.mu.RUnlock()

	self.mu.Lock()
	if v, ok := self.m[key]; ok {
		self.mu.Unlock()
		val, err := v.get()
		return val.(V), err
	}

	v := newMemVal(true)
	self.m[key] = v
	self.mu.Unlock()

	val, err := v.compute(func() (any, error) {

		called[key] = true
		defer delete(called, key)

		r, err := self.fun(key, func(k K) (V, error) {

			if _, p := called[k]; p {
				var nilV V
				return nilV, CyclicLoopError[K]{[]K{k}}
			}

			return self.get(k, called)

		})

		if err != nil {
			if cyclic, ok := err.(CyclicLoopError[K]); ok {
				return r, cyclic.Append(key)
			}
		}

		return r, err

	})

	return val.(V), err

}

func (self *memFunImpl[K, V]) Lookup(key K) (V, bool, error) {

	self.mu.RLock()
	defer self.mu.RUnlock()

	if m, ok := self.m[key]; ok {
		v, e := m.get()
		return v.(V), true, e
	} else {
		var v V
		return v, false, nil
	}

}

func (self *memFunImpl[K, V]) Keys() []K {

	self.mu.RLock()
	defer self.mu.RUnlock()

	keys := make([]K, 0, len(self.m))
	for k := range self.m {
		keys = append(keys, k.(K))
	}

	return keys

}

func (self *memFunImpl[K, V]) Delete(key K) {
	self.mu.Lock()
	delete(self.m, key)
	self.mu.Unlock()
}

func NewMemFun[K any, V any](fun Computer[K, V]) MemFun[K, V] {

	return &memFunImpl[K, V]{
		fun: fun,
		m:   make(map[any]*memVal),
	}

}

package core

import "fmt"

type sortedMap[T any] struct {
	m map[string]T
	s []string
}

func (self *sortedMap[T]) put(key string, elt T) {
	if _, p := self.m[key]; !p {
		self.s = append(self.s, key)
	}
	self.m[key] = elt
}

func (self *sortedMap[T]) len() int {
	return len(self.s)
}

func (self *sortedMap[T]) keys() []string {
	return self.s
}

func (self *sortedMap[T]) get(key string) (T, bool) {
	e, p := self.m[key]
	return e, p
}

func newEmptySortedMap[T any]() *sortedMap[T] {
	return &sortedMap[T]{make(map[string]T), make([]string, 0)}
}

func newSortedMap[T any](m map[string]T, s []string) *sortedMap[T] {

	u := make(map[string]bool, len(s))
	for _, k := range s {
		if _, ok := m[k]; !ok {
			panic(fmt.Sprintf("The key '%v' is in the sorted slice, but not in the given map.", k))
		}
		if _, ok := u[k]; ok {
			panic(fmt.Sprintf("The key '%v' is several times in the sorted slice.", k))
		}
		u[k] = true
	}

	return &sortedMap[T]{m, s}

}

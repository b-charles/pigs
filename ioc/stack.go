package ioc

import (
	"fmt"
	"strings"
	"sync"
)

type componentStack struct {
	stack   []*component
	present map[*component]bool
}

func newComponentStack() *componentStack {
	return &componentStack{
		[]*component{},
		map[*component]bool{},
	}
}

func (self *componentStack) push(component *component) error {

	if self.present[component] {

		i := 0
		for self.stack[i] != component {
			i++
		}

		return &cyclicError{
			components: self.stack[i:],
		}

	}

	self.stack = append(self.stack, component)
	self.present[component] = true

	return nil

}

func (self *componentStack) pop(component *component) {

	if last := self.stack[len(self.stack)-1]; last != component {
		panic(fmt.Sprintf("Error of using the components stack: last element: %v, expecting: %v", last, component))
	}

	self.stack = self.stack[:len(self.stack)-1]
	delete(self.present, component)

}

type cyclicError struct {
	components []*component
	message    string
	once       sync.Once
}

func (self *cyclicError) Error() string {

	self.once.Do(func() {

		var b strings.Builder
		fmt.Fprintf(&b, "Cyclic dependency detected: ")

		first := fmt.Sprintf("%v", self.components[0].String())
		fmt.Fprintf(&b, "%v -> ", first)

		for i := 1; i < len(self.components); i++ {
			fmt.Fprintf(&b, "%v -> ", self.components[i])
		}
		fmt.Fprintf(&b, "%v", first)

		self.message = b.String()

	})

	return self.message

}

func (self *cyclicError) String() string {
	return self.Error()
}

func (self *cyclicError) Components() []*component {
	return self.components
}

func isDirectCyclicError(err error) bool {
	cyclic, ok := err.(*cyclicError)
	return ok && len(cyclic.components) == 1
}

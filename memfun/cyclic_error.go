package memfun

import (
	"fmt"
	"strings"
)

type CyclicLoopError[K any] struct {
	Stack []K
}

func (self CyclicLoopError[K]) Error() string {

	var b strings.Builder

	b.WriteString("Cyclic loop detected: ")

	b.WriteString(fmt.Sprintf("%v", self.Stack[0]))
	for i := 1; i < len(self.Stack); i++ {
		b.WriteString(" -> ")
		b.WriteString(fmt.Sprintf("%v", self.Stack[i]))
	}

	return b.String()

}

func (self CyclicLoopError[K]) Append(elt K) CyclicLoopError[K] {
	return CyclicLoopError[K]{append([]K{elt}, self.Stack...)}
}

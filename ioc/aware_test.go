package ioc_test

import (
	"reflect"

	. "github.com/b-charles/pigs/ioc"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type AwaredCallback struct {
	reforder *int
	order    int
}

func (self *AwaredCallback) incr() error {
	(*self.reforder)++
	self.order = *self.reforder
	return nil
}

func (self *AwaredCallback) Precall(method reflect.Value) error {
	return self.incr()
}

func (self *AwaredCallback) Postinst(method reflect.Value, args []reflect.Value) error {
	return self.incr()
}

func (self *AwaredCallback) Postcall(method reflect.Value, outs []reflect.Value) error {
	return self.incr()
}

func (self *AwaredCallback) Preclose() {
	self.incr()
}

func (self *AwaredCallback) Postclose() {
	self.incr()
}

var _ = Describe("Container awareness", func() {

	var (
		container *Container
		reforder  *int
	)

	BeforeEach(func() {
		container = NewContainer()
		i := 0
		reforder = &i
	})

	It("should call the CallInjected callbacks in the correct order", func() {

		precall := &AwaredCallback{reforder, 0}
		container.PutNamed(precall, "PreCall", func(PreCallAwared) {})

		postinst := &AwaredCallback{reforder, 0}
		container.PutNamed(postinst, "PostInst", func(PostInstAwared) {})

		preclose := &AwaredCallback{reforder, 0}
		container.PutNamed(preclose, "PreClose", func(PreCloseAwared) {})

		postclose := &AwaredCallback{reforder, 0}
		container.PutNamed(postclose, "PostClose", func(PostCloseAwared) {})

		container.CallInjected(func() {})

		Expect(precall.order).To(Equal(1))
		Expect(postinst.order).To(Equal(2))
		Expect(preclose.order).To(Equal(3))
		Expect(postclose.order).To(Equal(4))

	})

})

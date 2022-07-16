package ioc_test

import (
	"reflect"

	. "github.com/b-charles/pigs/ioc"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var i int = 0
var reforder *int = &i

type AwaredCallback struct {
	order int
}

func (self *AwaredCallback) Incr() error {
	(*reforder)++
	self.order = *reforder
	return nil
}

func (self *AwaredCallback) Precall(method reflect.Value) error {
	return self.Incr()
}

func (self *AwaredCallback) Postinst(method reflect.Value, args []reflect.Value) error {
	return self.Incr()
}

func (self *AwaredCallback) Postcall(method reflect.Value, outs []reflect.Value) error {
	return self.Incr()
}

func (self *AwaredCallback) Preclose() {
	self.Incr()
}

func (self *AwaredCallback) Postclose() {
	self.Incr()
}

var _ = Describe("Container awareness", func() {

	var (
		container *Container
		err       error
	)

	BeforeEach(func() {
		container = NewContainer()
		(*reforder) = 0
	})

	It("should call the CallInjected callbacks in the correct order", func() {

		precall := &AwaredCallback{}
		err = container.PutNamed(precall, "PreCall", func(PreCallAwared) {})
		Expect(err).To(Succeed())

		postinst := &AwaredCallback{}
		err = container.PutNamed(postinst, "PostInst", func(PostInstAwared) {})
		Expect(err).To(Succeed())

		preclose := &AwaredCallback{}
		err = container.PutNamed(preclose, "PreClose", func(PreCloseAwared) {})
		Expect(err).To(Succeed())

		postclose := &AwaredCallback{}
		err = container.PutNamed(postclose, "PostClose", func(PostCloseAwared) {})
		Expect(err).To(Succeed())

		call := &AwaredCallback{}
		err = container.CallInjected(func() { call.Incr() })
		Expect(err).To(Succeed())

		Expect(precall.order).To(Equal(1))
		Expect(postinst.order).To(Equal(2))
		Expect(call.order).To(Equal(3))
		Expect(preclose.order).To(Equal(4))
		Expect(postclose.order).To(Equal(5))

	})

})

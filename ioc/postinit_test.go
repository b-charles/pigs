package ioc_test

import (
	. "github.com/b-charles/pigs/ioc"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("IOC post-init and close", func() {

	var container *Container

	BeforeEach(func() {
		container = NewContainer()
	})

	It("should call PostInit in the correct order", func() {

		container.PutFactory(NewOrderRegister)
		container.Put(&First{})
		container.Put(&Second{})
		container.Put(&Third{})

		var third *Third
		Expect(container.CallInjected(func(injected *Third) {
			third = injected
		})).To(Succeed())

		Expect(third.Register.PostInitOrder).To(Equal([]string{"FIRST", "SECOND", "THIRD"}))

	})

	It("should call Close in the correct order", func() {

		container.PutFactory(NewOrderRegister)
		container.Put(&First{})
		container.Put(&Second{})
		container.Put(&Third{})

		var third *Third
		Expect(container.CallInjected(func(injected *Third) {
			third = injected
		})).To(Succeed())

		Expect(third.Register.CloseOrder).To(Equal([]string{"FIRST", "SECOND", "THIRD"}))

	})

})

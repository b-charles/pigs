package ioc_test

import (
	. "github.com/b-charles/pigs/ioc"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("IOC post-init and close", func() {

	var container *Container

	BeforeEach(func() {
		container = NewContainer()
	})

	It("should call PostInit in the correct order", func() {

		container.RegisterFactory(Core, "REGISTER", NewOrderRegister)
		container.RegisterComponent(Core, "1", &First{})
		container.RegisterComponent(Core, "2", &Second{})
		container.RegisterComponent(Core, "3", &Third{})

		var third *Third
		Expect(container.CallInjected(func(injected *Third) {
			third = injected
		})).To(Succeed())

		Expect(third.Register.PostInitOrder).To(Equal([]string{"FIRST", "SECOND", "THIRD"}))

	})

	It("should call Close in the correct order", func() {

		container.RegisterFactory(Core, "REGISTER", NewOrderRegister)
		container.RegisterComponent(Core, "1", &First{})
		container.RegisterComponent(Core, "2", &Second{})
		container.RegisterComponent(Core, "3", &Third{})

		var third *Third
		Expect(container.CallInjected(func(injected *Third) {
			third = injected
		})).To(Succeed())

		Expect(third.Register.CloseOrder).To(Equal([]string{"THIRD", "SECOND", "FIRST"}))

	})

})

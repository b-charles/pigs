package ioc_test

import (
	. "github.com/b-charles/pigs/ioc"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("IOC post-init and close", func() {

	var (
		container *Container
	)

	BeforeEach(func() {
		container = NewContainer()
	})

	It("should call PostInit in the correct order", func() {

		container.PutFactory(NewOrderRegster, "OrderRegister")
		container.Put(&First{}, "First")
		container.Put(&Second{}, "Second")
		container.Put(&Third{}, "Third")

		var injectedThird *Third

		Expect(container.CallInjected(func(injected struct {
			Third *Third
		}) {
			injectedThird = injected.Third
		})).Should(Succeed())

		Expect(injectedThird.OrderRegister.PostInitOrder).Should(Equal([]string{"FIRST", "SECOND", "THIRD"}))

	})

	It("should call Close in the correct order", func() {

		container.PutFactory(NewOrderRegster, "OrderRegister")
		container.Put(&First{}, "First")
		container.Put(&Second{}, "Second")
		container.Put(&Third{}, "Third")

		var injectedThird *Third

		Expect(container.CallInjected(func(injected struct {
			Third *Third
		}) {
			injectedThird = injected.Third
		})).Should(Succeed())

		container.Close()

		Expect(injectedThird.OrderRegister.CloseOrder).Should(Equal([]string{"FIRST", "SECOND", "THIRD"}))

	})

})

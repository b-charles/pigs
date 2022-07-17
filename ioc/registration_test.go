package ioc_test

import (
	. "github.com/b-charles/pigs/ioc"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("IOC registration", func() {

	var container *Container

	BeforeEach(func() {
		container = NewContainer()
	})

	Describe("with factory", func() {

		It("should register", func() {
			Expect(container.PutFactory(SimpleFactory("A"))).To(Succeed())
		})

		It("should register with signatures", func() {
			Expect(container.PutFactory(SimpleFactory("A"), func(Doer, BigDoer) {})).To(Succeed())
		})

		It("should not register with incorrect signature", func() {
			Expect(container.PutFactory(SimpleFactory("A"), func(NotDoer) {})).To(HaveOccurred())
		})

	})

	Describe("with instanciated component", func() {

		It("should register", func() {
			Expect(container.Put(&Simple{"A"})).To(Succeed())
		})

		It("should register with signatures", func() {
			Expect(container.Put(&Simple{"A"}, func(Doer, BigDoer) {})).To(Succeed())
		})

		It("should not register with incorrect signature", func() {
			Expect(container.Put(&Simple{"A"}, func(NotDoer) {})).To(HaveOccurred())
		})

	})

	Describe("with test factory", func() {

		It("should register", func() {
			Expect(container.TestPutFactory(SimpleFactory("A"))).To(Succeed())
		})

		It("should register with signatures", func() {
			Expect(container.TestPutFactory(SimpleFactory("A"), func(Doer, BigDoer) {})).To(Succeed())
		})

		It("should not register with incorrect signature", func() {
			Expect(container.TestPutFactory(SimpleFactory("A"), func(NotDoer) {})).To(HaveOccurred())
		})

	})

	Describe("with instanciated test component", func() {

		It("should register", func() {
			Expect(container.TestPut(&Simple{"A"})).To(Succeed())
		})

		It("should register with signatures", func() {
			Expect(container.TestPut(&Simple{"A"}, func(Doer, BigDoer) {})).To(Succeed())
		})

		It("should not register with incorrect signature", func() {
			Expect(container.TestPut(&Simple{"A"}, func(NotDoer) {})).To(HaveOccurred())
		})

	})

})

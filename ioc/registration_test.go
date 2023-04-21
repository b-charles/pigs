package ioc_test

import (
	. "github.com/b-charles/pigs/ioc"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("IOC registration", func() {

	var container *Container

	BeforeEach(func() {
		container = NewContainer()
	})

	Describe("with default factory", func() {

		It("should register", func() {
			Expect(container.RegisterFactory(Def, "A", SimpleFactory("A"))).To(Succeed())
		})

		It("should register with signatures", func() {
			Expect(container.RegisterFactory(Def, "A", SimpleFactory("A"), func(Doer, BigDoer) {})).To(Succeed())
		})

		It("should not register with incorrect signature", func() {
			Expect(container.RegisterFactory(Def, "A", SimpleFactory("A"), func(NotDoer) {})).To(HaveOccurred())
		})

	})

	Describe("with instanciated default component", func() {

		It("should register", func() {
			Expect(container.RegisterComponent(Def, "A", &Simple{"A"})).To(Succeed())
		})

		It("should register with signatures", func() {
			Expect(container.RegisterComponent(Def, "A", &Simple{"A"}, func(Doer, BigDoer) {})).To(Succeed())
		})

		It("should not register with incorrect signature", func() {
			Expect(container.RegisterComponent(Def, "A", &Simple{"A"}, func(NotDoer) {})).To(HaveOccurred())
		})

	})

	Describe("with factory", func() {

		It("should register", func() {
			Expect(container.RegisterFactory(Core, "A", SimpleFactory("A"))).To(Succeed())
		})

		It("should register with signatures", func() {
			Expect(container.RegisterFactory(Core, "A", SimpleFactory("A"), func(Doer, BigDoer) {})).To(Succeed())
		})

		It("should not register with incorrect signature", func() {
			Expect(container.RegisterFactory(Core, "A", SimpleFactory("A"), func(NotDoer) {})).To(HaveOccurred())
		})

	})

	Describe("with instanciated component", func() {

		It("should register", func() {
			Expect(container.RegisterComponent(Core, "A", &Simple{"A"})).To(Succeed())
		})

		It("should register with signatures", func() {
			Expect(container.RegisterComponent(Core, "A", &Simple{"A"}, func(Doer, BigDoer) {})).To(Succeed())
		})

		It("should not register with incorrect signature", func() {
			Expect(container.RegisterComponent(Core, "A", &Simple{"A"}, func(NotDoer) {})).To(HaveOccurred())
		})

	})

	Describe("with test factory", func() {

		It("should register", func() {
			Expect(container.RegisterFactory(Test, "A", SimpleFactory("A"))).To(Succeed())
		})

		It("should register with signatures", func() {
			Expect(container.RegisterFactory(Test, "A", SimpleFactory("A"), func(Doer, BigDoer) {})).To(Succeed())
		})

		It("should not register with incorrect signature", func() {
			Expect(container.RegisterFactory(Test, "A", SimpleFactory("A"), func(NotDoer) {})).To(HaveOccurred())
		})

	})

	Describe("with instanciated test component", func() {

		It("should register", func() {
			Expect(container.RegisterComponent(Test, "A", &Simple{"A"})).To(Succeed())
		})

		It("should register with signatures", func() {
			Expect(container.RegisterComponent(Test, "A", &Simple{"A"}, func(Doer, BigDoer) {})).To(Succeed())
		})

		It("should not register with incorrect signature", func() {
			Expect(container.RegisterComponent(Test, "A", &Simple{"A"}, func(NotDoer) {})).To(HaveOccurred())
		})

	})

})

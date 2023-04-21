package ioc_test

import (
	. "github.com/b-charles/pigs/ioc"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("IOC factory", func() {

	var container *Container

	BeforeEach(func() {
		container = NewContainer()
	})

	Describe("with simple injection", func() {

		It("should create a component with dependency (factory)", func() {

			container.RegisterFactory(Core, "A", SimpleFactory("A"))

			container.RegisterFactory(Core, "INJECTED", InjectedFactory)

			container.CallInjected(func(injected *Injected) {
				Expect(injected).To(Equal(&Injected{&Simple{"A"}}))
			})

		})

		It("should create a component with dependency (init)", func() {

			container.RegisterFactory(Core, "A", SimpleFactory("A"))

			container.RegisterComponent(Core, "INJECTED", &Initialized{})

			container.CallInjected(func(injected *Initialized) {
				Expect(injected).To(Equal(&Initialized{&Simple{"A"}}))
			})

		})

		It("should create a component with dependency (postInit)", func() {

			container.RegisterFactory(Core, "A", SimpleFactory("A"))

			container.RegisterComponent(Core, "INJECTED", &PostInitialized{})

			container.CallInjected(func(injected *PostInitialized) {
				Expect(injected).To(Equal(&PostInitialized{&Simple{"A"}}))
			})

		})

		It("should inject test component if provided", func() {

			container.RegisterFactory(Core, "A", SimpleFactory("A"))
			container.RegisterFactory(Test, "TEST", SimpleFactory("TEST"))

			container.RegisterFactory(Core, "INJECTED", InjectedFactory)

			container.CallInjected(func(injected *Injected) {
				Expect(injected).To(Equal(&Injected{&Simple{"TEST"}}))
			})

		})

		It("should return an error if no component is provided", func() {

			container.RegisterFactory(Core, "INJECTED", InjectedFactory)

			Expect(container.CallInjected(func(injected *Injected) {})).To(HaveOccurred())

		})

		It("should return an error if too many components are provided", func() {

			container.RegisterFactory(Core, "A1", SimpleFactory("A1"))
			container.RegisterFactory(Core, "A2", SimpleFactory("A2"))

			container.RegisterFactory(Core, "INJECTED", InjectedFactory)

			Expect(container.CallInjected(func(injected *Injected) {})).To(HaveOccurred())

		})

		It("should inject the component in itself", func() {

			container.RegisterComponent(Core, "LOOPING", &Looping{})

			Expect(container.CallInjected(func(injected *Looping) {
				Expect(injected).To(Equal(injected.Looping))
			})).To(Succeed())

		})

		It("should throw an error if a cyclic dependency is detected", func() {

			container.RegisterFactory(Core, "LOOPING", func(looping *Looping) *Looping {
				return &Looping{looping}
			})

			Expect(container.CallInjected(func(injected *Looping) {
				Expect(injected).To(Equal(injected.Looping))
			})).To(HaveOccurred())

		})

		It("should be possible to promote a core component to test", func() {

			container.RegisterFactory(Core, "A1", SimpleFactory("A1"))
			container.RegisterFactory(Test, "PROM", func(a1 *Simple) *Simple { return a1 })

			container.RegisterComponent(Core, "INJECTED", &SliceInjected{})

			Expect(container.CallInjected(func(injected *SliceInjected) {
				Expect(injected.Simple).To(ConsistOf(&Simple{"A1"}))
			})).To(Succeed())

		})

	})

	Describe("with interface", func() {

		It("should create a component with dependency (factory)", func() {

			container.RegisterFactory(Core, "A", SimpleFactory("A"))

			container.RegisterFactory(Core, "INJECTED", InterfaceInjectedFactory)

			container.CallInjected(func(injected *InterfaceInjected) {
				Expect(injected).To(Equal(&InterfaceInjected{&Simple{"A"}}))
			})

		})

		It("should create a component with dependency (init)", func() {

			container.RegisterFactory(Core, "A", SimpleFactory("A"))

			container.RegisterComponent(Core, "INJECTED", &InterfaceInitialized{})

			container.CallInjected(func(injected *InterfaceInitialized) {
				Expect(injected).To(Equal(&InterfaceInitialized{&Simple{"A"}}))
			})

		})

		It("should create a component with dependency (postInit)", func() {

			container.RegisterFactory(Core, "A", SimpleFactory("A"))

			container.RegisterComponent(Core, "INJECTED", &InterfacePostInitialized{})

			container.CallInjected(func(injected *InterfacePostInitialized) {
				Expect(injected).To(Equal(&InterfacePostInitialized{&Simple{"A"}}))
			})

		})

	})

	Describe("with slice injection", func() {

		It("should create a component with a slice of dependencies", func() {

			container.RegisterFactory(Core, "A1", SimpleFactory("A1"), func(Doer) {})
			container.RegisterFactory(Core, "A2", SimpleFactory("A2"), func(Doer) {})
			container.RegisterFactory(Core, "A3", SimpleFactory("A3"), func(Doer) {})

			container.RegisterComponent(Core, "INJECTED", &SliceInjected{})

			container.CallInjected(func(injected *SliceInjected) {
				Expect(injected.Simple).To(ConsistOf(
					&Simple{"A1"},
					&Simple{"A2"},
					&Simple{"A3"}))
			})

		})

		It("should inject a slice discarding nil components", func() {

			container.RegisterFactory(Core, "A1", SimpleFactory("A1"))
			container.RegisterFactory(Core, "NIL", func() *Simple { return nil })
			container.RegisterFactory(Core, "A3", SimpleFactory("A3"))

			container.RegisterComponent(Core, "INJECTED", &SliceInjected{})

			container.CallInjected(func(injected *SliceInjected) {
				Expect(injected.Simple).To(ConsistOf(
					&Simple{"A1"},
					&Simple{"A3"}))
			})

		})

		It("should inject a slice of interface", func() {

			container.RegisterFactory(Core, "A1", SimpleFactory("A1"), func(Doer) {})
			container.RegisterFactory(Core, "A2", SimpleFactory("A2"), func(Doer) {})
			container.RegisterFactory(Core, "T1", TrivialFactory("T1"), func(Doer) {})

			container.RegisterComponent(Core, "INJECTED", &InterfaceSliceInjected{})

			container.CallInjected(func(injected *InterfaceSliceInjected) {
				Expect(injected.Doers).To(ConsistOf(
					&Simple{"A1"},
					&Simple{"A2"},
					Trivial("T1")))
			})

		})

		It("should inject an empty slice", func() {

			container.RegisterComponent(Core, "INJECTED", &InterfaceSliceInjected{})

			container.CallInjected(func(injected *InterfaceSliceInjected) {
				Expect(injected.Doers).To(BeEmpty())
			})

		})

	})

})

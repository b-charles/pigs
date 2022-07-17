package ioc_test

import (
	. "github.com/b-charles/pigs/ioc"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("IOC factory", func() {

	var container *Container

	BeforeEach(func() {
		container = NewContainer()
	})

	Describe("with simple injection", func() {

		It("should create a component with dependency (factory)", func() {

			container.PutFactory(SimpleFactory("A"))

			container.PutFactory(InjectedFactory)

			container.CallInjected(func(injected *Injected) {
				Expect(injected).To(Equal(&Injected{&Simple{"A"}}))
			})

		})

		It("should create a component with dependency (init)", func() {

			container.PutFactory(SimpleFactory("A"))

			container.Put(&Initialized{})

			container.CallInjected(func(injected *Initialized) {
				Expect(injected).To(Equal(&Initialized{&Simple{"A"}}))
			})

		})

		It("should create a component with dependency (postInit)", func() {

			container.PutFactory(SimpleFactory("A"))

			container.Put(&PostInitialized{})

			container.CallInjected(func(injected *PostInitialized) {
				Expect(injected).To(Equal(&PostInitialized{&Simple{"A"}}))
			})

		})

		It("should inject test component if provided", func() {

			container.PutFactory(SimpleFactory("A"))
			container.TestPutFactory(SimpleFactory("TEST"))

			container.PutFactory(InjectedFactory)

			container.CallInjected(func(injected *Injected) {
				Expect(injected).To(Equal(&Injected{&Simple{"TEST"}}))
			})

		})

		It("should return an error if no component is provided", func() {

			container.PutFactory(InjectedFactory)

			Expect(container.CallInjected(func(injected *Injected) {})).To(HaveOccurred())

		})

		It("should return an error if too many components are provided", func() {

			container.PutFactory(SimpleFactory("A1"))
			container.PutFactory(SimpleFactory("A2"))

			container.PutFactory(InjectedFactory)

			Expect(container.CallInjected(func(injected *Injected) {})).To(HaveOccurred())

		})

		It("should inject the component in itself", func() {

			container.Put(&Looping{})

			Expect(container.CallInjected(func(injected *Looping) {
				Expect(injected).To(Equal(injected.Looping))
			})).To(Succeed())

		})

	})

	Describe("with interface", func() {

		It("should create a component with dependency (factory)", func() {

			container.PutFactory(SimpleFactory("A"))

			container.PutFactory(InterfaceInjectedFactory)

			container.CallInjected(func(injected *InterfaceInjected) {
				Expect(injected).To(Equal(&InterfaceInjected{&Simple{"A"}}))
			})

		})

		It("should create a component with dependency (init)", func() {

			container.PutFactory(SimpleFactory("A"))

			container.Put(&InterfaceInitialized{})

			container.CallInjected(func(injected *InterfaceInitialized) {
				Expect(injected).To(Equal(&InterfaceInitialized{&Simple{"A"}}))
			})

		})

		It("should create a component with dependency (postInit)", func() {

			container.PutFactory(SimpleFactory("A"))

			container.Put(&InterfacePostInitialized{})

			container.CallInjected(func(injected *InterfacePostInitialized) {
				Expect(injected).To(Equal(&InterfacePostInitialized{&Simple{"A"}}))
			})

		})

	})

	Describe("with slice injection", func() {

		It("should create a component with a slice of dependencies", func() {

			container.PutFactory(SimpleFactory("A1"), func(Doer) {})
			container.PutFactory(SimpleFactory("A2"), func(Doer) {})
			container.PutFactory(SimpleFactory("A3"), func(Doer) {})

			container.Put(&SliceInjected{})

			container.CallInjected(func(injected *SliceInjected) {
				Expect(injected.Simple).To(ConsistOf(
					&Simple{"A1"},
					&Simple{"A2"},
					&Simple{"A3"}))
			})

		})

		It("should inject a slice with nil components", func() {

			container.PutFactory(SimpleFactory("A1"))
			container.PutFactory(func() *Simple { return nil })
			container.PutFactory(SimpleFactory("A3"))

			container.Put(&SliceInjected{})

			container.CallInjected(func(injected *SliceInjected) {
				Expect(injected.Simple).To(ConsistOf(
					&Simple{"A1"},
					(*Simple)(nil),
					&Simple{"A3"}))
			})

		})

		It("should inject a slice of interface", func() {

			container.PutFactory(SimpleFactory("A1"), func(Doer) {})
			container.PutFactory(SimpleFactory("A2"), func(Doer) {})
			container.PutFactory(TrivialFactory("T1"), func(Doer) {})

			container.Put(&InterfaceSliceInjected{})

			container.CallInjected(func(injected *InterfaceSliceInjected) {
				Expect(injected.Doers).To(ConsistOf(
					&Simple{"A1"},
					&Simple{"A2"},
					Trivial("T1")))
			})

		})

		It("should inject an empty slice", func() {

			container.Put(&InterfaceSliceInjected{})

			container.CallInjected(func(injected *InterfaceSliceInjected) {
				Expect(injected.Doers).To(BeEmpty())
			})

		})

	})

})

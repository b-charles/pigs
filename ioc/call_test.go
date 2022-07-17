package ioc_test

import (
	. "github.com/b-charles/pigs/ioc"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("IOC call", func() {

	var container *Container

	BeforeEach(func() {
		container = NewContainer()
	})

	It("should be able to call without injection", func() {

		called := false
		Expect(container.CallInjected(func() {
			called = true
		})).To(Succeed())
		Expect(called).To(BeTrue())

	})

	Describe("with simple injection", func() {

		It("should be able to inject itself", func() {

			called := false
			Expect(container.CallInjected(func(injected *Container) {
				called = true
				Expect(injected).To(Equal(container))
			})).To(Succeed())
			Expect(called).To(BeTrue())

		})

		It("should inject a simple component", func() {

			container.PutFactory(SimpleFactory("A"))

			called := false
			Expect(container.CallInjected(func(a *Simple) {
				called = true
				Expect(a).To(Equal(&Simple{"A"}))
			})).To(Succeed())
			Expect(called).To(BeTrue())

		})

		It("should inject two simple components", func() {

			container.PutFactory(SimpleFactory("A"))
			container.PutFactory(TrivialFactory("B"))

			called := false
			Expect(container.CallInjected(func(a *Simple, b Trivial) {
				called = true
				Expect(a).To(Equal(&Simple{"A"}))
				Expect(b).To(Equal(Trivial("B")))
			})).To(Succeed())
			Expect(called).To(BeTrue())

		})

		It("should inject nil", func() {

			container.PutFactory(func() *Simple { return nil })

			called := false
			Expect(container.CallInjected(func(a *Simple) {
				called = true
				Expect(a).To(BeNil())
			})).To(Succeed())
			Expect(called).To(BeTrue())

		})

		It("should inject test component if provided", func() {

			container.PutFactory(SimpleFactory("A"))
			container.TestPutFactory(SimpleFactory("TEST"))

			called := false
			Expect(container.CallInjected(func(a *Simple) {
				called = true
				Expect(a).To(Equal(&Simple{"TEST"}))
			})).To(Succeed())
			Expect(called).To(BeTrue())

		})

		It("should return an error if no component is provided", func() {

			Expect(container.CallInjected(func(a *Simple) {
			})).To(HaveOccurred())

		})

		It("should return an error if too many components are provided", func() {

			container.PutFactory(SimpleFactory("A1"))
			container.PutFactory(SimpleFactory("A2"))

			Expect(container.CallInjected(func(a *Simple) {
			})).To(HaveOccurred())

		})

		It("should inject an interface", func() {

			container.PutFactory(SimpleFactory("A"), func(Doer) {})

			called := false
			Expect(container.CallInjected(func(a Doer) {
				called = true
				Expect(a).To(Equal(&Simple{"A"}))
			})).To(Succeed())
			Expect(called).To(BeTrue())

		})

		It("should inject not a struct", func() {

			container.PutFactory(TrivialFactory("A"), func(Doer) {})

			called := false
			Expect(container.CallInjected(func(a Doer) {
				called = true
				Expect(a).To(Equal(Trivial("A")))
			})).To(Succeed())
			Expect(called).To(BeTrue())

		})

		It("should restore 'core' configuration after a CallInjected", func() {

			container.PutFactory(SimpleFactory("A"), func(Doer) {})
			container.TestPutFactory(SimpleFactory("TEST"), func(Doer) {})

			// test

			called := false
			Expect(container.CallInjected(func(a Doer) {
				called = true
				Expect(a).To(Equal(&Simple{"TEST"}))
			})).To(Succeed(), "First call (test) should be ok.")
			Expect(called).To(BeTrue())

			// core

			called = false
			Expect(container.CallInjected(func(a Doer) {
				called = true
				Expect(a).To(Equal(&Simple{"A"}))
			})).To(Succeed(), "Second call (core) should be ok.")
			Expect(called).To(BeTrue())

		})

	})

	Describe("with slice injection", func() {

		It("should inject a slice of pointers", func() {

			container.PutFactory(SimpleFactory("A1"))
			container.PutFactory(SimpleFactory("A2"))
			container.PutFactory(SimpleFactory("A3"))

			called := false
			Expect(container.CallInjected(func(slice []*Simple) {
				called = true
				Expect(slice).To(ConsistOf(
					&Simple{"A1"},
					&Simple{"A2"},
					&Simple{"A3"}))
			})).To(Succeed())
			Expect(called).To(BeTrue())

		})

		It("should inject a slice with nil components", func() {

			container.PutFactory(SimpleFactory("A1"))
			container.PutFactory(func() *Simple { return nil })
			container.PutFactory(SimpleFactory("A3"))

			called := false
			Expect(container.CallInjected(func(slice []*Simple) {
				called = true
				Expect(slice).To(ConsistOf(
					(*Simple)(nil),
					&Simple{"A1"},
					&Simple{"A3"}))
			})).To(Succeed())
			Expect(called).To(BeTrue())

		})

		It("should inject a slice of interface", func() {

			container.PutFactory(SimpleFactory("A1"), func(Doer) {})
			container.PutFactory(SimpleFactory("A2"), func(Doer) {})
			container.PutFactory(SimpleFactory("A3"), func(Doer) {})

			called := false
			Expect(container.CallInjected(func(slice []Doer) {
				called = true
				Expect(slice).To(ConsistOf(
					&Simple{"A1"},
					&Simple{"A2"},
					&Simple{"A3"}))
			})).To(Succeed())
			Expect(called).To(BeTrue())

		})

		It("should inject a slice of interface (mixed)", func() {

			container.PutFactory(SimpleFactory("A1"), func(Doer) {})
			container.PutFactory(TrivialFactory("T2"), func(Doer) {})
			container.PutFactory(SimpleFactory("A2"), func(Doer) {})

			called := false
			Expect(container.CallInjected(func(slice []Doer) {
				called = true
				Expect(slice).To(ConsistOf(
					&Simple{"A1"},
					Trivial("T2"),
					&Simple{"A2"}))
			})).To(Succeed())
			Expect(called).To(BeTrue())

		})

		It("should inject an empty slice", func() {

			called := false
			Expect(container.CallInjected(func(slice []Doer) {
				called = true
				Expect(slice).To(BeEmpty())
			})).To(Succeed())
			Expect(called).To(BeTrue())

		})

		It("should inject only test components", func() {

			container.PutFactory(SimpleFactory("A1"))
			container.PutFactory(SimpleFactory("A2"))
			container.PutFactory(SimpleFactory("A3"))

			container.TestPutFactory(SimpleFactory("TEST"))

			called := false
			Expect(container.CallInjected(func(slice []*Simple) {
				called = true
				Expect(slice).To(ConsistOf(&Simple{"TEST"}))
			})).To(Succeed())
			Expect(called).To(BeTrue())

		})

	})

})

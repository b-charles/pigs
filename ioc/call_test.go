package ioc_test

import (
	. "github.com/b-charles/pigs/ioc"
	. "github.com/onsi/ginkgo/v2"
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

		It("should be able to inject status", func() {

			called := false
			Expect(container.CallInjected(func(injected *ContainerStatus) {
				called = true
				Expect(injected.String()).To(Not(BeEmpty()))
			})).To(Succeed())
			Expect(called).To(BeTrue())

		})

		It("should inject a simple component", func() {

			container.RegisterFactory(Core, "A", SimpleFactory("A"))

			called := false
			Expect(container.CallInjected(func(a *Simple) {
				called = true
				Expect(a).To(Equal(&Simple{"A"}))
			})).To(Succeed())
			Expect(called).To(BeTrue())

		})

		It("should inject two simple components", func() {

			container.RegisterFactory(Core, "A", SimpleFactory("A"))
			container.RegisterFactory(Core, "B", TrivialFactory("B"))

			called := false
			Expect(container.CallInjected(func(a *Simple, b Trivial) {
				called = true
				Expect(a).To(Equal(&Simple{"A"}))
				Expect(b).To(Equal(Trivial("B")))
			})).To(Succeed())
			Expect(called).To(BeTrue())

		})

		It("should inject default component if nothing is provided", func() {

			container.RegisterFactory(Def, "DEFAULT", SimpleFactory("DEFAULT"))

			called := false
			Expect(container.CallInjected(func(a *Simple) {
				called = true
				Expect(a).To(Equal(&Simple{"DEFAULT"}))
			})).To(Succeed())
			Expect(called).To(BeTrue())

		})

		It("should inject default component if default component is nil", func() {

			container.RegisterFactory(Def, "DEFAULT", SimpleFactory("DEFAULT"))
			container.RegisterFactory(Core, "NIL", func() *Simple { return nil })

			called := false
			Expect(container.CallInjected(func(a *Simple) {
				called = true
				Expect(a).To(Equal(&Simple{"DEFAULT"}))
			})).To(Succeed())
			Expect(called).To(BeTrue())

		})

		It("should inject core component even if default component is provided", func() {

			container.RegisterFactory(Def, "DEFAULT", SimpleFactory("DEFAULT"))
			container.RegisterFactory(Core, "A", SimpleFactory("A"))

			called := false
			Expect(container.CallInjected(func(a *Simple) {
				called = true
				Expect(a).To(Equal(&Simple{"A"}))
			})).To(Succeed())
			Expect(called).To(BeTrue())

		})

		It("should inject core component if test component is nil", func() {

			container.RegisterFactory(Core, "A", SimpleFactory("A"))
			container.RegisterFactory(Test, "NIL", func() *Simple { return nil })

			called := false
			Expect(container.CallInjected(func(a *Simple) {
				called = true
				Expect(a).To(Equal(&Simple{"A"}))
			})).To(Succeed())
			Expect(called).To(BeTrue())

		})

		It("should inject test component if provided", func() {

			container.RegisterFactory(Core, "A", SimpleFactory("A"))
			container.RegisterFactory(Test, "TEST", SimpleFactory("TEST"))

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

		It("should return an error if the only component is nil", func() {

			container.RegisterFactory(Core, "NIL", func() *Simple { return nil })

			Expect(container.CallInjected(func(a *Simple) {
			})).To(HaveOccurred())

		})

		It("should return an error if too many components are provided", func() {

			container.RegisterFactory(Core, "A1", SimpleFactory("A1"))
			container.RegisterFactory(Core, "A2", SimpleFactory("A2"))

			Expect(container.CallInjected(func(a *Simple) {
			})).To(HaveOccurred())

		})

		It("should inject an interface", func() {

			container.RegisterFactory(Core, "A", SimpleFactory("A"), func(Doer) {})

			called := false
			Expect(container.CallInjected(func(a Doer) {
				called = true
				Expect(a).To(Equal(&Simple{"A"}))
			})).To(Succeed())
			Expect(called).To(BeTrue())

		})

		It("should inject not a struct", func() {

			container.RegisterFactory(Core, "A", TrivialFactory("A"), func(Doer) {})

			called := false
			Expect(container.CallInjected(func(a Doer) {
				called = true
				Expect(a).To(Equal(Trivial("A")))
			})).To(Succeed())
			Expect(called).To(BeTrue())

		})

		It("should restore 'core' configuration after a CallInjected", func() {

			container.RegisterFactory(Core, "A", SimpleFactory("A"), func(Doer) {})
			container.RegisterFactory(Test, "TEST", SimpleFactory("TEST"), func(Doer) {})

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

			container.RegisterFactory(Core, "A1", SimpleFactory("A1"))
			container.RegisterFactory(Core, "A2", SimpleFactory("A2"))
			container.RegisterFactory(Core, "A3", SimpleFactory("A3"))

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

		It("should inject a slice with discarding nil components", func() {

			container.RegisterFactory(Core, "A1", SimpleFactory("A1"))
			container.RegisterFactory(Core, "NIL", func() *Simple { return nil })
			container.RegisterFactory(Core, "A3", SimpleFactory("A3"))

			called := false
			Expect(container.CallInjected(func(slice []*Simple) {
				called = true
				Expect(slice).To(ConsistOf(
					&Simple{"A1"},
					&Simple{"A3"}))
			})).To(Succeed())
			Expect(called).To(BeTrue())

		})

		It("should inject an empty slice if all components are nil", func() {

			container.RegisterFactory(Core, "NIL1", func() *Simple { return nil })
			container.RegisterFactory(Core, "NIL2", func() *Simple { return nil })

			called := false
			Expect(container.CallInjected(func(slice []*Simple) {
				called = true
				Expect(slice).To(BeEmpty())
			})).To(Succeed())
			Expect(called).To(BeTrue())

		})

		It("should inject a slice of interface", func() {

			container.RegisterFactory(Core, "A1", SimpleFactory("A1"), func(Doer) {})
			container.RegisterFactory(Core, "A2", SimpleFactory("A2"), func(Doer) {})
			container.RegisterFactory(Core, "A3", SimpleFactory("A3"), func(Doer) {})

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

			container.RegisterFactory(Core, "A1", SimpleFactory("A1"), func(Doer) {})
			container.RegisterFactory(Core, "T2", TrivialFactory("T2"), func(Doer) {})
			container.RegisterFactory(Core, "A2", SimpleFactory("A2"), func(Doer) {})

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

			container.RegisterFactory(Core, "A1", SimpleFactory("A1"))
			container.RegisterFactory(Core, "A2", SimpleFactory("A2"))
			container.RegisterFactory(Core, "A3", SimpleFactory("A3"))

			container.RegisterFactory(Test, "TEST", SimpleFactory("TEST"))

			called := false
			Expect(container.CallInjected(func(slice []*Simple) {
				called = true
				Expect(slice).To(ConsistOf(&Simple{"TEST"}))
			})).To(Succeed())
			Expect(called).To(BeTrue())

		})

	})

})

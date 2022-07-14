package ioc_test

import (
	. "github.com/b-charles/pigs/ioc"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("IOC call", func() {

	var (
		container *Container
	)

	BeforeEach(func() {
		container = NewContainer()
	})

	Describe("with simple injection", func() {

		It("should be able to get itself", func() {

			var injectedContainer *Container

			Expect(container.CallInjected(func(injected struct {
				Container *Container
			}) {
				injectedContainer = injected.Container
			})).Should(Succeed())

			Expect(injectedContainer).To(Equal(container))

		})

		It("should be able to get itself by name", func() {

			var injectedContainer *Container

			Expect(container.CallInjected(func(injected struct {
				Container *Container `inject:"github.com/b-charles/pigs/ioc/Container"`
			}) {
				injectedContainer = injected.Container
			})).Should(Succeed())

			Expect(injectedContainer).To(Equal(container))

		})

		It("should inject a simple component", func() {

			container.PutNamedFactory(SimpleStructFactory("A"), "A")

			var injectedA *SimpleStruct

			Expect(container.CallInjected(func(injected struct {
				A *SimpleStruct
			}) {
				injectedA = injected.A
			})).Should(Succeed())

			Expect(injectedA).To(Equal(&SimpleStruct{"A"}))

		})

		It("should inject a simple deferenced component", func() {

			container.PutNamedFactory(SimpleStructFactory("A"), "A")

			var injectedA SimpleStruct

			Expect(container.CallInjected(func(injected struct {
				A SimpleStruct
			}) {
				injectedA = injected.A
			})).Should(Succeed())

			Expect(injectedA).To(Equal(SimpleStruct{"A"}))

		})

		It("should select not nil component", func() {

			container.PutNamedFactory(func() *SimpleStruct { return nil }, "NIL", "A")
			container.PutNamedFactory(SimpleStructFactory("A"), "NOT_NIL", "A")

			var injectedA *SimpleStruct

			Expect(container.CallInjected(func(injected struct {
				A *SimpleStruct
			}) {
				injectedA = injected.A
			})).Should(Succeed())

			Expect(injectedA).To(Equal(&SimpleStruct{"A"}))

		})

		It("should inject test component if provided", func() {

			container.PutNamedFactory(SimpleStructFactory("A"), "A")
			container.TestPutNamedFactory(SimpleStructFactory("TEST"), "TEST", "A")

			var injectedA *SimpleStruct

			Expect(container.CallInjected(func(injected struct {
				A *SimpleStruct
			}) {
				injectedA = injected.A
			})).Should(Succeed())

			Expect(injectedA).To(Equal(&SimpleStruct{"TEST"}))

		})

		It("should return an error if no component is provided", func() {

			var injectedA *SimpleStruct

			Expect(container.CallInjected(func(injected struct {
				A *SimpleStruct
			}) {
				injectedA = injected.A
				injectedA.doSomething()
			})).Should(HaveOccurred())

		})

		It("should return an error if too many components are provided", func() {

			container.PutNamedFactory(SimpleStructFactory("A1"), "A1", "A")
			container.PutNamedFactory(SimpleStructFactory("A2"), "A2", "A")

			var injectedA *SimpleStruct

			Expect(container.CallInjected(func(injected struct {
				A *SimpleStruct
			}) {
				injectedA = injected.A
				injectedA.doSomething()
			})).Should(HaveOccurred())

		})

		It("should inject an interface", func() {

			container.PutNamedFactory(SimpleStructFactory("A"), "A")

			var injectedA SomethingDoer

			Expect(container.CallInjected(func(injected struct {
				A SomethingDoer
			}) {
				injectedA = injected.A
			})).Should(Succeed())

			Expect(injectedA).To(Equal(&SimpleStruct{"A"}))

		})

		It("should restore 'core' configuration after a TestClear", func() {

			container.PutNamedFactory(SimpleStructFactory("A"), "A")
			container.TestPutNamedFactory(SimpleStructFactory("TEST"), "TEST", "A")

			var injectedA *SimpleStruct

			// test

			Expect(container.CallInjected(func(injected struct {
				A *SimpleStruct
			}) {
				injectedA = injected.A
			})).Should(Succeed())

			Expect(injectedA).To(Equal(&SimpleStruct{"TEST"}))

			// core

			Expect(container.CallInjected(func(injected struct {
				A *SimpleStruct
			}) {
				injectedA = injected.A
			})).Should(Succeed())

			Expect(injectedA).To(Equal(&SimpleStruct{"A"}))

		})

	})

	Describe("with slice injection", func() {

		It("should inject a slice of pointers", func() {

			container.PutNamedFactory(SimpleStructFactory("A1"), "A1", "A")
			container.PutNamedFactory(SimpleStructFactory("A2"), "A2", "A")
			container.PutNamedFactory(SimpleStructFactory("A3"), "A3", "A")

			var injectedA []*SimpleStruct

			Expect(container.CallInjected(func(injected struct {
				A []*SimpleStruct
			}) {
				injectedA = injected.A
			})).Should(Succeed())

			Expect(injectedA).Should(ConsistOf(
				&SimpleStruct{"A1"},
				&SimpleStruct{"A2"},
				&SimpleStruct{"A3"}))

		})

		It("should inject a slice of values", func() {

			container.PutNamedFactory(SimpleStructFactory("A1"), "A1", "A")
			container.PutNamedFactory(SimpleStructFactory("A2"), "A2", "A")
			container.PutNamedFactory(SimpleStructFactory("A3"), "A3", "A")

			var injectedA []SimpleStruct

			Expect(container.CallInjected(func(injected struct {
				A []SimpleStruct
			}) {
				injectedA = injected.A
			})).Should(Succeed())

			Expect(injectedA).Should(ConsistOf(
				SimpleStruct{"A1"},
				SimpleStruct{"A2"},
				SimpleStruct{"A3"}))

		})

		It("should inject a slice without nil components", func() {

			container.PutNamedFactory(func() *SimpleStruct {
				return nil
			}, "NIL", "A")

			container.PutNamedFactory(SimpleStructFactory("A1"), "A1", "A")
			container.PutNamedFactory(SimpleStructFactory("A2"), "A2", "A")
			container.PutNamedFactory(SimpleStructFactory("A3"), "A3", "A")

			var injectedA []*SimpleStruct

			Expect(container.CallInjected(func(injected struct {
				A []*SimpleStruct
			}) {
				injectedA = injected.A
			})).Should(Succeed())

			Expect(injectedA).Should(ConsistOf(
				&SimpleStruct{"A1"},
				&SimpleStruct{"A2"},
				&SimpleStruct{"A3"}))

		})

		It("should inject a slice of interface", func() {

			container.PutNamedFactory(SimpleStructFactory("A1"), "A1", "A")
			container.PutNamedFactory(SimpleStructFactory("A2"), "A2", "A")
			container.PutNamedFactory(SimpleStructFactory("A3"), "A3", "A")

			var injectedA []SomethingDoer

			Expect(container.CallInjected(func(injected struct {
				A []SomethingDoer
			}) {
				injectedA = injected.A
			})).Should(Succeed())

			Expect(injectedA).Should(ConsistOf(
				&SimpleStruct{"A1"},
				&SimpleStruct{"A2"},
				&SimpleStruct{"A3"}))

		})

	})

	Describe("with map injection", func() {

		It("should inject a map of pointers", func() {

			container.PutNamedFactory(SimpleStructFactory("A1"), "A1", "A")
			container.PutNamedFactory(SimpleStructFactory("A2"), "A2", "A")
			container.PutNamedFactory(SimpleStructFactory("A3"), "A3", "A")

			var injectedA map[string]*SimpleStruct

			Expect(container.CallInjected(func(injected struct {
				A map[string]*SimpleStruct
			}) {
				injectedA = injected.A
			})).Should(Succeed())

			Expect(injectedA).Should(HaveLen(3))
			Expect(injectedA).Should(HaveKeyWithValue("A1", &SimpleStruct{"A1"}))
			Expect(injectedA).Should(HaveKeyWithValue("A2", &SimpleStruct{"A2"}))
			Expect(injectedA).Should(HaveKeyWithValue("A3", &SimpleStruct{"A3"}))

		})

		It("should inject a map of values", func() {

			container.PutNamedFactory(SimpleStructFactory("A1"), "A1", "A")
			container.PutNamedFactory(SimpleStructFactory("A2"), "A2", "A")
			container.PutNamedFactory(SimpleStructFactory("A3"), "A3", "A")

			var injectedA map[string]SimpleStruct

			Expect(container.CallInjected(func(injected struct {
				A map[string]SimpleStruct
			}) {
				injectedA = injected.A
			})).Should(Succeed())

			Expect(injectedA).Should(HaveLen(3))
			Expect(injectedA).Should(HaveKeyWithValue("A1", SimpleStruct{"A1"}))
			Expect(injectedA).Should(HaveKeyWithValue("A2", SimpleStruct{"A2"}))
			Expect(injectedA).Should(HaveKeyWithValue("A3", SimpleStruct{"A3"}))

		})

		It("should inject a map without nil components", func() {

			container.PutNamedFactory(func() *SimpleStruct {
				return nil
			}, "NIL", "A")

			container.PutNamedFactory(SimpleStructFactory("A1"), "A1", "A")
			container.PutNamedFactory(SimpleStructFactory("A2"), "A2", "A")
			container.PutNamedFactory(SimpleStructFactory("A3"), "A3", "A")

			var injectedA map[string]*SimpleStruct

			Expect(container.CallInjected(func(injected struct {
				A map[string]*SimpleStruct
			}) {
				injectedA = injected.A
			})).Should(Succeed())

			Expect(injectedA).Should(HaveLen(3))
			Expect(injectedA).Should(HaveKeyWithValue("A1", &SimpleStruct{"A1"}))
			Expect(injectedA).Should(HaveKeyWithValue("A2", &SimpleStruct{"A2"}))
			Expect(injectedA).Should(HaveKeyWithValue("A3", &SimpleStruct{"A3"}))

		})

		It("should inject a map of interface", func() {

			container.PutNamedFactory(SimpleStructFactory("A1"), "A1", "A")
			container.PutNamedFactory(SimpleStructFactory("A2"), "A2", "A")
			container.PutNamedFactory(SimpleStructFactory("A3"), "A3", "A")

			var injectedA map[string]SomethingDoer

			Expect(container.CallInjected(func(injected struct {
				A map[string]SomethingDoer
			}) {
				injectedA = injected.A
			})).Should(Succeed())

			Expect(injectedA).Should(HaveLen(3))
			Expect(injectedA).Should(HaveKeyWithValue("A1", &SimpleStruct{"A1"}))
			Expect(injectedA).Should(HaveKeyWithValue("A2", &SimpleStruct{"A2"}))
			Expect(injectedA).Should(HaveKeyWithValue("A3", &SimpleStruct{"A3"}))

		})

	})

})

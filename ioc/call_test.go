package ioc_test

import (
	. "github.com/b-charles/pigs/ioc"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("IOC call", func() {

	var (
		container *Container
		err       error
	)

	BeforeEach(func() {
		container = NewContainer()
	})

	Describe("with simple injection", func() {

		It("should be able to inject itself", func() {

			err = container.CallInjected(func(injected *Container) {
				Expect(injected).To(Equal(container))
			})
			Expect(err).To(Succeed())

		})

		It("should be able to inject itself with an injected struct (without name)", func() {

			err = container.CallInjected(func(injected struct {
				Container *Container
			}) {
				Expect(injected.Container).To(Equal(container))
			})
			Expect(err).To(Succeed())

		})

		It("should be able to inject itself with an injected struct (with name)", func() {

			err = container.CallInjected(func(injected struct {
				Container *Container `inject:"github.com/b-charles/pigs/ioc/Container"`
			}) {
				Expect(injected.Container).To(Equal(container))
			})
			Expect(err).To(Succeed())

		})

		It("should inject a simple component (direct)", func() {

			err = container.PutFactory(SimpleStructFactory("A"))
			Expect(err).To(Succeed())

			err = container.CallInjected(func(a *SimpleStruct) {
				Expect(a).To(Equal(&SimpleStruct{"A"}))
			})
			Expect(err).To(Succeed())

		})

		It("should inject a simple component (injected struct without name)", func() {

			err = container.PutFactory(SimpleStructFactory("A"))
			Expect(err).To(Succeed())

			err = container.CallInjected(func(injected struct {
				Simple *SimpleStruct
			}) {
				Expect(injected.Simple).To(Equal(&SimpleStruct{"A"}))
			})
			Expect(err).To(Succeed())

		})

		It("should inject a simple component (injected struct with name)", func() {

			err = container.PutFactory(SimpleStructFactory("A"))
			Expect(err).To(Succeed())

			err = container.CallInjected(func(injected struct {
				Simple *SimpleStruct `inject:"github.com/b-charles/pigs/ioc_test/SimpleStruct"`
			}) {
				Expect(injected.Simple).To(Equal(&SimpleStruct{"A"}))
			})
			Expect(err).To(Succeed())

		})

		It("should inject a simple component (injected struct by field name)", func() {

			err = container.PutNamedFactory(SimpleStructFactory("A"), "A")
			Expect(err).To(Succeed())

			err = container.CallInjected(func(injected struct {
				A *SimpleStruct
			}) {
				Expect(injected.A).To(Equal(&SimpleStruct{"A"}))
			})
			Expect(err).To(Succeed())

		})

		It("should inject a simple deferenced component (direct)", func() {

			err = container.PutFactory(SimpleStructFactory("A"))
			Expect(err).To(Succeed())

			err = container.CallInjected(func(a SimpleStruct) {
				Expect(a).To(Equal(SimpleStruct{"A"}))
			})
			Expect(err).To(Succeed())

		})

		It("should inject a simple deferenced component (in a struct)", func() {

			err = container.PutFactory(SimpleStructFactory("A"))
			Expect(err).To(Succeed())

			err = container.CallInjected(func(injected struct {
				A SimpleStruct
			}) {
				Expect(injected.A).To(Equal(SimpleStruct{"A"}))
			})
			Expect(err).To(Succeed())

		})

		It("should select not nil component", func() {

			err = container.PutNamedFactory(func() *SimpleStruct { return nil }, "NIL", func(SomethingDoer) {})
			Expect(err).To(Succeed())

			err = container.PutNamedFactory(SimpleStructFactory("A"), "NOT_NIL", func(SomethingDoer) {})
			Expect(err).To(Succeed())

			err = container.CallInjected(func(a SomethingDoer) {
				Expect(a).To(Equal(&SimpleStruct{"A"}))
			})
			Expect(err).To(Succeed())

		})

		It("should inject test component if provided", func() {

			err = container.PutFactory(SimpleStructFactory("A"))
			Expect(err).To(Succeed())

			err = container.TestPutFactory(SimpleStructFactory("TEST"))
			Expect(err).To(Succeed())

			err = container.CallInjected(func(a *SimpleStruct) {
				Expect(a).To(Equal(&SimpleStruct{"TEST"}))
			})
			Expect(err).To(Succeed())

		})

		It("should return an error if no component is provided", func() {

			var injectedA *SimpleStruct

			Expect(container.CallInjected(func(injected struct {
				A *SimpleStruct
			}) {
				injectedA = injected.A
				injectedA.doSomething()
			})).To(HaveOccurred())

		})

		It("should return an error if too many components are provided", func() {

			Expect(container.PutNamedFactory(SimpleStructFactory("A1"), "A1", "A")).To(Succeed())
			Expect(container.PutNamedFactory(SimpleStructFactory("A2"), "A2", "A")).To(Succeed())

			var injectedA *SimpleStruct

			Expect(container.CallInjected(func(injected struct {
				A *SimpleStruct
			}) {
				injectedA = injected.A
				injectedA.doSomething()
			})).To(HaveOccurred())

		})

		It("should inject an interface", func() {

			Expect(container.PutNamedFactory(SimpleStructFactory("A"), "A")).To(Succeed())

			var injectedA SomethingDoer

			Expect(container.CallInjected(func(injected struct {
				A SomethingDoer
			}) {
				injectedA = injected.A
			})).To(Succeed())

			Expect(injectedA).To(Equal(&SimpleStruct{"A"}))

		})

		It("should restore 'core' configuration after a TestClear", func() {

			Expect(container.PutNamedFactory(SimpleStructFactory("A"), "A")).To(Succeed())
			Expect(container.TestPutNamedFactory(SimpleStructFactory("TEST"), "TEST", "A")).To(Succeed())

			var injectedA *SimpleStruct

			// test

			Expect(container.CallInjected(func(injected struct {
				A *SimpleStruct
			}) {
				injectedA = injected.A
			})).To(Succeed())

			Expect(injectedA).To(Equal(&SimpleStruct{"TEST"}))

			// core

			Expect(container.CallInjected(func(injected struct {
				A *SimpleStruct
			}) {
				injectedA = injected.A
			})).To(Succeed())

			Expect(injectedA).To(Equal(&SimpleStruct{"A"}))

		})

	})

	Describe("with slice injection", func() {

		It("should inject a slice of pointers", func() {

			Expect(container.PutNamedFactory(SimpleStructFactory("A1"), "A1", "A")).To(Succeed())
			Expect(container.PutNamedFactory(SimpleStructFactory("A2"), "A2", "A")).To(Succeed())
			Expect(container.PutNamedFactory(SimpleStructFactory("A3"), "A3", "A")).To(Succeed())

			var injectedA []*SimpleStruct

			Expect(container.CallInjected(func(injected struct {
				A []*SimpleStruct
			}) {
				injectedA = injected.A
			})).To(Succeed())

			Expect(injectedA).To(ConsistOf(
				&SimpleStruct{"A1"},
				&SimpleStruct{"A2"},
				&SimpleStruct{"A3"}))

		})

		It("should inject a slice of values", func() {

			Expect(container.PutNamedFactory(SimpleStructFactory("A1"), "A1", "A")).To(Succeed())
			Expect(container.PutNamedFactory(SimpleStructFactory("A2"), "A2", "A")).To(Succeed())
			Expect(container.PutNamedFactory(SimpleStructFactory("A3"), "A3", "A")).To(Succeed())

			var injectedA []SimpleStruct

			Expect(container.CallInjected(func(injected struct {
				A []SimpleStruct
			}) {
				injectedA = injected.A
			})).To(Succeed())

			Expect(injectedA).To(ConsistOf(
				SimpleStruct{"A1"},
				SimpleStruct{"A2"},
				SimpleStruct{"A3"}))

		})

		It("should inject a slice without nil components", func() {

			Expect(container.PutNamedFactory(func() *SimpleStruct {
				return nil
			}, "NIL", "A")).To(Succeed())

			Expect(container.PutNamedFactory(SimpleStructFactory("A1"), "A1", "A")).To(Succeed())
			Expect(container.PutNamedFactory(SimpleStructFactory("A2"), "A2", "A")).To(Succeed())
			Expect(container.PutNamedFactory(SimpleStructFactory("A3"), "A3", "A")).To(Succeed())

			var injectedA []*SimpleStruct

			Expect(container.CallInjected(func(injected struct {
				A []*SimpleStruct
			}) {
				injectedA = injected.A
			})).To(Succeed())

			Expect(injectedA).To(ConsistOf(
				&SimpleStruct{"A1"},
				&SimpleStruct{"A2"},
				&SimpleStruct{"A3"}))

		})

		It("should inject a slice of interface", func() {

			Expect(container.PutNamedFactory(SimpleStructFactory("A1"), "A1", "A")).To(Succeed())
			Expect(container.PutNamedFactory(SimpleStructFactory("A2"), "A2", "A")).To(Succeed())
			Expect(container.PutNamedFactory(SimpleStructFactory("A3"), "A3", "A")).To(Succeed())

			var injectedA []SomethingDoer

			Expect(container.CallInjected(func(injected struct {
				A []SomethingDoer
			}) {
				injectedA = injected.A
			})).To(Succeed())

			Expect(injectedA).To(ConsistOf(
				&SimpleStruct{"A1"},
				&SimpleStruct{"A2"},
				&SimpleStruct{"A3"}))

		})

	})

	Describe("with map injection", func() {

		It("should inject a map of pointers", func() {

			Expect(container.PutNamedFactory(SimpleStructFactory("A1"), "A1", "A")).To(Succeed())
			Expect(container.PutNamedFactory(SimpleStructFactory("A2"), "A2", "A")).To(Succeed())
			Expect(container.PutNamedFactory(SimpleStructFactory("A3"), "A3", "A")).To(Succeed())

			var injectedA map[string]*SimpleStruct

			Expect(container.CallInjected(func(injected struct {
				A map[string]*SimpleStruct
			}) {
				injectedA = injected.A
			})).To(Succeed())

			Expect(injectedA).To(HaveLen(3))
			Expect(injectedA).To(HaveKeyWithValue("A1", &SimpleStruct{"A1"}))
			Expect(injectedA).To(HaveKeyWithValue("A2", &SimpleStruct{"A2"}))
			Expect(injectedA).To(HaveKeyWithValue("A3", &SimpleStruct{"A3"}))

		})

		It("should inject a map of values", func() {

			Expect(container.PutNamedFactory(SimpleStructFactory("A1"), "A1", "A")).To(Succeed())
			Expect(container.PutNamedFactory(SimpleStructFactory("A2"), "A2", "A")).To(Succeed())
			Expect(container.PutNamedFactory(SimpleStructFactory("A3"), "A3", "A")).To(Succeed())

			var injectedA map[string]SimpleStruct

			Expect(container.CallInjected(func(injected struct {
				A map[string]SimpleStruct
			}) {
				injectedA = injected.A
			})).To(Succeed())

			Expect(injectedA).To(HaveLen(3))
			Expect(injectedA).To(HaveKeyWithValue("A1", SimpleStruct{"A1"}))
			Expect(injectedA).To(HaveKeyWithValue("A2", SimpleStruct{"A2"}))
			Expect(injectedA).To(HaveKeyWithValue("A3", SimpleStruct{"A3"}))

		})

		It("should inject a map without nil components", func() {

			Expect(container.PutNamedFactory(func() *SimpleStruct {
				return nil
			}, "NIL", "A")).To(Succeed())

			Expect(container.PutNamedFactory(SimpleStructFactory("A1"), "A1", "A")).To(Succeed())
			Expect(container.PutNamedFactory(SimpleStructFactory("A2"), "A2", "A")).To(Succeed())
			Expect(container.PutNamedFactory(SimpleStructFactory("A3"), "A3", "A")).To(Succeed())

			var injectedA map[string]*SimpleStruct

			Expect(container.CallInjected(func(injected struct {
				A map[string]*SimpleStruct
			}) {
				injectedA = injected.A
			})).To(Succeed())

			Expect(injectedA).To(HaveLen(3))
			Expect(injectedA).To(HaveKeyWithValue("A1", &SimpleStruct{"A1"}))
			Expect(injectedA).To(HaveKeyWithValue("A2", &SimpleStruct{"A2"}))
			Expect(injectedA).To(HaveKeyWithValue("A3", &SimpleStruct{"A3"}))

		})

		It("should inject a map of interface", func() {

			Expect(container.PutNamedFactory(SimpleStructFactory("A1"), "A1", "A")).To(Succeed())
			Expect(container.PutNamedFactory(SimpleStructFactory("A2"), "A2", "A")).To(Succeed())
			Expect(container.PutNamedFactory(SimpleStructFactory("A3"), "A3", "A")).To(Succeed())

			var injectedA map[string]SomethingDoer

			Expect(container.CallInjected(func(injected struct {
				A map[string]SomethingDoer
			}) {
				injectedA = injected.A
			})).To(Succeed())

			Expect(injectedA).To(HaveLen(3))
			Expect(injectedA).To(HaveKeyWithValue("A1", &SimpleStruct{"A1"}))
			Expect(injectedA).To(HaveKeyWithValue("A2", &SimpleStruct{"A2"}))
			Expect(injectedA).To(HaveKeyWithValue("A3", &SimpleStruct{"A3"}))

		})

	})

})

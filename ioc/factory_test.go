package ioc_test

import (
	. "github.com/b-charles/pigs/ioc"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("IOC factory", func() {

	var (
		container *Container
	)

	BeforeEach(func() {
		container = NewContainer()
	})

	Describe("with simple injection", func() {

		It("should create a component with dependency", func() {

			container.PutFactory(SimpleStructFactory("A"), "A")
			container.PutFactory(InjectedStructFactory, "B")

			var injectedB *InjectedStruct

			Expect(container.CallInjected(func(injected struct {
				B *InjectedStruct
			}) {
				injectedB = injected.B
			})).Should(Succeed())

			Expect(injectedB).To(Equal(&InjectedStruct{&SimpleStruct{"A"}}))

		})

		It("should create a component with dependency, another form of factory", func() {

			container.PutFactory(SimpleStructFactory("SimpleStruct"), "SimpleStruct")

			container.PutFactory(func(A *SimpleStruct) *InjectedStruct {
				return &InjectedStruct{A}
			}, "InjectedStruct")

			var injectedB *InjectedStruct

			Expect(container.CallInjected(func(injected *InjectedStruct) {
				injectedB = injected
			})).Should(Succeed())

			Expect(injectedB).To(Equal(&InjectedStruct{&SimpleStruct{"SimpleStruct"}}))

		})

		It("should select not nil component", func() {

			container.PutFactory(func() *SimpleStruct {
				return nil
			}, "NIL", "A")

			container.PutFactory(SimpleStructFactory("NOT_NIL"), "NOT_NIL", "A")
			container.PutFactory(InjectedStructFactory, "B")

			var injectedB *InjectedStruct

			Expect(container.CallInjected(func(injected struct {
				B *InjectedStruct
			}) {
				injectedB = injected.B
			})).Should(Succeed())

			Expect(injectedB).To(Equal(&InjectedStruct{&SimpleStruct{"NOT_NIL"}}))

		})

		It("should inject test component if provided", func() {

			container.PutFactory(SimpleStructFactory("A"), "A")
			container.TestPutFactory(SimpleStructFactory("TEST"), "TEST", "A")

			container.PutFactory(InjectedStructFactory, "B")

			var injectedB *InjectedStruct

			Expect(container.CallInjected(func(injected struct {
				B *InjectedStruct
			}) {
				injectedB = injected.B
			})).Should(Succeed())

			Expect(injectedB).To(Equal(&InjectedStruct{&SimpleStruct{"TEST"}}))

		})

		It("should return an error if no component is provided", func() {

			container.PutFactory(InjectedStructFactory, "B")

			var injectedB *InjectedStruct

			Expect(container.CallInjected(func(injected struct {
				B *InjectedStruct
			}) {
				injectedB = injected.B
				injectedB.doSomething()
			})).Should(HaveOccurred())

		})

		It("should return an error if too many components are provided", func() {

			container.PutFactory(SimpleStructFactory("A1"), "A1", "A")
			container.PutFactory(SimpleStructFactory("A2"), "A2", "A")

			container.PutFactory(InjectedStructFactory, "B")

			var injectedB *InjectedStruct

			Expect(container.CallInjected(func(injected struct {
				B *InjectedStruct
			}) {
				injectedB = injected.B
				injectedB.doSomething()
			})).Should(HaveOccurred())

		})

		It("should inject an interface", func() {

			container.PutFactory(SimpleStructFactory("A"), "A")

			container.PutFactory(InterfaceInjectedStructFactory, "B")

			var injectedB *InterfaceInjectedStruct

			Expect(container.CallInjected(func(injected struct {
				B *InterfaceInjectedStruct
			}) {
				injectedB = injected.B
			})).Should(Succeed())

			Expect(injectedB).To(Equal(&InterfaceInjectedStruct{&SimpleStruct{"A"}}))

		})

	})

	Describe("with slice injection", func() {

		It("should inject a slice", func() {

			container.PutFactory(SimpleStructFactory("A1"), "A1", "A")
			container.PutFactory(SimpleStructFactory("A2"), "A2", "A")
			container.PutFactory(SimpleStructFactory("A3"), "A3", "A")

			container.PutFactory(SliceInjectedStructFactory, "B")

			var injectedB *SliceInjectedStruct

			Expect(container.CallInjected(func(injected struct {
				B *SliceInjectedStruct
			}) {
				injectedB = injected.B
			})).Should(Succeed())

			Expect(injectedB.SimpleStructs).Should(ConsistOf(
				&SimpleStruct{"A1"},
				&SimpleStruct{"A2"},
				&SimpleStruct{"A3"}))

		})

		It("should inject a slice without nil components", func() {

			container.PutFactory(func() *SimpleStruct {
				return nil
			}, "NIL", "A")

			container.PutFactory(SimpleStructFactory("A1"), "A1", "A")
			container.PutFactory(SimpleStructFactory("A2"), "A2", "A")
			container.PutFactory(SimpleStructFactory("A3"), "A3", "A")

			container.PutFactory(SliceInjectedStructFactory, "B")

			var injectedB *SliceInjectedStruct

			Expect(container.CallInjected(func(injected struct {
				B *SliceInjectedStruct
			}) {
				injectedB = injected.B
			})).Should(Succeed())

			Expect(injectedB.SimpleStructs).Should(ConsistOf(
				&SimpleStruct{"A1"},
				&SimpleStruct{"A2"},
				&SimpleStruct{"A3"}))

		})

		It("should inject a slice of interface", func() {

			container.PutFactory(SimpleStructFactory("A1"), "A1", "A")
			container.PutFactory(SimpleStructFactory("A2"), "A2", "A")
			container.PutFactory(SimpleStructFactory("A3"), "A3", "A")

			container.PutFactory(InterfaceSliceInjectedStructFactory, "B")

			var injectedB *InterfaceSliceInjectedStruct

			Expect(container.CallInjected(func(injected struct {
				B *InterfaceSliceInjectedStruct
			}) {
				injectedB = injected.B
			})).Should(Succeed())

			Expect(injectedB.Doers).Should(ConsistOf(
				&SimpleStruct{"A1"},
				&SimpleStruct{"A2"},
				&SimpleStruct{"A3"}))

		})

	})

	Describe("with map injection", func() {

		It("should inject a map", func() {

			container.PutFactory(SimpleStructFactory("A1"), "A1", "A")
			container.PutFactory(SimpleStructFactory("A2"), "A2", "A")
			container.PutFactory(SimpleStructFactory("A3"), "A3", "A")

			container.PutFactory(MapInjectedStructFactory, "B")

			var injectedB *MapInjectedStruct

			Expect(container.CallInjected(func(injected struct {
				B *MapInjectedStruct
			}) {
				injectedB = injected.B
			})).Should(Succeed())

			Expect(injectedB.SimpleStructs).Should(HaveLen(3))
			Expect(injectedB.SimpleStructs).Should(HaveKeyWithValue("A1", &SimpleStruct{"A1"}))
			Expect(injectedB.SimpleStructs).Should(HaveKeyWithValue("A2", &SimpleStruct{"A2"}))
			Expect(injectedB.SimpleStructs).Should(HaveKeyWithValue("A3", &SimpleStruct{"A3"}))

		})

		It("should inject a map without nil components", func() {

			container.PutFactory(func() *SimpleStruct {
				return nil
			}, "NIL", "A")

			container.PutFactory(SimpleStructFactory("A1"), "A1", "A")
			container.PutFactory(SimpleStructFactory("A2"), "A2", "A")
			container.PutFactory(SimpleStructFactory("A3"), "A3", "A")

			container.PutFactory(MapInjectedStructFactory, "B")

			var injectedB *MapInjectedStruct

			Expect(container.CallInjected(func(injected struct {
				B *MapInjectedStruct
			}) {
				injectedB = injected.B
			})).Should(Succeed())

			Expect(injectedB.SimpleStructs).Should(HaveLen(3))
			Expect(injectedB.SimpleStructs).Should(HaveKeyWithValue("A1", &SimpleStruct{"A1"}))
			Expect(injectedB.SimpleStructs).Should(HaveKeyWithValue("A2", &SimpleStruct{"A2"}))
			Expect(injectedB.SimpleStructs).Should(HaveKeyWithValue("A3", &SimpleStruct{"A3"}))

		})

		It("should inject a map of interface", func() {

			container.PutFactory(SimpleStructFactory("A1"), "A1", "A")
			container.PutFactory(SimpleStructFactory("A2"), "A2", "A")
			container.PutFactory(SimpleStructFactory("A3"), "A3", "A")

			container.PutFactory(InterfaceMapInjectedStructFactory, "B")

			var injectedB *InterfaceMapInjectedStruct

			Expect(container.CallInjected(func(injected struct {
				B *InterfaceMapInjectedStruct
			}) {
				injectedB = injected.B
			})).Should(Succeed())

			Expect(injectedB.Doers).Should(HaveLen(3))
			Expect(injectedB.Doers).Should(HaveKeyWithValue("A1", &SimpleStruct{"A1"}))
			Expect(injectedB.Doers).Should(HaveKeyWithValue("A2", &SimpleStruct{"A2"}))
			Expect(injectedB.Doers).Should(HaveKeyWithValue("A3", &SimpleStruct{"A3"}))

		})

	})

})

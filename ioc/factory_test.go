package ioc_test

import (
	. "github.com/b-charles/pigs/ioc"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("IOC factory", func() {

	var (
		container *Container
		err       error
	)

	BeforeEach(func() {
		container = NewContainer()
	})

	Describe("with simple injection", func() {

		It("should create a component with dependency", func() {

			err = container.PutNamedFactory(SimpleStructFactory("A"), "A")
			Expect(err).To(Succeed())

			err = container.PutNamedFactory(InjectedStructFactory, "B")
			Expect(err).To(Succeed())

			Expect(container.CallInjected(func(b *InjectedStruct) {
				Expect(b).To(Equal(&InjectedStruct{&SimpleStruct{"A"}}))
			})).Should(Succeed())

		})

		It("should create a component with dependency, another form of factory", func() {

			Expect(container.PutNamedFactory(SimpleStructFactory("SimpleStruct"), "SimpleStruct")).Should(BeNil())

			Expect(container.PutNamedFactory(func(A *SimpleStruct) *InjectedStruct {
				return &InjectedStruct{A}
			}, "InjectedStruct")).Should(BeNil())

			var injectedB *InjectedStruct

			Expect(container.CallInjected(func(injected *InjectedStruct) {
				injectedB = injected
			})).Should(Succeed())

			Expect(injectedB).To(Equal(&InjectedStruct{&SimpleStruct{"SimpleStruct"}}))

		})

		It("should select not nil component", func() {

			Expect(container.PutNamedFactory(func() *SimpleStruct {
				return nil
			}, "NIL", "A")).Should(BeNil())

			Expect(container.PutNamedFactory(SimpleStructFactory("NOT_NIL"), "NOT_NIL", "A")).Should(BeNil())
			Expect(container.PutNamedFactory(InjectedStructFactory, "B")).Should(BeNil())

			var injectedB *InjectedStruct

			Expect(container.CallInjected(func(injected struct {
				B *InjectedStruct
			}) {
				injectedB = injected.B
			})).Should(Succeed())

			Expect(injectedB).To(Equal(&InjectedStruct{&SimpleStruct{"NOT_NIL"}}))

		})

		It("should inject test component if provided", func() {

			Expect(container.PutNamedFactory(SimpleStructFactory("A"), "A")).Should(BeNil())
			Expect(container.TestPutNamedFactory(SimpleStructFactory("TEST"), "TEST", "A")).Should(BeNil())

			Expect(container.PutNamedFactory(InjectedStructFactory, "B")).Should(BeNil())

			var injectedB *InjectedStruct

			Expect(container.CallInjected(func(injected struct {
				B *InjectedStruct
			}) {
				injectedB = injected.B
			})).Should(Succeed())

			Expect(injectedB).To(Equal(&InjectedStruct{&SimpleStruct{"TEST"}}))

		})

		It("should return an error if no component is provided", func() {

			Expect(container.PutNamedFactory(InjectedStructFactory, "B")).Should(BeNil())

			var injectedB *InjectedStruct

			Expect(container.CallInjected(func(injected struct {
				B *InjectedStruct
			}) {
				injectedB = injected.B
				injectedB.doSomething()
			})).Should(HaveOccurred())

		})

		It("should return an error if too many components are provided", func() {

			Expect(container.PutNamedFactory(SimpleStructFactory("A1"), "A1", "A")).Should(BeNil())
			Expect(container.PutNamedFactory(SimpleStructFactory("A2"), "A2", "A")).Should(BeNil())

			Expect(container.PutNamedFactory(InjectedStructFactory, "B")).Should(BeNil())

			var injectedB *InjectedStruct

			Expect(container.CallInjected(func(injected struct {
				B *InjectedStruct
			}) {
				injectedB = injected.B
				injectedB.doSomething()
			})).Should(HaveOccurred())

		})

		It("should inject an interface", func() {

			Expect(container.PutNamedFactory(SimpleStructFactory("A"), "A")).Should(BeNil())

			Expect(container.PutNamedFactory(InterfaceInjectedStructFactory, "B")).Should(BeNil())

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

			Expect(container.PutNamedFactory(SimpleStructFactory("A1"), "A1", "A")).Should(BeNil())
			Expect(container.PutNamedFactory(SimpleStructFactory("A2"), "A2", "A")).Should(BeNil())
			Expect(container.PutNamedFactory(SimpleStructFactory("A3"), "A3", "A")).Should(BeNil())

			Expect(container.PutFactory(SliceInjectedStructFactory, "B")).Should(BeNil())

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

			Expect(container.PutNamedFactory(func() *SimpleStruct {
				return nil
			}, "NIL", "A")).Should(BeNil())

			Expect(container.PutNamedFactory(SimpleStructFactory("A1"), "A1", "A")).Should(BeNil())
			Expect(container.PutNamedFactory(SimpleStructFactory("A2"), "A2", "A")).Should(BeNil())
			Expect(container.PutNamedFactory(SimpleStructFactory("A3"), "A3", "A")).Should(BeNil())

			Expect(container.PutNamedFactory(SliceInjectedStructFactory, "B")).Should(BeNil())

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

			Expect(container.PutNamedFactory(SimpleStructFactory("A1"), "A1", "A")).Should(BeNil())
			Expect(container.PutNamedFactory(SimpleStructFactory("A2"), "A2", "A")).Should(BeNil())
			Expect(container.PutNamedFactory(SimpleStructFactory("A3"), "A3", "A")).Should(BeNil())

			Expect(container.PutNamedFactory(InterfaceSliceInjectedStructFactory, "B")).Should(BeNil())

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

			Expect(container.PutNamedFactory(SimpleStructFactory("A1"), "A1", "A")).Should(BeNil())
			Expect(container.PutNamedFactory(SimpleStructFactory("A2"), "A2", "A")).Should(BeNil())
			Expect(container.PutNamedFactory(SimpleStructFactory("A3"), "A3", "A")).Should(BeNil())

			Expect(container.PutNamedFactory(MapInjectedStructFactory, "B")).Should(BeNil())

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

			Expect(container.PutFactory(func() *SimpleStruct {
				return nil
			}, "NIL", "A")).Should(BeNil())

			Expect(container.PutNamedFactory(SimpleStructFactory("A1"), "A1", "A")).Should(BeNil())
			Expect(container.PutNamedFactory(SimpleStructFactory("A2"), "A2", "A")).Should(BeNil())
			Expect(container.PutNamedFactory(SimpleStructFactory("A3"), "A3", "A")).Should(BeNil())

			Expect(container.PutFactory(MapInjectedStructFactory, "B")).Should(BeNil())

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

			Expect(container.PutNamedFactory(SimpleStructFactory("A1"), "A1", "A")).Should(BeNil())
			Expect(container.PutNamedFactory(SimpleStructFactory("A2"), "A2", "A")).Should(BeNil())
			Expect(container.PutNamedFactory(SimpleStructFactory("A3"), "A3", "A")).Should(BeNil())

			Expect(container.PutFactory(InterfaceMapInjectedStructFactory, "B")).Should(BeNil())

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

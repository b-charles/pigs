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

		It("should create a component without dependency", func() {

			Expect(container.Put(&SimpleStruct{"A"})).Should(BeNil())

			var injectedA *SimpleStruct

			Expect(container.CallInjected(func(injected struct {
				A *SimpleStruct
			}) {
				injectedA = injected.A
			})).Should(Succeed())

			Expect(injectedA).To(Equal(&SimpleStruct{"A"}))

		})

		It("should create a component with dependency", func() {

			Expect(container.PutNamedFactory(SimpleStructFactory("A"), "A")).Should(BeNil())
			Expect(container.Put(&InjectedObject{})).Should(BeNil())

			var injectedB *InjectedObject

			Expect(container.CallInjected(func(injected struct {
				B *InjectedObject
			}) {
				injectedB = injected.B
			})).Should(Succeed())

			Expect(injectedB).To(Equal(&InjectedObject{&SimpleStruct{"A"}}))

		})

		It("should inject test component if provided", func() {

			Expect(container.PutNamedFactory(SimpleStructFactory("A"), "A")).Should(BeNil())
			Expect(container.TestPutNamedFactory(SimpleStructFactory("TEST"), "TEST", "A")).Should(BeNil())

			Expect(container.Put(&InjectedObject{})).Should(BeNil())

			var injectedB *InjectedObject

			Expect(container.CallInjected(func(injected struct {
				B *InjectedObject
			}) {
				injectedB = injected.B
			})).Should(Succeed())

			Expect(injectedB).To(Equal(&InjectedObject{&SimpleStruct{"TEST"}}))

		})

		It("should return an error if no component is provided", func() {

			Expect(container.Put(&InjectedObject{})).Should(BeNil())

			var injectedB *InjectedObject

			Expect(container.CallInjected(func(injected struct {
				B *InjectedObject
			}) {
				injectedB = injected.B
				injectedB.doSomething()
			})).Should(HaveOccurred())

		})

		It("should return an error if too many components are provided", func() {

			Expect(container.PutNamedFactory(SimpleStructFactory("A1"), "A1", "A")).Should(BeNil())
			Expect(container.PutNamedFactory(SimpleStructFactory("A2"), "A2", "A")).Should(BeNil())

			Expect(container.Put(&InjectedObject{})).Should(BeNil())

			var injectedB *InjectedObject

			Expect(container.CallInjected(func(injected struct {
				B *InjectedObject
			}) {
				injectedB = injected.B
				injectedB.doSomething()
			})).Should(HaveOccurred())

		})

		It("should inject an interface by field name", func() {

			Expect(container.PutNamedFactory(SimpleStructFactory("A"), "A")).Should(BeNil())
			Expect(container.Put(&InterfaceInjectedObject{})).Should(BeNil())

			var injectedB *InterfaceInjectedObject

			Expect(container.CallInjected(func(injected struct {
				B *InterfaceInjectedObject
			}) {
				injectedB = injected.B
			})).Should(Succeed())

			Expect(injectedB).To(Equal(&InterfaceInjectedObject{&SimpleStruct{"A"}}))

		})

		It("should inject an interface by tag name", func() {

			Expect(container.PutNamedFactory(SimpleStructFactory("A"), "A")).Should(BeNil())
			Expect(container.Put(&NamedInterfaceInjectedObject{})).Should(BeNil())

			var injectedB *NamedInterfaceInjectedObject

			Expect(container.CallInjected(func(injected struct {
				B *NamedInterfaceInjectedObject
			}) {
				injectedB = injected.B
			})).Should(Succeed())

			Expect(injectedB).To(Equal(&NamedInterfaceInjectedObject{&SimpleStruct{"A"}}))

		})

		It("should inject the component in itself", func() {

			Expect(container.Put(&Looping{})).Should(BeNil())

			var injectedLooping *Looping

			Expect(container.CallInjected(func(injected struct {
				Looping *Looping
			}) {
				injectedLooping = injected.Looping
			})).Should(Succeed())

			Expect(injectedLooping.Looping).To(Equal(injectedLooping))

		})

	})

	Describe("with slice injection", func() {

		It("should create a component with a slice of dependencies", func() {

			Expect(container.PutNamedFactory(SimpleStructFactory("A1"), "A1", "A")).Should(BeNil())
			Expect(container.PutNamedFactory(SimpleStructFactory("A2"), "A2", "A")).Should(BeNil())
			Expect(container.PutNamedFactory(SimpleStructFactory("A3"), "A3", "A")).Should(BeNil())

			Expect(container.Put(&SliceInjectedObject{})).Should(BeNil())

			var injectedB *SliceInjectedObject

			Expect(container.CallInjected(func(injected struct {
				B *SliceInjectedObject
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

			Expect(container.Put(&SliceInjectedObject{})).Should(BeNil())

			var injectedB *SliceInjectedObject

			Expect(container.CallInjected(func(injected struct {
				B *SliceInjectedObject
			}) {
				injectedB = injected.B
			})).Should(Succeed())

			Expect(injectedB.SimpleStructs).Should(ConsistOf(
				&SimpleStruct{"A1"},
				&SimpleStruct{"A2"},
				&SimpleStruct{"A3"}))

		})

		It("should inject a slice of interface", func() {

			Expect(container.PutNamedFactory(SimpleStructFactory("A1"), "A1", "Doers")).Should(BeNil())
			Expect(container.PutNamedFactory(SimpleStructFactory("A2"), "A2", "Doers")).Should(BeNil())
			Expect(container.PutNamedFactory(SimpleStructFactory("A3"), "A3", "Doers")).Should(BeNil())

			Expect(container.Put(&InterfaceSliceInjectedObject{})).Should(BeNil())

			var injectedB *InterfaceSliceInjectedObject

			Expect(container.CallInjected(func(injected struct {
				B *InterfaceSliceInjectedObject
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

			Expect(container.Put(&MapInjectedObject{})).Should(BeNil())

			var injectedB *MapInjectedObject

			Expect(container.CallInjected(func(injected struct {
				B *MapInjectedObject
			}) {
				injectedB = injected.B
			})).Should(Succeed())

			Expect(injectedB.SimpleStructs).Should(HaveLen(3))
			Expect(injectedB.SimpleStructs).Should(HaveKeyWithValue("A1", &SimpleStruct{"A1"}))
			Expect(injectedB.SimpleStructs).Should(HaveKeyWithValue("A2", &SimpleStruct{"A2"}))
			Expect(injectedB.SimpleStructs).Should(HaveKeyWithValue("A3", &SimpleStruct{"A3"}))

		})

		It("should inject a map without nil components", func() {

			Expect(container.PutNamedFactory(func() *SimpleStruct {
				return nil
			}, "NIL", "A")).Should(BeNil())

			Expect(container.PutNamedFactory(SimpleStructFactory("A1"), "A1", "A")).Should(BeNil())
			Expect(container.PutNamedFactory(SimpleStructFactory("A2"), "A2", "A")).Should(BeNil())
			Expect(container.PutNamedFactory(SimpleStructFactory("A3"), "A3", "A")).Should(BeNil())

			Expect(container.Put(&MapInjectedObject{})).Should(BeNil())

			var injectedB *MapInjectedObject

			Expect(container.CallInjected(func(injected struct {
				B *MapInjectedObject
			}) {
				injectedB = injected.B
			})).Should(Succeed())

			Expect(injectedB.SimpleStructs).Should(HaveLen(3))
			Expect(injectedB.SimpleStructs).Should(HaveKeyWithValue("A1", &SimpleStruct{"A1"}))
			Expect(injectedB.SimpleStructs).Should(HaveKeyWithValue("A2", &SimpleStruct{"A2"}))
			Expect(injectedB.SimpleStructs).Should(HaveKeyWithValue("A3", &SimpleStruct{"A3"}))

		})

		It("should inject a map of interface", func() {

			Expect(container.PutNamedFactory(SimpleStructFactory("A1"), "A1", "Doers")).Should(BeNil())
			Expect(container.PutNamedFactory(SimpleStructFactory("A2"), "A2", "Doers")).Should(BeNil())
			Expect(container.PutNamedFactory(SimpleStructFactory("A3"), "A3", "Doers")).Should(BeNil())

			Expect(container.Put(&InterfaceMapInjectedObject{})).Should(BeNil())

			var injectedB *InterfaceMapInjectedObject

			Expect(container.CallInjected(func(injected struct {
				B *InterfaceMapInjectedObject
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

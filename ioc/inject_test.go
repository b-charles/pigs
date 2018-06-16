package ioc_test

import (
	. "github.com/l3eegbee/pigs/ioc"
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

			container.Put(&SimpleStruct{"A"}, "A")

			var injectedA *SimpleStruct

			Expect(container.CallInjected(func(injected struct {
				A *SimpleStruct
			}) {
				injectedA = injected.A
			})).Should(Succeed())

			Expect(injectedA).To(Equal(&SimpleStruct{"A"}))

		})

		It("should create a component with dependency", func() {

			container.PutFactory(SimpleStructFactory("A"), "A")
			container.Put(&InjectedObject{}, "B")

			var injectedB *InjectedObject

			Expect(container.CallInjected(func(injected struct {
				B *InjectedObject
			}) {
				injectedB = injected.B
			})).Should(Succeed())

			Expect(injectedB).To(Equal(&InjectedObject{&SimpleStruct{"A"}}))

		})

		It("should select not nil component", func() {

			container.PutFactory(func() *SimpleStruct {
				return nil
			}, "NIL", "A")

			container.PutFactory(SimpleStructFactory("A"), "A")
			container.Put(&InjectedObject{}, "B")

			var injectedB *InjectedObject

			Expect(container.CallInjected(func(injected struct {
				B *InjectedObject
			}) {
				injectedB = injected.B
			})).Should(Succeed())

			Expect(injectedB).To(Equal(&InjectedObject{&SimpleStruct{"A"}}))

		})

		It("should inject test component if provided", func() {

			container.PutFactory(SimpleStructFactory("A"), "A")
			container.TestPutFactory(SimpleStructFactory("TEST"), "TEST", "A")

			container.Put(&InjectedObject{}, "B")

			var injectedB *InjectedObject

			Expect(container.CallInjected(func(injected struct {
				B *InjectedObject
			}) {
				injectedB = injected.B
			})).Should(Succeed())

			Expect(injectedB).To(Equal(&InjectedObject{&SimpleStruct{"TEST"}}))

		})

		It("should return an error if no component is provided", func() {

			container.Put(&InjectedObject{}, "B")

			var injectedB *InjectedObject

			Expect(container.CallInjected(func(injected struct {
				B *InjectedObject
			}) {
				injectedB = injected.B
			})).Should(HaveOccurred())

		})

		It("should return an error if too many components are provided", func() {

			container.PutFactory(SimpleStructFactory("A1"), "A1", "A")
			container.PutFactory(SimpleStructFactory("A2"), "A2", "A")

			container.Put(&InjectedObject{}, "B")

			var injectedB *InjectedObject

			Expect(container.CallInjected(func(injected struct {
				B *InjectedObject
			}) {
				injectedB = injected.B
			})).Should(HaveOccurred())

		})

		It("should inject an interface by field name", func() {

			container.PutFactory(SimpleStructFactory("A"), "A")
			container.Put(&InterfaceInjectedObject{}, "B")

			var injectedB *InterfaceInjectedObject

			Expect(container.CallInjected(func(injected struct {
				B *InterfaceInjectedObject
			}) {
				injectedB = injected.B
			})).Should(Succeed())

			Expect(injectedB).To(Equal(&InterfaceInjectedObject{&SimpleStruct{"A"}}))

		})

		It("should inject an interface by tag name", func() {

			container.PutFactory(SimpleStructFactory("A"), "A")
			container.Put(&NamedInterfaceInjectedObject{}, "B")

			var injectedB *NamedInterfaceInjectedObject

			Expect(container.CallInjected(func(injected struct {
				B *NamedInterfaceInjectedObject
			}) {
				injectedB = injected.B
			})).Should(Succeed())

			Expect(injectedB).To(Equal(&NamedInterfaceInjectedObject{&SimpleStruct{"A"}}))

		})

		It("should inject the component in itself", func() {

			container.Put(&Looping{}, "Looping")

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

			container.PutFactory(SimpleStructFactory("A1"), "A1", "A")
			container.PutFactory(SimpleStructFactory("A2"), "A2", "A")
			container.PutFactory(SimpleStructFactory("A3"), "A3", "A")

			container.Put(&SliceInjectedObject{}, "B")

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

			container.PutFactory(func() *SimpleStruct {
				return nil
			}, "NIL", "A")

			container.PutFactory(SimpleStructFactory("A1"), "A1", "A")
			container.PutFactory(SimpleStructFactory("A2"), "A2", "A")
			container.PutFactory(SimpleStructFactory("A3"), "A3", "A")

			container.Put(&SliceInjectedObject{}, "B")

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

			container.PutFactory(SimpleStructFactory("A1"), "A1", "Doers")
			container.PutFactory(SimpleStructFactory("A2"), "A2", "Doers")
			container.PutFactory(SimpleStructFactory("A3"), "A3", "Doers")

			container.Put(&InterfaceSliceInjectedObject{}, "B")

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

			container.PutFactory(SimpleStructFactory("A1"), "A1", "A")
			container.PutFactory(SimpleStructFactory("A2"), "A2", "A")
			container.PutFactory(SimpleStructFactory("A3"), "A3", "A")

			container.Put(&MapInjectedObject{}, "B")

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

			container.PutFactory(func() *SimpleStruct {
				return nil
			}, "NIL", "A")

			container.PutFactory(SimpleStructFactory("A1"), "A1", "A")
			container.PutFactory(SimpleStructFactory("A2"), "A2", "A")
			container.PutFactory(SimpleStructFactory("A3"), "A3", "A")

			container.Put(&MapInjectedObject{}, "B")

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

			container.PutFactory(SimpleStructFactory("A1"), "A1", "Doers")
			container.PutFactory(SimpleStructFactory("A2"), "A2", "Doers")
			container.PutFactory(SimpleStructFactory("A3"), "A3", "Doers")

			container.Put(&InterfaceMapInjectedObject{}, "B")

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

package ioc_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/l3eegbee/pigs/ioc"
)

// Simple interface

type SomethingDoer interface {
	doSomething()
}

// Class SimpleStruct

type SimpleStruct struct {
	Tag string
}

func (self *SimpleStruct) doSomething() {}

func SimpleStructFactory(tag string) func() *SimpleStruct {
	return func() *SimpleStruct { return &SimpleStruct{tag} }
}

// Class InjectedStruct

type InjectedStruct struct {
	SimpleStruct *SimpleStruct
}

func InjectedStructFactory(injected struct {
	SimpleStruct *SimpleStruct `inject:"A"`
}) *InjectedStruct {
	return &InjectedStruct{injected.SimpleStruct}
}

// Class SliceInjectedStruct

type SliceInjectedStruct struct {
	SimpleStructs []*SimpleStruct
}

func SliceInjectedStructFactory(injected struct {
	SimpleStructs []*SimpleStruct `inject:"A"`
}) *SliceInjectedStruct {
	return &SliceInjectedStruct{injected.SimpleStructs}
}

// Class MapInjectedStruct

type MapInjectedStruct struct {
	SimpleStructs map[string]*SimpleStruct
}

func MapInjectedStructFactory(injected struct {
	SimpleStructs map[string]*SimpleStruct `inject:"A"`
}) *MapInjectedStruct {
	return &MapInjectedStruct{injected.SimpleStructs}
}

// Class InterfaceInjectedStruct

type InterfaceInjectedStruct struct {
	Doer SomethingDoer
}

func InterfaceInjectedStructFactory(injected struct {
	Doer SomethingDoer `inject:"A"`
}) *InterfaceInjectedStruct {
	return &InterfaceInjectedStruct{injected.Doer}
}

type InterfaceInjectedObject struct {
	A SomethingDoer `inject:""`
}

type NamedInterfaceInjectedObject struct {
	Doer SomethingDoer `inject:"A"`
}

// Class InterfaceSliceInjectedStruct

type InterfaceSliceInjectedStruct struct {
	Doers []SomethingDoer
}

func InterfaceSliceInjectedStructFactory(injected struct {
	Doers []SomethingDoer `inject:"A"`
}) *InterfaceSliceInjectedStruct {
	return &InterfaceSliceInjectedStruct{injected.Doers}
}

type InterfaceSliceInjectedObject struct {
	Doers []SomethingDoer `inject:""`
}

// Class InterfaceMapInjectedStruct

type InterfaceMapInjectedStruct struct {
	Doers map[string]SomethingDoer
}

func InterfaceMapInjectedStructFactory(injected struct {
	Doers map[string]SomethingDoer `inject:"A"`
}) *InterfaceMapInjectedStruct {
	return &InterfaceMapInjectedStruct{injected.Doers}
}

type InterfaceMapInjectedObject struct {
	Doers map[string]SomethingDoer `inject:""`
}

// Class Looping (inject itself)

type Looping struct {
	Looping *Looping `inject:""`
}

// Ordered classes

type Ordered interface {
	Name() string
}

type OrderRegister struct {
	PostInitOrder []string
	CloseOrder    []string
}

func NewOrderRegster() *OrderRegister {
	return &OrderRegister{make([]string, 0), make([]string, 0)}
}

func (self *OrderRegister) registerPostInit(obj Ordered) {
	self.PostInitOrder = append(self.PostInitOrder, obj.Name())
}
func (self *OrderRegister) registerClose(obj Ordered) {
	self.CloseOrder = append(self.CloseOrder, obj.Name())
}

type First struct {
	OrderRegister *OrderRegister `inject:""`
}

func (self *First) Name() string { return "FIRST" }
func (self *First) PostInit()    { self.OrderRegister.registerPostInit(self) }
func (self *First) Close()       { self.OrderRegister.registerClose(self) }

type Second struct {
	OrderRegister *OrderRegister `inject:""`
	Injected      *First         `inject:"First"`
}

func (self *Second) Name() string { return "SECOND" }
func (self *Second) PostInit()    { self.OrderRegister.registerPostInit(self) }
func (self *Second) Close()       { self.OrderRegister.registerClose(self) }

type Third struct {
	OrderRegister *OrderRegister `inject:""`
	Injected      *Second        `inject:"Second"`
}

func (self *Third) Name() string { return "THIRD" }
func (self *Third) PostInit()    { self.OrderRegister.registerPostInit(self) }
func (self *Third) Close()       { self.OrderRegister.registerClose(self) }

// Not pointers injections

type MapProvider interface {
	GetMap() map[string]string
}

type MyMap map[string]string

func NewMap(key, value string) MyMap {
	m := make(map[string]string, 0)
	m[key] = value
	return m
}
func (self MyMap) GetMap() map[string]string { return self }

func CreateWrapMap(injected struct {
	Map MyMap `inject:"MyMap"`
}) MyMap {
	return injected.Map
}

type MySuperMap map[string]string

func CreateSuperMap(injected struct {
	Maps []MapProvider `inject:"Maps"`
}) MySuperMap {
	super := make(map[string]string, 0)
	for _, m := range injected.Maps {
		for k, v := range m.GetMap() {
			super[k] = v
		}
	}
	return super
}

type MyHyperMap map[string]string

func CreateHyperMap(injected struct {
	Maps map[string]MapProvider `inject:"Maps"`
}) MyHyperMap {
	hyper := make(map[string]string, 0)
	for name, m := range injected.Maps {
		for k, v := range m.GetMap() {
			hyper[name+"."+k] = v
		}
	}
	return hyper
}

// Tests

var _ = Describe("Container", func() {

	var (
		container *Container
	)

	BeforeEach(func() {
		container = NewContainer()
	})

	Describe("Get itself", func() {

		It("should be able to extract itself", func() {
			Expect(container.GetComponent("ApplicationContainer")).To(Equal(container))
		})

	})

	Describe("Factory injection", func() {

		It("should create a component without dependency", func() {

			container.PutFactory(SimpleStructFactory("A"), "A")

			a := container.GetComponent("A").(*SimpleStruct)

			Expect(a).Should(Equal(&SimpleStruct{"A"}))

		})

		Context("For a simple depency", func() {

			Context("For struct dependency", func() {

				It("should create a component with a dependency", func() {

					container.PutFactory(SimpleStructFactory("A"), "A")
					container.PutFactory(InjectedStructFactory, "B")

					b := container.GetComponent("B").(*InjectedStruct)

					Expect(b.SimpleStruct).Should(Equal(&SimpleStruct{"A"}))

				})

				It("should select not nil component", func() {

					container.PutFactory(func() *SimpleStruct {
						return nil
					}, "NIL", "A")

					container.PutFactory(SimpleStructFactory("NOT_NIL"), "NOT_NIL", "A")
					container.PutFactory(InjectedStructFactory, "B")

					b := container.GetComponent("B").(*InjectedStruct)

					Expect(b.SimpleStruct).Should(Equal(&SimpleStruct{"NOT_NIL"}))

				})

				It("should inject test component if provided", func() {

					container.PutFactory(SimpleStructFactory("A"), "A")
					container.PutFactory(InjectedStructFactory, "B")

					container.TestPutFactory(SimpleStructFactory("TEST"), "TEST", "A")

					b := container.GetComponent("B").(*InjectedStruct)

					Expect(b.SimpleStruct).Should(Equal(&SimpleStruct{"TEST"}))

				})

				It("should panic if no component is provided", func() {

					container.PutFactory(InjectedStructFactory, "B")

					Expect(func() {
						_ = container.GetComponent("B")
					}).Should(Panic())

				})

				It("should panic if many components are provided", func() {

					container.PutFactory(SimpleStructFactory("A1"), "A1", "A")
					container.PutFactory(SimpleStructFactory("A2"), "A2", "A")
					container.PutFactory(InjectedStructFactory, "B")

					Expect(func() {
						_ = container.GetComponent("B")
					}).Should(Panic())

				})

			})

			Context("For interface dependency", func() {

				It("should create a component with a dependency", func() {

					container.PutFactory(SimpleStructFactory("A"), "A")
					container.PutFactory(InterfaceInjectedStructFactory, "B")

					b := container.GetComponent("B").(*InterfaceInjectedStruct)

					Expect(b.Doer).Should(Equal(&SimpleStruct{"A"}))

				})

				It("should select not nil component", func() {

					container.PutFactory(func() *SimpleStruct {
						return nil
					}, "NIL", "A")

					container.PutFactory(SimpleStructFactory("NOT_NIL"), "NOT_NIL", "A")
					container.PutFactory(InterfaceInjectedStructFactory, "B")

					b := container.GetComponent("B").(*InterfaceInjectedStruct)

					Expect(b.Doer).Should(Equal(&SimpleStruct{"NOT_NIL"}))

				})

				It("should inject test component if provided", func() {

					container.PutFactory(SimpleStructFactory("A"), "A")
					container.PutFactory(InterfaceInjectedStructFactory, "B")

					container.TestPutFactory(SimpleStructFactory("TEST"), "TEST", "A")

					b := container.GetComponent("B").(*InterfaceInjectedStruct)

					Expect(b.Doer).Should(Equal(&SimpleStruct{"TEST"}))

				})

				It("should panic if no component is provided", func() {

					container.PutFactory(InterfaceInjectedStructFactory, "B")

					Expect(func() {
						_ = container.GetComponent("B")
					}).Should(Panic())

				})

				It("should panic if many components are provided", func() {

					container.PutFactory(SimpleStructFactory("A1"), "A1", "A")
					container.PutFactory(SimpleStructFactory("A2"), "A2", "A")
					container.PutFactory(InterfaceInjectedStructFactory, "B")

					Expect(func() {
						_ = container.GetComponent("B")
					}).Should(Panic())

				})

			})

			It("should restore 'core' configuration after a TestClear", func() {

				container.PutFactory(SimpleStructFactory("A"), "A")
				container.PutFactory(InjectedStructFactory, "B")

				// test

				container.TestPutFactory(SimpleStructFactory("TEST"), "TEST", "A")

				b := container.GetComponent("B").(*InjectedStruct)

				Expect(b.SimpleStruct.Tag).To(Equal("TEST"))

				// restore

				container.ClearTests()

				// core

				b = container.GetComponent("B").(*InjectedStruct)

				Expect(b.SimpleStruct.Tag).To(Equal("A"))

			})

		})

		Context("For a depency of type slice", func() {

			Context("For struct dependency", func() {

				It("should create a component with auto-discovery", func() {

					container.PutFactory(SimpleStructFactory("A1"), "A1", "A")
					container.PutFactory(SimpleStructFactory("A2"), "A2", "A")
					container.PutFactory(SimpleStructFactory("A3"), "A3", "A")

					container.PutFactory(SliceInjectedStructFactory, "B")

					b := container.GetComponent("B").(*SliceInjectedStruct)

					Expect(b.SimpleStructs).Should(ConsistOf(
						&SimpleStruct{"A1"},
						&SimpleStruct{"A2"},
						&SimpleStruct{"A3"}))

				})

				It("should select not nil components", func() {

					container.PutFactory(func() *SimpleStruct {
						return nil
					}, "NIL", "A")

					container.PutFactory(SimpleStructFactory("A1"), "A1", "A")
					container.PutFactory(SimpleStructFactory("A2"), "A2", "A")
					container.PutFactory(SimpleStructFactory("A3"), "A3", "A")

					container.PutFactory(SliceInjectedStructFactory, "B")

					b := container.GetComponent("B").(*SliceInjectedStruct)

					Expect(b.SimpleStructs).Should(ConsistOf(
						&SimpleStruct{"A1"},
						&SimpleStruct{"A2"},
						&SimpleStruct{"A3"}))

				})

				It("should inject test component if provided", func() {

					container.PutFactory(SimpleStructFactory("A1"), "A1", "A")
					container.PutFactory(SimpleStructFactory("A2"), "A2", "A")
					container.PutFactory(SimpleStructFactory("A3"), "A3", "A")

					container.TestPutFactory(SimpleStructFactory("TEST"), "TEST", "A")

					container.PutFactory(SliceInjectedStructFactory, "B")

					b := container.GetComponent("B").(*SliceInjectedStruct)

					Expect(b.SimpleStructs).Should(ConsistOf(&SimpleStruct{"TEST"}))

				})

			})

			Context("For interface dependency", func() {

				It("should create a component with auto-discovery", func() {

					container.PutFactory(SimpleStructFactory("A1"), "A1", "A")
					container.PutFactory(SimpleStructFactory("A2"), "A2", "A")
					container.PutFactory(SimpleStructFactory("A3"), "A3", "A")

					container.PutFactory(InterfaceSliceInjectedStructFactory, "B")

					b := container.GetComponent("B").(*InterfaceSliceInjectedStruct)

					Expect(b.Doers).Should(ConsistOf(
						&SimpleStruct{"A1"},
						&SimpleStruct{"A2"},
						&SimpleStruct{"A3"}))

				})

				It("should select not nil components", func() {

					container.PutFactory(func() *SimpleStruct {
						return nil
					}, "NIL", "A")

					container.PutFactory(SimpleStructFactory("A1"), "A1", "A")
					container.PutFactory(SimpleStructFactory("A2"), "A2", "A")
					container.PutFactory(SimpleStructFactory("A3"), "A3", "A")

					container.PutFactory(InterfaceSliceInjectedStructFactory, "B")

					b := container.GetComponent("B").(*InterfaceSliceInjectedStruct)

					Expect(b.Doers).Should(ConsistOf(
						&SimpleStruct{"A1"},
						&SimpleStruct{"A2"},
						&SimpleStruct{"A3"}))

				})

				It("should inject test component if provided", func() {

					container.PutFactory(SimpleStructFactory("A1"), "A1", "A")
					container.PutFactory(SimpleStructFactory("A2"), "A2", "A")
					container.PutFactory(SimpleStructFactory("A3"), "A3", "A")

					container.TestPutFactory(SimpleStructFactory("TEST"), "TEST", "A")

					container.PutFactory(InterfaceSliceInjectedStructFactory, "B")

					b := container.GetComponent("B").(*InterfaceSliceInjectedStruct)

					Expect(b.Doers).Should(ConsistOf(&SimpleStruct{"TEST"}))

				})

			})

		})

		Context("For a depency of type map", func() {

			Context("For struct dependency", func() {

				It("should create a component with auto-discovery", func() {

					container.PutFactory(SimpleStructFactory("A1"), "A1", "A")
					container.PutFactory(SimpleStructFactory("A2"), "A2", "A")
					container.PutFactory(SimpleStructFactory("A3"), "A3", "A")

					container.PutFactory(MapInjectedStructFactory, "B")

					b := container.GetComponent("B").(*MapInjectedStruct)

					Expect(b.SimpleStructs).Should(HaveLen(3))
					Expect(b.SimpleStructs).Should(HaveKeyWithValue("A1", &SimpleStruct{"A1"}))
					Expect(b.SimpleStructs).Should(HaveKeyWithValue("A2", &SimpleStruct{"A2"}))
					Expect(b.SimpleStructs).Should(HaveKeyWithValue("A3", &SimpleStruct{"A3"}))

				})

				It("should select not nil components", func() {

					container.PutFactory(func() *SimpleStruct {
						return nil
					}, "NIL", "A")

					container.PutFactory(SimpleStructFactory("A1"), "A1", "A")
					container.PutFactory(SimpleStructFactory("A2"), "A2", "A")
					container.PutFactory(SimpleStructFactory("A3"), "A3", "A")

					container.PutFactory(MapInjectedStructFactory, "B")

					b := container.GetComponent("B").(*MapInjectedStruct)

					Expect(b.SimpleStructs).Should(HaveLen(3))
					Expect(b.SimpleStructs).Should(HaveKeyWithValue("A1", &SimpleStruct{"A1"}))
					Expect(b.SimpleStructs).Should(HaveKeyWithValue("A2", &SimpleStruct{"A2"}))
					Expect(b.SimpleStructs).Should(HaveKeyWithValue("A3", &SimpleStruct{"A3"}))

				})

				It("should inject test component if provided", func() {

					container.PutFactory(SimpleStructFactory("A1"), "A1", "A")
					container.PutFactory(SimpleStructFactory("A2"), "A2", "A")
					container.PutFactory(SimpleStructFactory("A3"), "A3", "A")

					container.TestPutFactory(SimpleStructFactory("TEST"), "TEST", "A")

					container.PutFactory(MapInjectedStructFactory, "B")

					b := container.GetComponent("B").(*MapInjectedStruct)

					Expect(b.SimpleStructs).Should(HaveLen(1))
					Expect(b.SimpleStructs).Should(HaveKeyWithValue("TEST", &SimpleStruct{"TEST"}))

				})

			})

			Context("For interface dependency", func() {

				It("should create a component with auto-discovery", func() {

					container.PutFactory(SimpleStructFactory("A1"), "A1", "A")
					container.PutFactory(SimpleStructFactory("A2"), "A2", "A")
					container.PutFactory(SimpleStructFactory("A3"), "A3", "A")

					container.PutFactory(InterfaceMapInjectedStructFactory, "B")

					b := container.GetComponent("B").(*InterfaceMapInjectedStruct)

					Expect(b.Doers).Should(HaveLen(3))
					Expect(b.Doers).Should(HaveKeyWithValue("A1", &SimpleStruct{"A1"}))
					Expect(b.Doers).Should(HaveKeyWithValue("A2", &SimpleStruct{"A2"}))
					Expect(b.Doers).Should(HaveKeyWithValue("A3", &SimpleStruct{"A3"}))

				})

				It("should select not nil components", func() {

					container.PutFactory(func() *SimpleStruct {
						return nil
					}, "NIL", "A")

					container.PutFactory(SimpleStructFactory("A1"), "A1", "A")
					container.PutFactory(SimpleStructFactory("A2"), "A2", "A")
					container.PutFactory(SimpleStructFactory("A3"), "A3", "A")

					container.PutFactory(InterfaceMapInjectedStructFactory, "B")

					b := container.GetComponent("B").(*InterfaceMapInjectedStruct)

					Expect(b.Doers).Should(HaveLen(3))
					Expect(b.Doers).Should(HaveKeyWithValue("A1", &SimpleStruct{"A1"}))
					Expect(b.Doers).Should(HaveKeyWithValue("A2", &SimpleStruct{"A2"}))
					Expect(b.Doers).Should(HaveKeyWithValue("A3", &SimpleStruct{"A3"}))

				})

				It("should inject test component if provided", func() {

					container.PutFactory(SimpleStructFactory("A1"), "A1", "A")
					container.PutFactory(SimpleStructFactory("A2"), "A2", "A")
					container.PutFactory(SimpleStructFactory("A3"), "A3", "A")

					container.TestPutFactory(SimpleStructFactory("TEST"), "TEST", "A")

					container.PutFactory(InterfaceMapInjectedStructFactory, "B")

					b := container.GetComponent("B").(*InterfaceMapInjectedStruct)

					Expect(b.Doers).Should(HaveLen(1))
					Expect(b.Doers).Should(HaveKeyWithValue("TEST", &SimpleStruct{"TEST"}))

				})

			})

		})

	})

	Describe("Instance injection", func() {

		It("should create a component without dependency", func() {

			container.Put(&SimpleStruct{"A"}, "A")

			a := container.GetComponent("A").(*SimpleStruct)

			Expect(a).Should(Equal(&SimpleStruct{"A"}))

		})

		It("should create a component with a simple dependency", func() {

			container.PutFactory(SimpleStructFactory("A"), "A")
			container.Put(&InterfaceInjectedObject{}, "B")

			b := container.GetComponent("B").(*InterfaceInjectedObject)

			Expect(b.A).Should(Equal(&SimpleStruct{"A"}))

		})

		It("should create a component with a simple dependency defined by name", func() {

			container.PutFactory(SimpleStructFactory("A"), "A")
			container.Put(&NamedInterfaceInjectedObject{}, "B")

			b := container.GetComponent("B").(*NamedInterfaceInjectedObject)

			Expect(b.Doer).Should(Equal(&SimpleStruct{"A"}))

		})

		It("should create a component with a slice of dependencies", func() {

			container.PutFactory(SimpleStructFactory("A1"), "A1", "Doers")
			container.PutFactory(SimpleStructFactory("A2"), "A2", "Doers")
			container.PutFactory(SimpleStructFactory("A3"), "A3", "Doers")

			container.Put(&InterfaceSliceInjectedObject{}, "B")

			b := container.GetComponent("B").(*InterfaceSliceInjectedObject)

			Expect(b.Doers).Should(ConsistOf(
				&SimpleStruct{"A1"},
				&SimpleStruct{"A2"},
				&SimpleStruct{"A3"}))

		})

		It("should create a component with a map of dependencies", func() {

			container.PutFactory(SimpleStructFactory("A1"), "A1", "Doers")
			container.PutFactory(SimpleStructFactory("A2"), "A2", "Doers")
			container.PutFactory(SimpleStructFactory("A3"), "A3", "Doers")

			container.Put(&InterfaceMapInjectedObject{}, "B")

			b := container.GetComponent("B").(*InterfaceMapInjectedObject)

			Expect(b.Doers).Should(HaveLen(3))
			Expect(b.Doers).Should(HaveKeyWithValue("A1", &SimpleStruct{"A1"}))
			Expect(b.Doers).Should(HaveKeyWithValue("A2", &SimpleStruct{"A2"}))
			Expect(b.Doers).Should(HaveKeyWithValue("A3", &SimpleStruct{"A3"}))

		})

	})

	Describe("Loop dependencies", func() {

		It("should inject a component in itself", func() {

			container.Put(&Looping{}, "Looping")

			looping := container.GetComponent("Looping").(*Looping)

			Expect(looping.Looping).To(Equal(looping))

		})

	})

	Describe("PostInit and Close", func() {

		It("should call PostInit in the correct order", func() {

			container.PutFactory(NewOrderRegster, "OrderRegister")
			container.Put(&First{}, "First")
			container.Put(&Second{}, "Second")
			container.Put(&Third{}, "Third")

			third := container.GetComponent("Third").(*Third)

			Expect(third.OrderRegister.PostInitOrder).Should(Equal([]string{"FIRST", "SECOND", "THIRD"}))

		})

		It("should call Close in the correct order", func() {

			container.PutFactory(NewOrderRegster, "OrderRegister")
			container.Put(&First{}, "First")
			container.Put(&Second{}, "Second")
			container.Put(&Third{}, "Third")

			third := container.GetComponent("Third").(*Third)

			container.Close()

			Expect(third.OrderRegister.CloseOrder).Should(Equal([]string{"FIRST", "SECOND", "THIRD"}))

		})

	})

	Describe("Freaks show", func() {

		It("should handle not struct component", func() {

			container.Put(NewMap("Hello", "World"), "MyMap")

			myMap := container.GetComponent("MyMap").(MyMap)

			Expect(myMap).Should(HaveLen(1))
			Expect(myMap).Should(HaveKeyWithValue("Hello", "World"))

		})

		It("should inject not struct components", func() {

			container.Put(NewMap("Hello", "World"), "MyMap")
			container.PutFactory(CreateWrapMap, "Wrap")

			wrap := container.GetComponent("Wrap").(MyMap)

			Expect(wrap).Should(HaveLen(1))
			Expect(wrap).Should(HaveKeyWithValue("Hello", "World"))

		})

		It("should inject not struct components in a slice", func() {

			container.Put(NewMap("Hello", "World"), "MyMap1", "Maps")
			container.Put(NewMap("Hi", "Everybody"), "MyMap2", "Maps")

			container.PutFactory(CreateSuperMap, "Super")

			super := container.GetComponent("Super").(MySuperMap)

			Expect(super).Should(HaveLen(2))
			Expect(super).Should(HaveKeyWithValue("Hello", "World"))
			Expect(super).Should(HaveKeyWithValue("Hi", "Everybody"))

		})

		It("should inject not struct components in a map", func() {

			container.Put(NewMap("Hello", "World"), "MyMap1", "Maps")
			container.Put(NewMap("Hi", "Everybody"), "MyMap2", "Maps")

			container.PutFactory(CreateHyperMap, "Hyper")

			hyper := container.GetComponent("Hyper").(MyHyperMap)

			Expect(hyper).Should(HaveLen(2))
			Expect(hyper).Should(HaveKeyWithValue("MyMap1.Hello", "World"))
			Expect(hyper).Should(HaveKeyWithValue("MyMap2.Hi", "Everybody"))

		})

	})

})

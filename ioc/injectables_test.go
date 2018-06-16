package ioc_test

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

type InjectedObject struct {
	A *SimpleStruct `inject:""`
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

type SliceInjectedObject struct {
	SimpleStructs []*SimpleStruct `inject:"A"`
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

type MapInjectedObject struct {
	SimpleStructs map[string]*SimpleStruct `inject:"A"`
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
func (self *OrderRegister) registerClose(obj Ordered) error {
	self.CloseOrder = append(self.CloseOrder, obj.Name())
	return nil
}

type First struct {
	OrderRegister *OrderRegister `inject:""`
}

func (self *First) Name() string { return "FIRST" }
func (self *First) PostInit()    { self.OrderRegister.registerPostInit(self) }
func (self *First) Close() error { return self.OrderRegister.registerClose(self) }

type Second struct {
	OrderRegister *OrderRegister `inject:""`
	Injected      *First         `inject:"First"`
}

func (self *Second) Name() string { return "SECOND" }
func (self *Second) PostInit()    { self.OrderRegister.registerPostInit(self) }
func (self *Second) Close() error { return self.OrderRegister.registerClose(self) }

type Third struct {
	OrderRegister *OrderRegister `inject:""`
	Injected      *Second        `inject:"Second"`
}

func (self *Third) Name() string { return "THIRD" }
func (self *Third) PostInit()    { self.OrderRegister.registerPostInit(self) }
func (self *Third) Close() error { return self.OrderRegister.registerClose(self) }

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

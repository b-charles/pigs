package ioc_test

// Doer

type Doer interface {
	do()
}

type BigDoer interface {
	bigDo()
}

type NotDoer interface {
	notDo()
}

// Simple

type Simple struct {
	Tag string
}

func (self *Simple) String() string {
	return self.Tag
}

func (self *Simple) do() {}

func (self *Simple) bigDo() {}

func SimpleFactory(tag string) func() *Simple {
	return func() *Simple { return &Simple{tag} }
}

// Trivial

type Trivial string

func (self Trivial) do() {}

func TrivialFactory(tag string) func() Trivial {
	return func() Trivial { return (Trivial)(tag) }
}

// Injected

type Injected struct {
	Simple *Simple
}

func InjectedFactory(simple *Simple) *Injected {
	return &Injected{simple}
}

// Initialized

type Initialized struct {
	Simple *Simple `inject:""`
}

// PostInitialized

type PostInitialized struct {
	Simple *Simple
}

func (self *PostInitialized) PostInit(simple *Simple) {
	self.Simple = simple
}

// InterfaceInjected

type InterfaceInjected struct {
	Doer Doer
}

func InterfaceInjectedFactory(doer Doer) *InterfaceInjected {
	return &InterfaceInjected{doer}
}

// InterfaceInitialized

type InterfaceInitialized struct {
	Doer Doer `inject:""`
}

// InterfacePostInitialized

type InterfacePostInitialized struct {
	Doer Doer
}

func (self *InterfacePostInitialized) PostInit(doer Doer) {
	self.Doer = doer
}

// SliceInjectedStruct

type SliceInjected struct {
	Simple []*Simple `inject:""`
}

// InterfaceSliceInjected

type InterfaceSliceInjected struct {
	Doers []Doer `inject:""`
}

// Looping (inject itself)

type Looping struct {
	Looping *Looping `inject:""`
}

// Ordered

type OrderedComp interface {
	Name() string
}

type OrderRegister struct {
	PostInitOrder []string
	CloseOrder    []string
}

func NewOrderRegister() *OrderRegister {
	return &OrderRegister{make([]string, 0), make([]string, 0)}
}

func (self *OrderRegister) registerPostInit(obj OrderedComp) {
	self.PostInitOrder = append(self.PostInitOrder, obj.Name())
}

func (self *OrderRegister) registerClose(obj OrderedComp) error {
	self.CloseOrder = append(self.CloseOrder, obj.Name())
	return nil
}

type First struct {
	Register *OrderRegister `inject:""`
}

func (self *First) Name() string { return "FIRST" }
func (self *First) PostInit()    { self.Register.registerPostInit(self) }
func (self *First) Close() error { return self.Register.registerClose(self) }

type Second struct {
	Register *OrderRegister `inject:""`
	Injected *First         `inject:""`
}

func (self *Second) Name() string { return "SECOND" }
func (self *Second) PostInit()    { self.Register.registerPostInit(self) }
func (self *Second) Close() error { return self.Register.registerClose(self) }

type Third struct {
	Register *OrderRegister `inject:""`
	Injected *Second        `inject:""`
}

func (self *Third) Name() string { return "THIRD" }
func (self *Third) PostInit()    { self.Register.registerPostInit(self) }
func (self *Third) Close() error { return self.Register.registerClose(self) }

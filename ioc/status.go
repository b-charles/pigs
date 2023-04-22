package ioc

import (
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

var (
	colorReset  = "\033[0m"
	colorFaint  = "\033[2m"
	colorStrike = "\033[9m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorPurple = "\033[35m"
	colorCyan   = "\033[36m"
)

// Types utilities

func typeElem(typ reflect.Type) (reflect.Type, bool) {
	if typ.Kind() == reflect.Pointer {
		return typ.Elem(), true
	} else {
		return typ, false
	}
}

func typeLess(a, b reflect.Type) bool {
	ae, _ := typeElem(a)
	be, _ := typeElem(b)
	return ae.String() < be.String()
}

// Component wrapper

type ComponentRecord interface {
	Name() string
	Type() reflect.Type
	Instanciated() bool
}

type componentRecordImpl struct {
	name         string
	typ          reflect.Type
	instanciated bool
}

func (self *componentRecordImpl) Name() string {
	return self.name
}

func (self *componentRecordImpl) Type() reflect.Type {
	return self.typ
}

func (self *componentRecordImpl) Instanciated() bool {
	return self.instanciated
}

type ComponentType interface {
	Type() reflect.Type
	NumComponents() int
	GetComponent(int) ComponentRecord
	Overloaded() bool
}

type componentTypeImpl struct {
	typ        reflect.Type
	components []*componentRecordImpl
	overloaded bool
}

func (self *componentTypeImpl) Type() reflect.Type {
	return self.typ
}

func (self *componentTypeImpl) NumComponents() int {
	return len(self.components)
}

func (self *componentTypeImpl) GetComponent(i int) ComponentRecord {
	return self.components[i]
}

func (self *componentTypeImpl) Overloaded() bool {
	return self.overloaded
}

func (self *componentTypeImpl) write(builder *strings.Builder) {

	if len(self.components) == 1 && self.typ == self.components[0].typ {

		name := self.components[0].name
		instanciated := self.components[0].instanciated

		if _, p := typeElem(self.typ); !p {
			builder.WriteString(" ")
		}

		if self.overloaded {
			builder.WriteString(colorStrike)
		} else if !instanciated {
			builder.WriteString(colorFaint)
		}

		builder.WriteString(self.typ.String())

		if name != "" {
			builder.WriteString(" - ")
			builder.WriteString(name)
		}

		if self.overloaded {
			builder.WriteString(colorReset)
			builder.WriteString(" ")
			builder.WriteString(colorFaint)
			builder.WriteString("Overloaded")
			builder.WriteString(colorReset)
		} else if !instanciated {
			builder.WriteString(colorReset)
		}

		builder.WriteString("\n")

	} else {

		if _, p := typeElem(self.typ); !p {
			builder.WriteString(" ")
		}

		if self.overloaded {
			builder.WriteString(colorStrike)
		}

		builder.WriteString(self.typ.String())

		if self.overloaded {
			builder.WriteString(colorReset)
			builder.WriteString(" ")
			builder.WriteString(colorFaint)
			builder.WriteString("Overloaded")
			builder.WriteString(colorReset)
		}

		builder.WriteString(" (")
		builder.WriteString(strconv.Itoa(len(self.components)))
		builder.WriteString(")\n")

		for _, t := range self.components {
			builder.WriteString("   |-> ")
			if _, p := typeElem(t.typ); !p {
				builder.WriteString(" ")
			}
			if !t.instanciated {
				builder.WriteString(colorFaint)
			}
			builder.WriteString(t.typ.String())
			if name := t.name; name != "" {
				builder.WriteString(" - ")
				builder.WriteString(name)
			}
			if !t.instanciated {
				builder.WriteString(colorReset)
			}
			builder.WriteString("\n")
		}

	}

}

// Instance wrapper

type ComponentInstance interface {
	Scope() Scope
	Name() string
	Type() reflect.Type
	Value() any
	Closable() bool
}

type componentInstanceImpl struct {
	scope    Scope
	name     string
	typ      reflect.Type
	value    any
	closable bool
}

func (self *componentInstanceImpl) Scope() Scope {
	return self.scope
}

func (self *componentInstanceImpl) Name() string {
	return self.name
}

func (self *componentInstanceImpl) Type() reflect.Type {
	return self.typ
}

func (self *componentInstanceImpl) Value() any {
	return self.value
}

func (self *componentInstanceImpl) Closable() bool {
	return self.closable
}

func (self *componentInstanceImpl) write(builder *strings.Builder) {

	if self.scope == Core {
		builder.WriteString("[")
		builder.WriteString(colorCyan)
		builder.WriteString("Core")
		builder.WriteString(colorReset)
		builder.WriteString("] ")
	} else if self.scope == Def {
		builder.WriteString("[")
		builder.WriteString(colorYellow)
		builder.WriteString("Def")
		builder.WriteString(colorReset)
		builder.WriteString("]  ")
	} else { // scope == test
		builder.WriteString("[")
		builder.WriteString(colorPurple)
		builder.WriteString("Test")
		builder.WriteString(colorReset)
		builder.WriteString("] ")
	}

	if _, p := typeElem(self.typ); !p {
		builder.WriteString(" ")
	}
	builder.WriteString(self.typ.String())
	if name := self.name; name != "" {
		builder.WriteString(" - ")
		builder.WriteString(name)
	}

	if self.closable {
		builder.WriteString(" (")
		builder.WriteString(colorYellow)
		builder.WriteString("Closable")
		builder.WriteString(colorReset)
		builder.WriteString(")")
	}

	if stringer, ok := self.value.(fmt.Stringer); ok {
		builder.WriteString(":\n    ")
		if str := stringer.String(); str != "" {
			builder.WriteString(strings.ReplaceAll(str, "\n", "\n    "))
		}
	}
	builder.WriteString("\n")

}

// Container status

type ContainerStatus interface {
	NumDefaultTypes() int
	DefaultType(int) ComponentType
	NumCoreTypes() int
	CoreType(int) ComponentType
	NumTestTypes() int
	TestType(int) ComponentType
	NumInstances() int
	Instance(int) ComponentInstance
	String() string
	Print()
}

type containerStatusImpl struct {
	def       []*componentTypeImpl
	core      []*componentTypeImpl
	test      []*componentTypeImpl
	instances []*componentInstanceImpl
}

func (self *containerStatusImpl) NumDefaultTypes() int {
	return len(self.def)
}

func (self *containerStatusImpl) DefaultType(i int) ComponentType {
	return self.def[i]
}

func (self *containerStatusImpl) NumCoreTypes() int {
	return len(self.core)
}

func (self *containerStatusImpl) CoreType(i int) ComponentType {
	return self.core[i]
}

func (self *containerStatusImpl) NumTestTypes() int {
	return len(self.test)
}

func (self *containerStatusImpl) TestType(i int) ComponentType {
	return self.test[i]
}

func (self *containerStatusImpl) NumInstances() int {
	return len(self.instances)
}

func (self *containerStatusImpl) Instance(i int) ComponentInstance {
	return self.instances[i]
}

func (self *containerStatusImpl) update(container *Container) {

	// Default

	self.def = make([]*componentTypeImpl, 0, len(container.defaultComponents))
	for typ, comps := range container.defaultComponents {

		components := make([]*componentRecordImpl, 0, len(comps))
		for _, comp := range comps {
			_, instanciated := container.instances[comp]
			components = append(components, &componentRecordImpl{
				name:         comp.name,
				typ:          comp.main,
				instanciated: instanciated,
			})
		}
		sort.Slice(components, func(i, j int) bool {
			return typeLess(components[i].typ, components[j].typ)
		})

		_, overloaded := container.coreComponents[typ]
		if !overloaded {
			_, overloaded = container.testComponents[typ]
		}

		self.def = append(self.def, &componentTypeImpl{
			typ:        typ,
			components: components,
			overloaded: overloaded,
		})

	}
	sort.Slice(self.def, func(i, j int) bool {
		return typeLess(self.def[i].typ, self.def[j].typ)
	})

	// Core

	self.core = make([]*componentTypeImpl, 0, len(container.coreComponents))
	for typ, comps := range container.coreComponents {

		components := make([]*componentRecordImpl, 0, len(comps))
		for _, comp := range comps {
			_, instanciated := container.instances[comp]
			components = append(components, &componentRecordImpl{
				name:         comp.name,
				typ:          comp.main,
				instanciated: instanciated,
			})
		}
		sort.Slice(components, func(i, j int) bool {
			return typeLess(components[i].typ, components[j].typ)
		})

		_, overloaded := container.testComponents[typ]

		self.core = append(self.core, &componentTypeImpl{
			typ:        typ,
			components: components,
			overloaded: overloaded,
		})

	}
	sort.Slice(self.core, func(i, j int) bool {
		return typeLess(self.core[i].typ, self.core[j].typ)
	})

	// Test

	self.test = make([]*componentTypeImpl, 0, len(container.testComponents))
	for typ, comps := range container.testComponents {

		components := make([]*componentRecordImpl, 0, len(comps))
		for _, comp := range comps {
			_, instanciated := container.instances[comp]
			components = append(components, &componentRecordImpl{
				name:         comp.name,
				typ:          comp.main,
				instanciated: instanciated,
			})
		}
		sort.Slice(components, func(i, j int) bool {
			return typeLess(components[i].typ, components[j].typ)
		})

		self.test = append(self.test, &componentTypeImpl{
			typ:        typ,
			components: components,
			overloaded: false,
		})

	}
	sort.Slice(self.test, func(i, j int) bool {
		return typeLess(self.test[i].typ, self.test[j].typ)
	})

	// Instances

	self.instances = make([]*componentInstanceImpl, 0, len(container.instances))
	for comp, inst := range container.instances {

		value := inst.value.Interface()
		if value == self {
			value = ""
		}

		var scope Scope
		if _, p := container.testComponents[comp.main]; p {
			scope = Test
		} else if _, p := container.coreComponents[comp.main]; p {
			scope = Core
		} else {
			scope = Def
		}

		self.instances = append(self.instances, &componentInstanceImpl{
			scope:    scope,
			name:     comp.name,
			typ:      comp.main,
			value:    value,
			closable: inst.isClosable(),
		})

	}
	sort.Slice(self.instances, func(i, j int) bool {
		return typeLess(self.instances[i].typ, self.instances[j].typ)
	})

}

func (self *containerStatusImpl) String() string {

	var builder strings.Builder

	builder.WriteString(colorReset)
	builder.WriteString("\n\n Container status \n")
	builder.WriteString(colorRed)
	builder.WriteString("------------------\n")
	builder.WriteString(colorReset)

	builder.WriteString("\n")
	builder.WriteString(colorGreen)
	builder.WriteString("###")
	builder.WriteString(colorReset)
	builder.WriteString(" Default components (")
	builder.WriteString(fmt.Sprintf("%d", len(self.def)))
	builder.WriteString(")\n\n")
	for _, def := range self.def {
		def.write(&builder)
	}

	builder.WriteString("\n")
	builder.WriteString(colorGreen)
	builder.WriteString("###")
	builder.WriteString(colorReset)
	builder.WriteString(" Core components (")
	builder.WriteString(fmt.Sprintf("%d", len(self.core)))
	builder.WriteString(")\n\n")
	for _, core := range self.core {
		core.write(&builder)
	}

	builder.WriteString("\n")
	builder.WriteString(colorGreen)
	builder.WriteString("###")
	builder.WriteString(colorReset)
	builder.WriteString(" Test components (")
	builder.WriteString(fmt.Sprintf("%d", len(self.test)))
	builder.WriteString(")\n\n")
	for _, test := range self.test {
		test.write(&builder)
	}

	builder.WriteString("\n")
	builder.WriteString(colorGreen)
	builder.WriteString("###")
	builder.WriteString(colorReset)
	builder.WriteString(" Component instances (")
	builder.WriteString(fmt.Sprintf("%d", len(self.instances)))
	builder.WriteString(")\n\n")
	for _, inst := range self.instances {
		inst.write(&builder)
	}

	builder.WriteString("\n")

	return builder.String()

}

func (self *containerStatusImpl) Print() {
	fmt.Printf(self.String())
}

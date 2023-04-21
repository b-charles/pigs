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

func typeAlignString(typ reflect.Type) string {
	if t, p := typeElem(typ); p {
		return typ.String()
	} else {
		return fmt.Sprintf(" %s", t.String())
	}
}

func typeLess(a, b reflect.Type) bool {
	ae, _ := typeElem(a)
	be, _ := typeElem(b)
	return ae.String() < be.String()
}

// Component wrapper

type ComponentRecords struct {
	Type       reflect.Type
	From       []reflect.Type
	Overloaded bool
}

func (self ComponentRecords) write(builder *strings.Builder) {

	builder.WriteString(typeAlignString(self.Type))

	if self.Overloaded {
		builder.WriteString(" (")
		builder.WriteString(colorPurple)
		builder.WriteString("Overloaded")
		builder.WriteString(colorReset)
		builder.WriteString(")")
	}

	if len(self.From) == 1 && self.Type == self.From[0] {
		builder.WriteString("\n")
	} else {

		builder.WriteString(" (")
		builder.WriteString(strconv.Itoa(len(self.From)))
		builder.WriteString(")\n")

		for _, t := range self.From {
			builder.WriteString("   |-> ")
			builder.WriteString(typeAlignString(t))
			builder.WriteString("\n")
		}

	}

}

func sortComponentsRecords(slice []ComponentRecords) {
	sort.Slice(slice, func(i, j int) bool {
		return typeLess(slice[i].Type, slice[j].Type)
	})
}

// Instance wrapper

type InstanceRecords struct {
	Type     reflect.Type
	Instance any
	Closable bool
}

func (self InstanceRecords) write(builder *strings.Builder) {

	builder.WriteString(typeAlignString(self.Type))
	if self.Closable {
		builder.WriteString(" (")
		builder.WriteString(colorYellow)
		builder.WriteString("Closable")
		builder.WriteString(colorReset)
		builder.WriteString(")")
	}
	builder.WriteString(": ")
	builder.WriteString(fmt.Sprintf("%v", self.Instance))
	builder.WriteString("\n")

}

// Container status

type ContainerStatus struct {
	container *Container
	Default   []ComponentRecords
	Core      []ComponentRecords
	Test      []ComponentRecords
	Instances []InstanceRecords
}

var containerStatus_type reflect.Type = reflect.TypeOf(&ContainerStatus{})

func (self *ContainerStatus) update() {

	// Default

	self.Default = make([]ComponentRecords, 0, len(self.container.defaultComponents))
	for typ, comps := range self.container.defaultComponents {

		from := make([]reflect.Type, 0, len(comps))
		for _, comp := range comps {
			from = append(from, comp.main)
		}
		sort.Slice(from, func(i, j int) bool {
			return typeLess(from[i], from[j])
		})

		_, over := self.container.coreComponents[typ]
		if !over {
			_, over = self.container.testComponents[typ]
		}

		self.Default = append(self.Default, ComponentRecords{
			Type:       typ,
			From:       from,
			Overloaded: over,
		})

	}
	sort.Slice(self.Default, func(i, j int) bool {
		return typeLess(self.Default[i].Type, self.Default[j].Type)
	})

	// Core

	self.Core = make([]ComponentRecords, 0, len(self.container.coreComponents))
	for typ, comps := range self.container.coreComponents {

		from := make([]reflect.Type, 0, len(comps))
		for _, comp := range comps {
			from = append(from, comp.main)
		}
		sort.Slice(from, func(i, j int) bool {
			return typeLess(from[i], from[j])
		})

		_, over := self.container.testComponents[typ]

		self.Core = append(self.Core, ComponentRecords{
			Type:       typ,
			From:       from,
			Overloaded: over,
		})

	}
	sort.Slice(self.Core, func(i, j int) bool {
		return typeLess(self.Core[i].Type, self.Core[j].Type)
	})

	// Test

	self.Test = make([]ComponentRecords, 0, len(self.container.testComponents))
	for typ, comps := range self.container.testComponents {

		from := make([]reflect.Type, 0, len(comps))
		for _, comp := range comps {
			from = append(from, comp.main)
		}
		sort.Slice(from, func(i, j int) bool {
			return typeLess(from[i], from[j])
		})

		self.Test = append(self.Test, ComponentRecords{
			Type:       typ,
			From:       from,
			Overloaded: false,
		})

	}
	sort.Slice(self.Test, func(i, j int) bool {
		return typeLess(self.Test[i].Type, self.Test[j].Type)
	})

	// Instances

	self.Instances = make([]InstanceRecords, 0, len(self.container.instances))
	for comp, inst := range self.container.instances {

		value := inst.value.Interface()
		if value == self {
			value = "<Container status>"
		}

		self.Instances = append(self.Instances, InstanceRecords{
			Type:     comp.main,
			Instance: value,
			Closable: inst.isClosable(),
		})

	}
	sort.Slice(self.Instances, func(i, j int) bool {
		return typeLess(self.Instances[i].Type, self.Instances[j].Type)
	})

}

func (self *ContainerStatus) String() string {

	var builder strings.Builder

	builder.WriteString(colorReset)
	builder.WriteString("\n Container status \n")
	builder.WriteString(colorRed)
	builder.WriteString("------------------\n")
	builder.WriteString(colorReset)

	builder.WriteString("\n")
	builder.WriteString(colorGreen)
	builder.WriteString("###")
	builder.WriteString(colorReset)
	builder.WriteString(" Default components (")
	builder.WriteString(fmt.Sprintf("%d", len(self.Default)))
	builder.WriteString(")\n")
	for _, def := range self.Default {
		def.write(&builder)
	}

	builder.WriteString("\n")
	builder.WriteString(colorGreen)
	builder.WriteString("###")
	builder.WriteString(colorReset)
	builder.WriteString(" Core components (")
	builder.WriteString(fmt.Sprintf("%d", len(self.Core)))
	builder.WriteString(")\n")
	for _, core := range self.Core {
		core.write(&builder)
	}

	builder.WriteString("\n")
	builder.WriteString(colorGreen)
	builder.WriteString("###")
	builder.WriteString(colorReset)
	builder.WriteString(" Test components (")
	builder.WriteString(fmt.Sprintf("%d", len(self.Test)))
	builder.WriteString(")\n")
	for _, test := range self.Test {
		test.write(&builder)
	}

	builder.WriteString("\n")
	builder.WriteString(colorGreen)
	builder.WriteString("###")
	builder.WriteString(colorReset)
	builder.WriteString(" Component instances (")
	builder.WriteString(fmt.Sprintf("%d", len(self.Instances)))
	builder.WriteString(")\n")
	for _, inst := range self.Instances {
		inst.write(&builder)
	}

	builder.WriteString("\n")

	return builder.String()

}

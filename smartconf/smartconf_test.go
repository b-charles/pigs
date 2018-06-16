package smartconf_test

import (
	"fmt"
	"strconv"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/l3eegbee/pigs/smartconf"
)

func TestConfig(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Smart Config Suite")
}

// Simple types

type SimpleHolder struct {
	BoolValue    bool
	IntValue     int
	Int8Value    int8
	Int16Value   int16
	Int32Value   int32
	Int64Value   int64
	UintValue    uint
	Uint8Value   uint8
	Uint16Value  uint16
	Uint32Value  uint32
	Uint64Value  uint64
	Float32Value float32
	Float64Value float64
	StringValue  string
}

// Parse interfaces

type ConfInterface interface {
	DoSomething() bool
}

func ConfInterfaceParser(value string) (ConfInterface, error) {
	if b, ok := strconv.ParseBool(value); ok == nil {
		return &BooleanHolder{b}, nil
	} else if i, ok := strconv.ParseInt(value, 0, 0); ok == nil {
		return &IntHolder{int(i)}, nil
	} else {
		return nil, fmt.Errorf("Error parse generic interface.")
	}
}

type BooleanHolder struct {
	Value bool
}

func (self *BooleanHolder) DoSomething() bool {
	return self.Value
}

type IntHolder struct {
	Value int
}

func (self *IntHolder) DoSomething() bool {
	return self.Value > 0
}

type InterfaceHolder struct {
	Value ConfInterface
}

// Map

type MapInterfaceHolder struct {
	Value map[string]ConfInterface
}

// Ptr

type PtrHolder struct {
	Int *IntHolder
}

// Slice

type SliceHolder struct {
	Elt []IntHolder
}

// Struct

type StructHolder struct {
	Inter  ConfInterface
	Simple IntHolder
	Yesno  *BooleanHolder
}

var _ = Describe("Smart Config", func() {

	var smart *SmartConf = NewSmartConf(map[string]string{
		"simple.bool.value":      "true",
		"simple.int.value":       "-42",
		"simple.int8.value":      "0x01",
		"simple.int16.value":     "0x0123",
		"simple.int32.value":     "0x01234567",
		"simple.int64.value":     "0x0123456789ABCDEF",
		"simple.uint.value":      "42",
		"simple.uint8.value":     "0x42",
		"simple.uint16.value":    "0xDEAD",
		"simple.uint32.value":    "0xDEADBEEF",
		"simple.uint64.value":    "0xDEADFACEBEEFCAFE",
		"simple.float32.value":   "3.141592653589793238462643383279",
		"simple.float64.value":   "3.14159265358979323846264338327950288419716939937510582097494459",
		"simple.string.value":    "Hello, my name is Pigs.",
		"interface.first.value":  "true",
		"interface.second.value": "42",
		"map.value.one":          "54",
		"map.value.two":          "true",
		"map.value.three":        "-5",
		"map.value.four":         "false",
		"ptr.int.value":          "5555",
		"slice.elt.0.value":      "5",
		"slice.elt.1.value":      "4",
		"slice.elt.2.value":      "3",
		"slice.elt.3.value":      "2",
		"slice.elt.4.value":      "1",
		"struct.inter":           "36",
		"struct.simple.value":    "27",
		"struct.yesno.value":     "true",
	}, []interface{}{
		ConfInterfaceParser,
	})

	It("Should config a simple struct", func() {

		valueHolder := smart.MustConfigure("simple", (*SimpleHolder)(nil)).(*SimpleHolder)

		Expect(valueHolder).Should(Equal(&SimpleHolder{
			BoolValue:    true,
			IntValue:     -42,
			Int8Value:    0x01,
			Int16Value:   0x0123,
			Int32Value:   0x01234567,
			Int64Value:   0x0123456789ABCDEF,
			UintValue:    42,
			Uint8Value:   0x42,
			Uint16Value:  0xDEAD,
			Uint32Value:  0xDEADBEEF,
			Uint64Value:  0xDEADFACEBEEFCAFE,
			Float32Value: 3.141592653589793238462643383279,
			Float64Value: 3.14159265358979323846264338327950288419716939937510582097494459,
			StringValue:  "Hello, my name is Pigs.",
		}))

	})

	It("Should parse to interfaces", func() {

		first := smart.MustConfigure("interface.first", (*InterfaceHolder)(nil)).(*InterfaceHolder)
		Expect(first).Should(Equal(&InterfaceHolder{&BooleanHolder{true}}))

		second := smart.MustConfigure("interface.second", (*InterfaceHolder)(nil)).(*InterfaceHolder)
		Expect(second).Should(Equal(&InterfaceHolder{&IntHolder{42}}))

	})

	It("Should parse a map", func() {

		mapHolder := smart.MustConfigure("map", (*MapInterfaceHolder)(nil)).(*MapInterfaceHolder)

		Expect(mapHolder.Value).Should(HaveLen(4))
		Expect(mapHolder.Value).Should(HaveKeyWithValue("one", &IntHolder{54}))
		Expect(mapHolder.Value).Should(HaveKeyWithValue("two", &BooleanHolder{true}))
		Expect(mapHolder.Value).Should(HaveKeyWithValue("three", &IntHolder{-5}))
		Expect(mapHolder.Value).Should(HaveKeyWithValue("four", &BooleanHolder{false}))

	})

	It("Should handle pointer", func() {

		ptrHolder := smart.MustConfigure("ptr", (*PtrHolder)(nil)).(*PtrHolder)
		Expect(ptrHolder).Should(Equal(&PtrHolder{&IntHolder{5555}}))

	})

	It("Should parse slices", func() {

		sliceHolder := smart.MustConfigure("slice", (*SliceHolder)(nil)).(*SliceHolder)

		Expect(sliceHolder.Elt).Should(HaveLen(5))
		Expect(sliceHolder.Elt[0]).Should(Equal(IntHolder{5}))
		Expect(sliceHolder.Elt[1]).Should(Equal(IntHolder{4}))
		Expect(sliceHolder.Elt[2]).Should(Equal(IntHolder{3}))
		Expect(sliceHolder.Elt[3]).Should(Equal(IntHolder{2}))
		Expect(sliceHolder.Elt[4]).Should(Equal(IntHolder{1}))

	})

	It("Should do struct and all", func() {

		structHolder := smart.MustConfigure("struct", (*StructHolder)(nil)).(*StructHolder)

		Expect(structHolder).Should(Equal(&StructHolder{
			Inter:  &IntHolder{36},
			Simple: IntHolder{27},
			Yesno:  &BooleanHolder{true},
		}))

	})

})

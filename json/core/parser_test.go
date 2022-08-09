package core_test

import (
	"strconv"

	. "github.com/b-charles/pigs/json/core"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func testInt(s string) {
	if v, err := strconv.Atoi(s); err != nil {
		panic(err)
	} else {
		Expect(ParseString(s)).WithOffset(1).To(Equal(JsonInt(v)))
	}
}

func testFloat(s string) {
	if v, err := strconv.ParseFloat(s, 64); err != nil {
		panic(err)
	} else {
		Expect(ParseString(s)).WithOffset(1).To(Equal(JsonFloat(v)))
	}
}

var _ = Describe("Json parser", func() {

	It("should be able to parse simple constants", func() {

		Expect(ParseString(`true`)).To(Equal(JSON_TRUE))
		Expect(ParseString(`false`)).To(Equal(JSON_FALSE))
		Expect(ParseString(`null`)).To(Equal(JSON_NULL))

	})

	It("should be able to parse integers", func() {

		testInt(`0`)
		testInt(`4`)
		testInt(`-4`)
		testInt(`42`)
		testInt(`-42`)

	})

	It("should be able to parse floats", func() {

		testFloat(`0.0`)
		testFloat(`10.04`)
		testFloat(`-10.04`)
		testFloat(`2e04`)
		testFloat(`0.2e-4`)
		testFloat(`0.2e+4`)
		testFloat(`-10.04e7`)

	})

	It("sould be able to parse empty object", func() {
		Expect(ParseString(`{}`)).To(Equal(JSON_EMPTY_OBJECT))
		Expect(ParseString(`{ }`)).To(Equal(JSON_EMPTY_OBJECT))
		Expect(ParseString(` { } `)).To(Equal(JSON_EMPTY_OBJECT))
	})

	It("sould be able to parse simple object", func() {
		expected := NewJsonBuilder().SetBool(`a`, true).Build()
		Expect(ParseString(`{"a":true}`)).To(Equal(expected))
	})

	It("sould be able to parse less simple object", func() {
		expected := NewJsonBuilder().
			SetBool(`a`, true).
			SetInt(`b`, 42).
			Build()
		Expect(ParseString(`{"a":true,"b":42}`)).To(Equal(expected))
	})

	It("sould be able to parse nested object", func() {
		expected := NewJsonBuilder().SetBool(`a.b`, true).Build()
		Expect(ParseString(`{"a":{"b":true}}`)).To(Equal(expected))
	})

	It("sould be able to parse empty array", func() {
		Expect(ParseString(`[]`)).To(Equal(JSON_EMPTY_ARRAY))
		Expect(ParseString(`[ ]`)).To(Equal(JSON_EMPTY_ARRAY))
		Expect(ParseString(` [ ] `)).To(Equal(JSON_EMPTY_ARRAY))
	})

	It("sould be able to parse simple array", func() {
		expected := NewJsonBuilder().SetBool(`[0]`, true).Build()
		Expect(ParseString(`[true]`)).To(Equal(expected))
	})

	It("sould be able to parse simple array", func() {
		expected := NewJsonBuilder().SetBool(`[0]`, true).Build()
		Expect(ParseString(`[true]`)).To(Equal(expected))
	})

	It("sould be able to parse less simple arryay", func() {
		expected := NewJsonBuilder().
			SetBool(`[0]`, true).
			SetInt(`[1]`, 42).
			Build()
		Expect(ParseString(`[true, 42]`)).To(Equal(expected))
	})

	It("sould be able to parse nested array", func() {
		expected := NewJsonBuilder().SetBool(`[0].[0]`, true).Build()
		Expect(ParseString(`[[true]]`)).To(Equal(expected))
	})

	It("sould be able to parse complicated object", func() {
		expected := NewJsonBuilder().
			SetString(`terran[0].name`, "Raynor, Jim").
			SetBool(`terran[0].alive`, true).
			SetString(`terran[1].name`, "Mengsk, Arcturus").
			SetBool(`terran[1].alive`, false).
			SetString(`zerg[0].name`, "Kerrigan, Sarah").
			SetBool(`zerg[0].alive`, true).
			SetString(`protoss[0].name`, "Tassadar").
			SetBool(`protoss[0].alive`, false).
			Build()
		Expect(ParseString(`{
			"terran":[
				{"name": "Raynor, Jim", "alive": true},
				{"name": "Mengsk, Arcturus", "alive": false}
			],
			"zerg":[
				{"name": "Kerrigan, Sarah", "alive": true}
			],
			"protoss":[
				{"name": "Tassadar", "alive": false}
			]}`)).To(Equal(expected))
	})

})

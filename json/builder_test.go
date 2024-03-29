package json_test

import (
	. "github.com/b-charles/pigs/json"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Json builder", func() {

	It("should be able to build empty object", func() {

		b := NewJsonBuilder()
		Expect(b.Build().String()).To(Equal("null"))

	})

	It("should be able to build a simple value", func() {

		b := NewJsonBuilder()
		b.SetInt("val", 42)
		Expect(b.Build().String()).To(Equal(`{"val":42}`))

	})

	It("should handle empty field", func() {

		b := NewJsonBuilder()
		b.SetInt("a..b", 42)
		Expect(b.Build().String()).To(Equal(`{"a":{"":{"b":42}}}`))

	})

	It("should handle protected fields (dot)", func() {

		b := NewJsonBuilder()
		b.SetInt("a\\.b", 42)
		Expect(b.Build().String()).To(Equal(`{"a.b":42}`))

	})

	It("should handle protected fields (bracket)", func() {

		b := NewJsonBuilder()
		b.SetInt("a\\[b]", 42)
		Expect(b.Build().String()).To(Equal(`{"a[b]":42}`))

	})

	It("should be able to build complex value", func() {

		b := NewJsonBuilder()
		b.SetString("a.b.c", "hello")
		b.SetString("a.b.d\\.e", "world")
		b.SetInt("a.a", 42)
		b.SetBool("a.b.a[4].e", true)
		Expect(b.Build().String()).To(Equal(`{"a":{"b":{"c":"hello","d.e":"world","a":[null,null,null,null,{"e":true}]},"a":42}}`))

	})

	It("should be able to handle weird strings", func() {

		b := NewJsonBuilder()
		b.SetString("a.a", "♪🤠\n")
		Expect(b.Build().String()).To(Equal(`{"a":{"a":"\u266a\ud83e\udd20\n"}}`))

	})

	It("should overwrite values", func() {

		b := NewJsonBuilder()
		b.SetInt("a.a", 42)
		b.SetBool("a.b[0]", true)
		b.SetInt("a.a", 45)
		b.SetBool("a.b[0]", false)
		Expect(b.Build().String()).To(Equal(`{"a":{"a":45,"b":[false]}}`))

	})

})

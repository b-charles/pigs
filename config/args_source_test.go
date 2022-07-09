package config_test

import (
	. "github.com/b-charles/pigs/config"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Args", func() {

	It("should parse value between simple quote", func() {

		source := ParseArgs([]string{
			"--jamiroquai='Virtual Insanity'",
		})

		Expect(source).Should(HaveKeyWithValue("jamiroquai", "Virtual Insanity"))

	})

	It("should parse value between double quote", func() {

		source := ParseArgs([]string{
			"--santana=\"Flor D'Luna\"",
		})

		Expect(source).Should(HaveKeyWithValue("santana", "Flor D'Luna"))

	})

	It("should parse boolean", func() {

		source := ParseArgs([]string{
			"--yes",
		})

		Expect(source).Should(HaveKeyWithValue("yes", "true"))

	})

	It("should parse false boolean", func() {

		source := ParseArgs([]string{
			"--no-yes",
		})

		Expect(source).Should(HaveKeyWithValue("yes", "false"))

	})

})

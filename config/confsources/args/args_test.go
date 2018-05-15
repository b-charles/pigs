package args_test

import (
	"testing"

	. "github.com/l3eegbee/pigs/config/confsources/args"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestConfsources(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Args config sources Suite")
}

var _ = Describe("Args", func() {

	It("should parse value between simple quote", func() {

		source := NewArgsConfigSource([]string{
			"--jamiroquai='Virtual Insanity'",
		})

		Expect(source.LoadEnv()).Should(HaveKeyWithValue("jamiroquai", "Virtual Insanity"))

	})

	It("should parse value between double quote", func() {

		source := NewArgsConfigSource([]string{
			"--santana=\"Flor D'Luna\"",
		})

		Expect(source.LoadEnv()).Should(HaveKeyWithValue("santana", "Flor D'Luna"))

	})

	It("should parse boolean", func() {

		source := NewArgsConfigSource([]string{
			"--yes",
		})

		Expect(source.LoadEnv()).Should(HaveKeyWithValue("yes", "true"))

	})

	It("should parse false boolean", func() {

		source := NewArgsConfigSource([]string{
			"--no-yes",
		})

		Expect(source.LoadEnv()).Should(HaveKeyWithValue("yes", "false"))

	})

})

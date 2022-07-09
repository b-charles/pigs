package config_test

import (
	. "github.com/b-charles/pigs/config"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Default config source", func() {

	var (
		source *DefaultConfigSource
	)

	BeforeEach(func() {
		source = NewDefaultConfigSource()
	})

	It("should do his stuff", func() {

		source.Set("somebody", "that I used to know")

		Expect(source.LoadEnv()).Should(HaveKeyWithValue("somebody", "that I used to know"))

	})

	It("should panic if too much Ed Sheeran", func() {

		source.Set("ed.sheeran", "Nancy Moligan")
		Expect(func() {
			source.Set("ed.sheeran", "What do I know?")
		}).Should(Panic())

	})

})

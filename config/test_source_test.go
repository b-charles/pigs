package config_test

import (
	. "github.com/b-charles/pigs/config"
	"github.com/b-charles/pigs/ioc"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Test", func() {

	It("should record entry for testing", func() {

		Test("the.dead.south", "In Hell I'll Be in Good Company")
		ioc.CallInjected(func(config Configuration) {
			Expect(config.Get("the.dead.south")).To(Equal("In Hell I'll Be in Good Company"))
		})

	})

	It("should record several entrie for testing", func() {

		TestMap(map[string]string{
			"hermans.hermits": "No Milk Today",
			"milky.chance":    "Stolen Dance",
		})
		ioc.CallInjected(func(config Configuration) {
			Expect(config.Get("hermans.hermits")).To(Equal("No Milk Today"))
			Expect(config.Get("milky.chance")).To(Equal("Stolen Dance"))
		})

	})

	It("should allow the overloading of keys", func() {

		TestMap(map[string]string{
			"black.eye.peas": "I Gotta Feeling",
		})
		Test("black.eye.peas", "The Apl song")
		ioc.CallInjected(func(config Configuration) {
			Expect(config.Get("black.eye.peas")).To(Equal("The Apl song"))
		})

	})

})

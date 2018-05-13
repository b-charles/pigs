package config_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/l3eegbee/pigs/config"
	"github.com/l3eegbee/pigs/config/confsources"
	"github.com/l3eegbee/pigs/ioc"
)

var _ = Describe("Configuration", func() {

	BeforeEach(func() {
		ioc.ClearTests()
	})

	Describe("Injection and merge", func() {

		It("Should get a simple source", func() {

			confsources.SetEnvForTests(map[string]string{
				"hello": "world",
			})

			config := ioc.GetComponent("Configuration").(*Configuration)

			Expect(config.Env).Should(HaveKeyWithValue("hello", "world"))

		})

		It("Should merge multiple sources", func() {

			confsources.SetEnvForTestsWithPriority(1, map[string]string{
				"hello": "bob",
			})
			confsources.SetEnvForTestsWithPriority(0, map[string]string{
				"hello": "world",
			})

			config := ioc.GetComponent("Configuration").(*Configuration)

			Expect(config.Env).Should(HaveKeyWithValue("hello", "bob"))

		})

	})

	Describe("Resolve placeholders", func() {

		It("Should resolve simple placeholder", func() {

			confsources.SetEnvForTests(map[string]string{
				"name":   "Batman",
				"whoami": "I'm ${name}",
			})

			config := ioc.GetComponent("Configuration").(*Configuration)

			Expect(config.Env).Should(HaveKeyWithValue("whoami", "I'm Batman"))

		})

		It("Should resolve complex placeholder", func() {

			confsources.SetEnvForTests(map[string]string{
				"egg":             "oeuf",
				"ham":             "jambon",
				"cheese":          "fromage",
				"recipe-complete": "${egg}, ${ham}, ${cheese}",
				"order":           "complete",
				"plate":           "${recipe-${order}}",
			})

			config := ioc.GetComponent("Configuration").(*Configuration)

			Expect(config.Env).Should(HaveKeyWithValue("plate", "oeuf, jambon, fromage"))

		})

	})

})

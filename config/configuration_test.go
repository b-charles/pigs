package config_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/b-charles/pigs/config"
	"github.com/b-charles/pigs/ioc"
)

func TestConfig(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Config Suite")
}

// Simple impl

type SimpleConfigSource struct {
	Priority int
	Env      map[string]string
}

func (self *SimpleConfigSource) GetPriority() int {
	return self.Priority
}

func (self *SimpleConfigSource) LoadEnv(config MutableConfig) error {
	for key, value := range self.Env {
		config.Set(key, value)
	}
	return nil
}

var _ = Describe("Configuration", func() {

	Describe("Injection and merge", func() {

		It("Should get a simple source", func() {

			SetTest(map[string]string{
				"hello": "world",
			})

			ioc.CallInjected(func(config Configuration) {
				Expect(config.Get("hello")).To(Equal("world"))
			})

		})

		It("Should merge multiple sources", func() {

			ioc.TestPut(&SimpleConfigSource{10, map[string]string{
				"hello": "bill",
			}}, func(ConfigSource) {})
			ioc.TestPut(&SimpleConfigSource{20, map[string]string{
				"hello": "bob",
			}}, func(ConfigSource) {})

			ioc.CallInjected(func(config Configuration) {
				Expect(config.Get("hello")).To(Equal("bob"))
			})

		})

	})

	Describe("Resolve placeholders", func() {

		It("Should resolve simple placeholder", func() {

			SetTest(map[string]string{
				"name":   "Batman",
				"whoami": "I'm ${name}",
			})

			ioc.CallInjected(func(config Configuration) {
				Expect(config.Get("whoami")).To(Equal("I'm Batman"))
			})

		})

		It("Should resolve complex placeholder", func() {

			SetTest(map[string]string{
				"egg":             "oeuf",
				"ham":             "jambon",
				"cheese":          "fromage",
				"recipe-complete": "${egg}, ${ham}, ${cheese}",
				"order":           "complete",
				"plate":           "${recipe-${order}}",
			})

			ioc.CallInjected(func(config Configuration) {
				Expect(config.Get("plate")).To(Equal("oeuf, jambon, fromage"))
			})

		})

	})

})

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

	Describe("Default config source", func() {

		var backup map[string]string

		BeforeEach(func() {
			ioc.TestPut("ioc test flag")
			backup = BackupDefault()
		})

		AfterEach(func() {
			RestoreDefault(backup)
		})

		It("should record default entries", func() {

			Set("somebody", "that I used to know")

			ioc.CallInjected(func(config Configuration) {
				Expect(config.Get("somebody")).To(Equal("that I used to know"))
			})

		})

		It("should panic if a default entry is defined twice", func() {

			Set("ed.sheeran", "Nancy Moligan")
			Expect(func() {
				Set("ed.sheeran", "What do I know?")
			}).Should(Panic())

		})

		It("should merge default and tests settings", func() {

			Set("daft.punk", "Something About Us")
			Test("justice", "D.A.N.C.E")

			ioc.CallInjected(func(config Configuration) {
				Expect(config.Get("daft.punk")).To(Equal("Something About Us"))
				Expect(config.Get("justice")).To(Equal("D.A.N.C.E"))
			})

		})

	})

	Describe("Injection and merge", func() {

		BeforeEach(func() {
			ioc.TestPut("ioc test flag")
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

			Test("name", "Batman")
			Test("whoami", "I'm ${name}")

			ioc.CallInjected(func(config Configuration) {
				Expect(config.Get("whoami")).To(Equal("I'm Batman"))
			})

		})

		It("Should resolve complex placeholder", func() {

			TestMap(map[string]string{
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

		It("Should not touch placeholders which are not found", func() {

			Test("whoami", "I'm ${name}")

			ioc.CallInjected(func(config Configuration) {
				Expect(config.Get("whoami")).To(Equal("I'm ${name}"))
			})

		})

	})

})

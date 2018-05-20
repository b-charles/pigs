package config_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/l3eegbee/pigs/config"
)

func TestConfig(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Config Suite")
}

type Cfg struct {
	Priority int
	Env      map[string]string
}

func NewCfg(priority int) *Cfg {
	return &Cfg{
		Priority: priority,
		Env:      make(map[string]string),
	}
}

func (self *Cfg) Put(key, value string) *Cfg {
	self.Env[key] = value
	return self
}

func (self *Cfg) GetPriority() int {
	return self.Priority
}

func (self *Cfg) LoadEnv() map[string]string {
	return self.Env
}

var _ = Describe("Configuration", func() {

	Describe("Injection and merge", func() {

		It("Should get a simple source", func() {

			cfg := NewCfg(10).Put("hello", "world")

			config := CreateConfiguration([]ConfigSource{cfg})

			Expect(config).Should(HaveKeyWithValue("hello", "world"))

		})

		It("Should merge multiple sources", func() {

			cfg1 := NewCfg(20).Put("hello", "bob")
			cfg2 := NewCfg(10).Put("hello", "world")

			config := CreateConfiguration([]ConfigSource{cfg1, cfg2})

			Expect(config).Should(HaveKeyWithValue("hello", "bob"))

		})

	})

	Describe("Resolve placeholders", func() {

		It("Should resolve simple placeholder", func() {

			cfg := NewCfg(10)
			cfg.Put("name", "Batman")
			cfg.Put("whoami", "I'm ${name}")

			config := CreateConfiguration([]ConfigSource{cfg})

			Expect(config).Should(HaveKeyWithValue("whoami", "I'm Batman"))

		})

		It("Should resolve complex placeholder", func() {

			cfg := NewCfg(10)
			cfg.Put("egg", "oeuf")
			cfg.Put("ham", "jambon")
			cfg.Put("cheese", "fromage")
			cfg.Put("recipe-complete", "${egg}, ${ham}, ${cheese}")
			cfg.Put("order", "complete")
			cfg.Put("plate", "${recipe-${order}}")

			config := CreateConfiguration([]ConfigSource{cfg})

			Expect(config).Should(HaveKeyWithValue("plate", "oeuf, jambon, fromage"))

		})

	})

})

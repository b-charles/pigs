package config_test

import (
	. "github.com/b-charles/pigs/config"
	"github.com/b-charles/pigs/ioc"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Default config source", func() {

	BeforeEach(func() {
		ioc.TestPut("test flag")
	})

	It("should record default entries", func() {

		SetDefault("somebody", "that I used to know")

		ioc.CallInjected(func(config DefaultConfigSource) {
			Expect(config["somebody"]).To(Equal("that I used to know"))
		})

	})

	It("should panic if too much Ed Sheeran", func() {

		SetDefault("ed.sheeran", "Nancy Moligan")
		Expect(func() {
			SetDefault("ed.sheeran", "What do I know?")
		}).Should(Panic())

	})

	It("should merge default and tests settings", func() {

		SetDefault("daft.punk", "Something About Us")
		SetTest(map[string]string{
			"justice": "D.A.N.C.E",
		})

		ioc.CallInjected(func(config Configuration) {
			Expect(config.Get("daft.punk")).To(Equal("Something About Us"))
			Expect(config.Get("justice")).To(Equal("D.A.N.C.E"))
		})

	})

})

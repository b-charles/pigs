package config_test

import (
	. "github.com/b-charles/pigs/config"
	"github.com/b-charles/pigs/ioc"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/spf13/afero"
)

var _ = Describe("Json", func() {

	var backup map[string]string

	BeforeEach(func() {
		backup = BackupDefault()
	})

	AfterEach(func() {
		RestoreDefault(backup)
	})

	It("should parse simple value", func() {

		appFs := afero.NewMemMapFs()
		afero.WriteFile(appFs, "application.json", []byte("{\"hello\":\"world\"}"), 0644)

		ioc.TestPut(appFs, func(afero.Fs) {})

		ioc.CallInjected(func(config Configuration) {
			Expect(config.Get("hello")).To(Equal("world"))
		})

	})

	It("should load the different files", func() {

		Set("config.json.01", "file1.json")
		Set("config.json.02", "file2.json")

		appFs := afero.NewMemMapFs()
		afero.WriteFile(appFs, "file1.json", []byte("{\"first\":\"james\"}"), 0644)
		afero.WriteFile(appFs, "file2.json", []byte("{\"last\":\"bond\"}"), 0644)

		ioc.TestPut(appFs, func(afero.Fs) {})

		ioc.CallInjected(func(config Configuration) {
			Expect(config.Get("first")).To(Equal("james"))
			Expect(config.Get("last")).To(Equal("bond"))
		})

	})

})

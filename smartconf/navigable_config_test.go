package smartconf_test

import (
	"fmt"

	"github.com/b-charles/pigs/config"
	"github.com/b-charles/pigs/ioc"
	. "github.com/b-charles/pigs/smartconf"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("IOC registration", func() {

	It("should parse empty config", func() {

		config.SetTest(map[string]string{})

		ioc.CallInjected(func(config *NavConfig, status *ioc.ContainerStatus) {
			fmt.Print(status)
			Expect(config.Keys()).To(BeEmpty())
		})

	})

})

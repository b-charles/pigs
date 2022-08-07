package smartconf_test

import (
	"github.com/b-charles/pigs/config"
	"github.com/b-charles/pigs/ioc"
	. "github.com/b-charles/pigs/smartconf"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type empty_config struct{}

type simple_config struct {
	Property string
}

func alert_parser(value string) (int, error) {
	if value == "apocalypse" {
		return 4, nil
	} else {
		return 0, nil
	}
}

type parsed_config struct {
	Level int
}

type complex_root_config struct {
	Param1 string
	Sub    complex_nested_config
}

type complex_nested_config struct {
	List []int
	Map  map[string]bool
}

var _ = Describe("Smart configuration", func() {

	It("should accept empty config", func() {

		config := &empty_config{}
		TestConfigure("", config)

		ioc.CallInjected(func(injected *empty_config) {
			Expect(injected).To(Equal(config))
		})

	})

	It("should accept simple string config", func() {

		TestConfigure("my", &simple_config{})
		config.SetTest(map[string]string{
			"my.property": "Hello, World!",
		})

		ioc.CallInjected(func(injected *simple_config) {
			Expect(injected).To(Equal(&simple_config{
				Property: "Hello, World!",
			}))
		})

	})

	It("should use defined parser", func() {

		ioc.TestPut(alert_parser, func(Parser) {})
		TestConfigure("threat", &parsed_config{})
		config.SetTest(map[string]string{
			"threat.level": "apocalypse",
		})

		ioc.CallInjected(func(injected *parsed_config) {
			Expect(injected).To(Equal(&parsed_config{
				Level: 4,
			}))
		})

	})

	It("should configure a complex struct", func() {

		TestConfigure("", &complex_root_config{})
		config.SetTest(map[string]string{
			"param1":      "great value",
			"sub.list.2":  "9",
			"sub.list.12": "8",
			"sub.list.a":  "7",
			"sub.list.b":  "6",
			"sub.list.c":  "5",
			"sub.map.yes": "true",
			"sub.map.no":  "false",
		})

		ioc.CallInjected(func(injected *complex_root_config) {
			Expect(injected).To(Equal(&complex_root_config{
				Param1: "great value",
				Sub: complex_nested_config{
					List: []int{9, 8, 7, 6, 5},
					Map: map[string]bool{
						"yes": true,
						"no":  false,
					},
				},
			}))
		})

	})

})

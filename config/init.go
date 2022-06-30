package config

import (
	_ "github.com/b-charles/pigs/config/confsources/args"
	_ "github.com/b-charles/pigs/config/confsources/envvar"
	_ "github.com/b-charles/pigs/config/confsources/programmatic"
	"github.com/b-charles/pigs/ioc"
)

func init() {

	ioc.PutFactory(
		func(injected struct {
			ConfigSources []ConfigSource
		}) Configuration {
			return CreateConfiguration(injected.ConfigSources)
		}, "Configuration")

}

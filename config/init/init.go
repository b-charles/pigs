package init

import (
	. "github.com/b-charles/pigs/config"
	"github.com/b-charles/pigs/ioc"
)

func init() {

	ioc.PutFactory(
		func(injected struct {
			DefaultSources []ConfigSource
			ConfigSources  []ConfigSource
		}) *Configuration {
			return CreateConfiguration(append(injected.DefaultSources, injected.ConfigSources...))
		}, "Configuration")

}

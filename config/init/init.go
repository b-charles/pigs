package init

import (
	. "github.com/l3eegbee/pigs/config"
	"github.com/l3eegbee/pigs/ioc"
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

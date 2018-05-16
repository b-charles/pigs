package init

import (
	. "github.com/l3eegbee/pigs/config"
	"github.com/l3eegbee/pigs/ioc"
)

func init() {

	ioc.PutFactory(
		func(defaultSources []ConfigSource, sources []ConfigSource) *Configuration {
			return CreateConfiguration(append(defaultSources, sources...))
		}, []string{"DefaultConfigSources", "ConfigSources"}, "Configuration")

}

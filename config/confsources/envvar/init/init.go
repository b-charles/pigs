package init

import (
	. "github.com/l3eegbee/pigs/config/confsources/envvar"
	"github.com/l3eegbee/pigs/ioc"
)

func init() {
	ioc.Put(NewEnvVarConfigSource(), "EnvVarConfigSource", "EnvVar", "ConfigSources")
}

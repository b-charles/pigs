package init

import (
	. "github.com/b-charles/pigs/config/confsources/envvar"
	"github.com/b-charles/pigs/ioc"
)

func init() {
	ioc.Put(NewEnvVarConfigSource(), "EnvVarConfigSource", "EnvVar", "ConfigSources")
}

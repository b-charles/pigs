package envvar

import (
	"github.com/b-charles/pigs/ioc"
)

func init() {
	ioc.Put(NewEnvVarConfigSource(), "EnvVarConfigSource", "ConfigSource")
}

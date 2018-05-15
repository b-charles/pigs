package envvar

import (
	"os"
	"strings"

	. "github.com/l3eegbee/pigs/config/confsources"
	"github.com/l3eegbee/pigs/ioc"
)

func NewEnvVarConfigSource() *SimpleConfigSource {

	env := make(map[string]string)

	for _, e := range os.Environ() {
		pair := strings.Split(e, "=")
		env[ConvertEnvVarKey(pair[0])] = pair[1]
	}

	return &SimpleConfigSource{
		Priority: CONFIG_SOURCE_PRIORITY_ENV_VAR,
		Env:      env,
	}

}

func init() {
	ioc.Put(NewEnvVarConfigSource(), "EnvVarConfigSource", "EnvVar", "ConfigSources")
}

package confsources

import (
	"os"
	"strings"

	"github.com/l3eegbee/pigs/config"
	"github.com/l3eegbee/pigs/ioc"
)

func ConvertEnvVarKey(key string) string {
	return strings.Replace(strings.ToLower(key), "_", ".", -1)
}

func NewEnvVarConfigSource() *config.SimpleConfigSource {

	env := make(map[string]string)

	for _, e := range os.Environ() {
		pair := strings.Split(e, "=")
		env[ConvertEnvVarKey(pair[0])] = pair[1]
	}

	return &config.SimpleConfigSource{
		Priority: CONFIG_SOURCE_PRIORITY_ENV_VAR,
		Env:      env,
	}

}

func init() {
	ioc.Put(NewEnvVarConfigSource(), "EnvVarConfigSource", "ConfigSources")
}

package config

import (
	"os"
	"strings"

	"github.com/b-charles/pigs/ioc"
)

func convertEnvVarKey(key string) string {
	return strings.Replace(strings.ToLower(key), "_", ".", -1)
}

func ParseEnvVar(envvar []string) map[string]string {

	env := make(map[string]string)

	for _, e := range envvar {
		pair := strings.Split(e, "=")
		env[convertEnvVarKey(pair[0])] = pair[1]
	}

	return env

}

type EnvVarConfigSource map[string]string

func (self EnvVarConfigSource) GetPriority() int {
	return CONFIG_SOURCE_PRIORITY_ENV_VAR
}

func (self EnvVarConfigSource) LoadEnv(config MutableConfig) error {
	for k, v := range self {
		config.Set(k, v)
	}
	return nil
}

func (self EnvVarConfigSource) String() string {
	return stringify(self)
}

func init() {
	ioc.Put(EnvVarConfigSource(ParseEnvVar(os.Environ())), func(ConfigSource) {})
}

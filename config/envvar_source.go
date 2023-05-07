package config

import (
	"os"
	"strings"

	"github.com/b-charles/pigs/ioc"
	"github.com/b-charles/pigs/json"
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

type EnvVarConfigSource struct {
	source map[string]string
}

func (self *EnvVarConfigSource) GetPriority() int {
	return CONFIG_SOURCE_PRIORITY_ENV_VAR
}

func (self *EnvVarConfigSource) LoadEnv(config MutableConfig) error {
	for k, v := range self.source {
		config.Set(k, v)
	}
	return nil
}

func (self *EnvVarConfigSource) Json() json.JsonNode {
	return json.NewJsonObjectStrings(self.source)
}

func (self *EnvVarConfigSource) String() string {
	return self.Json().String()
}

func init() {

	ioc.PutNamed("Env var config source",
		&EnvVarConfigSource{ParseEnvVar(os.Environ())},
		func(ConfigSource) {})

}

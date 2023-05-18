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

	env := make(map[string]string, len(envvar))

	for _, e := range envvar {
		pair := strings.Split(e, "=")
		env[convertEnvVarKey(pair[0])] = pair[1]
	}

	return env

}

type EnvVarConfigSource ConfigSource

var CONFIG_SOURCE_PRIORITY_ENV_VAR = 0

type EnvVarConfigSourceImpl struct {
	source map[string]string
}

func (self *EnvVarConfigSourceImpl) GetPriority() int {
	return CONFIG_SOURCE_PRIORITY_ENV_VAR
}

func (self *EnvVarConfigSourceImpl) LoadEnv(config MutableConfig) error {
	for k, v := range self.source {
		config.Set(k, v)
	}
	return nil
}

func (self *EnvVarConfigSourceImpl) Json() json.JsonNode {
	return json.NewJsonObjectStrings(self.source)
}

func (self *EnvVarConfigSourceImpl) String() string {
	return self.Json().String()
}

func init() {

	ioc.DefaultPutNamed("Env var config source (default)",
		&EnvVarConfigSourceImpl{ParseEnvVar(os.Environ())},
		func(EnvVarConfigSource) {})

	ioc.PutNamedFactory("Env var config source (promoter)",
		func(v EnvVarConfigSource) (ConfigSource, error) { return v, nil })

}

package envvar

import (
	"os"
	"strings"

	. "github.com/l3eegbee/pigs/config/confsources"
)

func ParseEnvVar(envvar []string) map[string]string {

	env := make(map[string]string)

	for _, e := range envvar {
		pair := strings.Split(e, "=")
		env[ConvertEnvVarKey(pair[0])] = pair[1]
	}

	return env

}

func NewEnvVarConfigSource() *SimpleConfigSource {

	env := ParseEnvVar(os.Environ())

	return &SimpleConfigSource{
		Priority: CONFIG_SOURCE_PRIORITY_ENV_VAR,
		Env:      env,
	}

}

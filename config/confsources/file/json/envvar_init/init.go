package envvar_init

import (
	. "github.com/b-charles/pigs/config/confsources"
	. "github.com/b-charles/pigs/config/confsources/file"
	. "github.com/b-charles/pigs/config/confsources/file/json"
)

func init() {

	RegisterFormatedEnvVarConfig(
		CONFIG_SOURCE_PRIORITY_ENV_VAR_JSON,
		"_APPLICATION_JSON",
		ParseJsonToEnv,
		"JsonEnvVarConfigSource")

}

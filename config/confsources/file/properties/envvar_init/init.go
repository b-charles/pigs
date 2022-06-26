package envvar_init

import (
	. "github.com/b-charles/pigs/config/confsources"
	. "github.com/b-charles/pigs/config/confsources/file"
	. "github.com/b-charles/pigs/config/confsources/file/properties"
)

func init() {

	RegisterFormatedEnvVarConfig(
		CONFIG_SOURCE_PRIORITY_ENV_VAR_PROPERTIES,
		"_APPLICATION_PROPERTIES",
		ParsePropertiesToEnv,
		"PropertiesEnvVarConfigSource")

}

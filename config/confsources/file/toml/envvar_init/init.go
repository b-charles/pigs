package envvar_init

import (
	. "github.com/l3eegbee/pigs/config/confsources"
	. "github.com/l3eegbee/pigs/config/confsources/file"
	. "github.com/l3eegbee/pigs/config/confsources/file/toml"
)

func init() {

	RegisterFormatedEnvVarConfig(
		CONFIG_SOURCE_PRIORITY_ENV_VAR_TOML,
		"_APPLICATION_TOML",
		ParseTomlToEnv,
		"TomlEnvVarConfigSource")

}

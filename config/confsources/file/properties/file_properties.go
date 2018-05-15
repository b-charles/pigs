package properties

import (
	enc "github.com/magiconair/properties"

	. "github.com/l3eegbee/pigs/config/confsources"
	. "github.com/l3eegbee/pigs/config/confsources/file"
)

func ParsePropertiesToEnv(content string) map[string]string {

	p := enc.MustLoadString(content)
	p.Prefix = ""
	p.Postfix = ""

	return p.Map()

}

func init() {

	RegisterFileConfig(
		CONFIG_SOURCE_PRIORITY_FILE_PROPERTIES,
		".properties",
		ParsePropertiesToEnv,
		"PropertiesFileConfigSource")

	RegisterFormatedEnvVarConfig(
		CONFIG_SOURCE_PRIORITY_ENV_VAR_PROPERTIES,
		"_APPLICATION_PROPERTIES",
		ParsePropertiesToEnv,
		"PropertiesEnvVarConfigSource")

}

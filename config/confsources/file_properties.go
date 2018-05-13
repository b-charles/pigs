package confsources

import (
	"github.com/magiconair/properties"

	_ "github.com/l3eegbee/pigs/filesystem"
)

func ParsePropertiesToEnv(content string) map[string]string {

	p := properties.MustLoadString(content)
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

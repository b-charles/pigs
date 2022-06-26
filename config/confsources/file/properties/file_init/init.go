package file_init

import (
	. "github.com/b-charles/pigs/config/confsources"
	. "github.com/b-charles/pigs/config/confsources/file"
	. "github.com/b-charles/pigs/config/confsources/file/properties"
)

func init() {

	RegisterFileConfig(
		CONFIG_SOURCE_PRIORITY_FILE_PROPERTIES,
		".properties",
		ParsePropertiesToEnv,
		"PropertiesFileConfigSource")

}

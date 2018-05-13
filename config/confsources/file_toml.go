package confsources

import "github.com/BurntSushi/toml"

func ParseTomlToEnv(content string) map[string]string {

	root := make(map[string]interface{})
	if err := toml.Unmarshal([]byte(content), &root); err != nil {
		panic(err)
	}

	env := make(map[string]string)
	convertObjectInEnv(env, "", root)

	return env

}

func init() {

	RegisterFileConfig(
		CONFIG_SOURCE_PRIORITY_FILE_TOML,
		".toml",
		ParseTomlToEnv,
		"TomlFileConfigSource")

	RegisterFormatedEnvVarConfig(
		CONFIG_SOURCE_PRIORITY_ENV_VAR_TOML,
		"_APPLICATION_TOML",
		ParseTomlToEnv,
		"TomlEnvVarConfigSource")

}

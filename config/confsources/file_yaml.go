package confsources

import (
	"gopkg.in/yaml.v2"
)

func ParseYamlToEnv(content string) map[string]string {

	root := make(map[string]interface{})
	if err := yaml.Unmarshal([]byte(content), &root); err != nil {
		panic(err)
	}

	env := make(map[string]string)
	convertObjectInEnv(env, "", root)

	return env

}

func init() {

	RegisterFileConfig(
		CONFIG_SOURCE_PRIORITY_FILE_YAML,
		".yaml",
		ParseYamlToEnv,
		"YamlFileConfigSource")

	RegisterFormatedEnvVarConfig(
		CONFIG_SOURCE_PRIORITY_ENV_VAR_YAML,
		"_APPLICATION_YAML",
		ParseYamlToEnv,
		"YamlEnvVarConfigSource")

}

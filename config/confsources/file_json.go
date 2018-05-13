package confsources

import (
	"encoding/json"
)

func ParseJsonToEnv(content string) map[string]string {

	root := make(map[string]interface{})
	if err := json.Unmarshal([]byte(content), &root); err != nil {
		panic(err)
	}

	env := make(map[string]string)
	convertObjectInEnv(env, "", root)

	return env

}

func init() {

	RegisterFileConfig(
		CONFIG_SOURCE_PRIORITY_FILE_JSON,
		".json",
		ParseJsonToEnv,
		"JsonFileConfigSource")

	RegisterFormatedEnvVarConfig(
		CONFIG_SOURCE_PRIORITY_ENV_VAR_JSON,
		"_APPLICATION_JSON",
		ParseJsonToEnv,
		"JsonEnvVarConfigSource")

}

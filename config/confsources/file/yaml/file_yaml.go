package yaml

import (
	enc "gopkg.in/yaml.v2"

	. "github.com/l3eegbee/pigs/config/confsources/file"
)

func ParseYamlToEnv(content string) map[string]string {

	root := make(map[string]interface{})
	if err := enc.Unmarshal([]byte(content), &root); err != nil {
		panic(err)
	}

	env := make(map[string]string)
	ConvertObjectInEnv(env, "", root)

	return env

}

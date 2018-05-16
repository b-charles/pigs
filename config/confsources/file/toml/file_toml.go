package toml

import (
	enc "github.com/BurntSushi/toml"

	. "github.com/l3eegbee/pigs/config/confsources/file"
)

func ParseTomlToEnv(content string) map[string]string {

	root := make(map[string]interface{})
	if err := enc.Unmarshal([]byte(content), &root); err != nil {
		panic(err)
	}

	env := make(map[string]string)
	ConvertObjectInEnv(env, "", root)

	return env

}

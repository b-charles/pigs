package json

import (
	enc "encoding/json"

	. "github.com/b-charles/pigs/config/confsources/file"
)

func ParseJsonToEnv(content string) map[string]string {

	root := make(map[string]interface{})
	if err := enc.Unmarshal([]byte(content), &root); err != nil {
		panic(err)
	}

	env := make(map[string]string)
	ConvertObjectInEnv(env, "", root)

	return env

}

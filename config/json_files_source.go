package config

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/b-charles/pigs/ioc"
	"github.com/b-charles/pigs/json"
	"github.com/spf13/afero"
)

func mergeIn(config MutableConfig, path string, json json.JsonNode) {

	if json.IsString() {

		config.Set(path, json.AsString())

	} else if json.IsFloat() {

		config.Set(path, fmt.Sprintf("%f", json.AsFloat()))

	} else if json.IsInt() {

		config.Set(path, fmt.Sprintf("%d", json.AsInt()))

	} else if json.IsBool() {

		if json.AsBool() {
			config.Set(path, "true")
		} else {
			config.Set(path, "false")
		}

	} else if json.IsObject() {

		for _, key := range json.GetKeys() {

			sub := key
			if path != "" {
				sub = fmt.Sprintf("%s.%s", path, key)
			}

			mergeIn(config, sub, json.GetMember(key))

		}

	} else if json.IsArray() {

		l := json.GetLen()
		for i := 0; i < l; i++ {

			sub := ""
			if path != "" {
				sub = fmt.Sprintf("%s.%d", path, i)
			} else {
				sub = fmt.Sprintf("%d", i)
			}

			mergeIn(config, sub, json.GetElement(i))

		}

	} else { // null

		config.Set(path, "null")

	}

}

var CONFIG_SOURCE_JSON_PREFIX = "config.json"

type JsonFilesConfigSource struct {
	fs afero.Fs
}

func (self *JsonFilesConfigSource) GetPriority() int {
	return CONFIG_SOURCE_PRIORITY_JSON_FILES
}

func (self *JsonFilesConfigSource) LoadEnv(config MutableConfig) error {

	keys := []string{}
	for _, key := range config.Keys() {
		if strings.HasPrefix(key, CONFIG_SOURCE_JSON_PREFIX) {
			keys = append(keys, key)
		}
	}

	for _, key := range keys {

		if path, _, err := config.Lookup(key); err != nil {

			return err

		} else {

			if b, err := afero.ReadFile(self.fs, path); err != nil {
				continue
			} else if json, err := json.Parse(bytes.NewReader(b)); err != nil {
				return err
			} else {
				mergeIn(config, "", json)
			}

		}

	}

	return nil

}

func init() {

	Set(CONFIG_SOURCE_JSON_PREFIX, "application.json")

	ioc.PutNamedFactory("Json config source",
		func(fs afero.Fs) (*JsonFilesConfigSource, error) {
			return &JsonFilesConfigSource{fs}, nil
		}, func(ConfigSource) {})

}

package config

import (
	"fmt"
	"sync"

	"github.com/b-charles/pigs/ioc"
)

/*
 * Defaults
 */

type DefaultConfigSource map[string]string

func (self DefaultConfigSource) GetPriority() int {
	return CONFIG_SOURCE_PRIORITY_DEFAULT
}

func (self DefaultConfigSource) LoadEnv(config MutableConfig) error {
	for k, v := range self {
		config.Set(k, v)
	}
	return nil
}

func (self DefaultConfigSource) Set(key, value string) {

	if oldValue, present := self[key]; present {
		panic(fmt.Sprintf("The default value '%s' can't be overwrited from '%s' to '%s'.", key, oldValue, value))
	}

	self[key] = value

}

func (self DefaultConfigSource) String() string {
	return stringify(self)
}

var defaultConfigSourceInstance DefaultConfigSource
var onceDefaultConfigSource sync.Once

func DefaultConfigSourceInstance() DefaultConfigSource {
	onceDefaultConfigSource.Do(func() {
		defaultConfigSourceInstance = map[string]string{}
	})
	return defaultConfigSourceInstance
}

func SetDefault(key, value string) {
	DefaultConfigSourceInstance().Set(key, value)
}

func init() {
	ioc.PutFactory(DefaultConfigSourceInstance, func(ConfigSource) {})
}

/*
 * Tests
 */

type TestConfigSource struct {
	Base DefaultConfigSource `inject:""`
	env  map[string]string
}

func (self *TestConfigSource) GetPriority() int {
	return CONFIG_SOURCE_PRIORITY_TESTS
}

func (self *TestConfigSource) LoadEnv(config MutableConfig) error {

	for k, v := range self.Base {
		if !config.HasKey(k) {
			config.Set(k, v)
		}
	}

	for k, v := range self.env {
		config.Set(k, v)
	}

	return nil

}

func (self *TestConfigSource) Set(key, value string) {
	self.env[key] = value
}

func (self *TestConfigSource) String() string {
	return stringify(self.env)
}

func SetTest(values map[string]string) {

	config := &TestConfigSource{
		env: map[string]string{},
	}

	for k, v := range values {
		config.Set(k, v)
	}

	ioc.TestPut(config, func(ConfigSource) {})

}

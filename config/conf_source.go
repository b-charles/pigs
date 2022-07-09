package config

import (
	"fmt"
	"sync"

	"github.com/b-charles/pigs/ioc"
)

/*
 * Defaults
 */

type DefaultConfigSource struct {
	*SimpleConfigSource
}

func NewDefaultConfigSource() *DefaultConfigSource {
	return &DefaultConfigSource{
		&SimpleConfigSource{
			Priority: CONFIG_SOURCE_PRIORITY_DEFAULT,
			Env:      make(map[string]string),
		}}
}

func (self *DefaultConfigSource) Set(key, value string) {

	if oldValue, ok := self.Env[key]; ok {
		panic(fmt.Sprintf("The default value '%s' can't be overwrited from '%s' to '%s'.", key, oldValue, value))
	}

	self.Env[key] = value

}

var defaultConfigSourceInstance *DefaultConfigSource
var once sync.Once

func DefaultConfigSourceInstance() *DefaultConfigSource {
	once.Do(func() {
		defaultConfigSourceInstance = NewDefaultConfigSource()
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
	*SimpleConfigSource
}

func NewTestConfigSource(base *DefaultConfigSource, priority int, env map[string]string) *TestConfigSource {

	res := make(map[string]string)
	for k, v := range base.Env {
		res[k] = v
	}
	for k, v := range env {
		res[k] = v
	}

	return &TestConfigSource{
		&SimpleConfigSource{
			Priority: priority,
			Env:      res,
		}}

}

func SetEnvForTestsWithPriority(priority int, env map[string]string) {

	ioc.TestPutFactory(func(injected struct {
		DefaultConfigSource *DefaultConfigSource
	}) *TestConfigSource {
		return NewTestConfigSource(injected.DefaultConfigSource, priority, env)
	}, func(ConfigSource) {})

}

func SetEnvForTests(env map[string]string) {
	SetEnvForTestsWithPriority(CONFIG_SOURCE_PRIORITY_TESTS, env)
}

package config

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
	"sync"

	"github.com/b-charles/pigs/ioc"
	"github.com/b-charles/pigs/json"
)

/*
 * ConfigSource
 */

type ConfigSource interface {
	GetPriority() int
	LoadEnv(MutableConfig) error
}

/*
 * Mutable config
 */

type MutableConfig interface {
	HasKey(string) bool
	Keys() []string
	GetRaw(string) (string, bool)
	Lookup(string) (string, bool, error)
	Get(string) string
	Set(string, string)
}

type Configuration interface {
	HasKey(string) bool
	Keys() []string
	GetRaw(string) (string, bool)
	Lookup(string) (string, bool, error)
	Get(string) string
}

type configImpl struct {
	mutable  bool
	raws     map[string]string
	resolved map[string]string
}

func (self *configImpl) HasKey(key string) bool {
	_, p := self.raws[key]
	return p
}

func (self *configImpl) Keys() []string {
	keys := make([]string, 0, len(self.raws))
	for k := range self.raws {
		keys = append(keys, k)
	}
	return keys
}

func (self *configImpl) GetRaw(key string) (string, bool) {
	value, p := self.raws[key]
	return value, p
}

var placeholderRegexp *regexp.Regexp = regexp.MustCompile("\\${\\s*([^{}]+)\\s*}")

func (self *configImpl) resolveValue(key string, traveled map[string]bool) (string, bool, *CyclicLoopError) {

	if resolved, p := self.resolved[key]; p {
		return resolved, true, nil
	}
	if value, p := self.raws[key]; !p {
		return "", false, nil
	} else {

		if t, p := traveled[key]; p && t {
			return value, true, newCyclicLoopError(key, value)
		}
		traveled[key] = true

		replace := true
		for replace {

			replace = false
			if matches := placeholderRegexp.FindAllStringSubmatch(value, -1); matches != nil {
				for _, match := range matches {
					if subvalue, found, err := self.resolveValue(match[1], traveled); err != nil {
						return "", true, err.push(key, value)
					} else if found {
						value = strings.Replace(value, match[0], subvalue, -1)
						replace = true
						break
					}
				}
			}

		}

		traveled[key] = false

		if !self.mutable {
			self.resolved[key] = value
		}

		return value, true, nil

	}

}

func (self *configImpl) Lookup(key string) (string, bool, error) {
	if value, found, err := self.resolveValue(key, map[string]bool{}); err == (*CyclicLoopError)(nil) {
		return value, found, nil
	} else {
		return value, found, error(err)
	}
}

func (self *configImpl) Get(key string) string {
	if v, _, err := self.Lookup(key); err != nil {
		panic(err)
	} else {
		return v
	}
}

func (self *configImpl) Set(key, value string) {
	self.raws[key] = value
}

func (self *configImpl) Json() json.JsonNode {
	return json.NewJsonObjectStrings(self.raws)
}

func (self *configImpl) String() string {
	return self.Json().String()
}

/*
 * Cyclic loop error
 */

type CyclicLoopElement struct {
	key   string
	value string
}

type CyclicLoopError struct {
	loop []CyclicLoopElement
}

func newCyclicLoopError(key, value string) *CyclicLoopError {
	err := &CyclicLoopError{make([]CyclicLoopElement, 0)}
	return err.push(key, value)
}

func (self *CyclicLoopError) push(key, value string) *CyclicLoopError {
	self.loop = append(self.loop, CyclicLoopElement{key, value})
	return self
}

func (self *CyclicLoopError) Error() string {

	if self == nil {
		return "Nil cyclic loop error. Should not be raised."
	}

	var b strings.Builder

	b.WriteString("Cyclic loop detected: ")

	b.WriteString(self.loop[0].key)
	b.WriteString(":'")
	b.WriteString(self.loop[0].value)
	for i := 1; i < len(self.loop); i++ {
		b.WriteString("' -> ")
		b.WriteString(self.loop[i].key)
		b.WriteString(":'")
		b.WriteString(self.loop[i].value)
	}
	b.WriteString("'")

	return b.String()

}

/*
 * Default
 */

var (
	defaultConfigMap     map[string]string
	onceDefaultConfigMap sync.Once
)

func getDefaultConfigMap() map[string]string {

	onceDefaultConfigMap.Do(func() {
		defaultConfigMap = map[string]string{}
	})

	return defaultConfigMap

}

func Set(key, value string) {
	SetMap(map[string]string{key: value})
}

func SetMap(values map[string]string) {

	config := getDefaultConfigMap()

	for key, value := range values {
		if oldValue, present := config[key]; present {
			panic(fmt.Sprintf("The default value '%s' can't be overwrited from '%s' to '%s'.", key, oldValue, value))
		}
		config[key] = value
	}

}

func BackupDefault() map[string]string {
	config := getDefaultConfigMap()
	backup := make(map[string]string, len(config))
	for k, v := range config {
		backup[k] = v
	}
	return backup
}

func RestoreDefault(backup map[string]string) {
	config := getDefaultConfigMap()
	for k := range config {
		delete(config, k)
	}
	for k, v := range backup {
		config[k] = v
	}
}

/*
 * Factory
 */

func CreateConfiguration(sources []ConfigSource) (Configuration, error) {

	sort.Slice(sources, func(i, j int) bool {
		return sources[i].GetPriority() < sources[j].GetPriority()
	})

	conf := &configImpl{
		mutable:  true,
		raws:     make(map[string]string),
		resolved: make(map[string]string),
	}

	for k, v := range getDefaultConfigMap() {
		conf.Set(k, v)
	}
	for _, source := range sources {
		err := source.LoadEnv(conf)
		if err != nil {
			return nil, fmt.Errorf("Error during loading configuration from '%v': %w", source, err)
		}
	}

	conf.mutable = false

	return conf, nil

}

func init() {

	ioc.DefaultPutNamedFactory("Configuration", CreateConfiguration)

}

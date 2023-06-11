package config

import (
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strings"
	"sync"

	"github.com/b-charles/pigs/ioc"
	"github.com/b-charles/pigs/json"
	"github.com/b-charles/pigs/memfun"
)

/*
 * ConfigSource
 */

type ConfigSource interface {
	GetPriority() int
	LoadEnv(MutableConfig) error
}

/*
 * Config
 */

var placeholderRegexp *regexp.Regexp = regexp.MustCompile("\\${\\s*([^{}]+)\\s*}")

type pstring struct {
	p   bool
	str string
}

func resolveValue(raws sync.Map, key string, recfun func(string) (pstring, error)) (pstring, error) {

	if uncastedValue, p := raws.Load(key); !p {

		return pstring{false, ""}, nil

	} else {

		value := uncastedValue.(string)

		replace := true
		for replace {

			replace = false
			if matches := placeholderRegexp.FindAllStringSubmatch(value, -1); matches != nil {
				for _, match := range matches {

					if subpvalue, err := recfun(match[1]); err != nil {
						return subpvalue, err
					} else if subpvalue.p {
						value = strings.Replace(value, match[0], subpvalue.str, -1)
						replace = true
						break
					}

				}
			}

		}

		return pstring{true, value}, nil

	}

}

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
	raws     sync.Map
	resolved memfun.MemFun[string, pstring]
}

func newConfigImpl() *configImpl {

	config := new(configImpl)
	config.mutable = false
	config.resolved = memfun.NewMemFun(func(key string, recfun func(string) (pstring, error)) (pstring, error) {
		return resolveValue(config.raws, key, recfun)
	})

	return config

}

func (self *configImpl) HasKey(key string) bool {
	_, p := self.raws.Load(key)
	return p
}

func (self *configImpl) Keys() []string {
	keys := make([]string, 0)
	self.raws.Range(func(k, v any) bool {
		keys = append(keys, k.(string))
		return true
	})
	return keys
}

func (self *configImpl) GetRaw(key string) (string, bool) {
	value, p := self.raws.Load(key)
	return value.(string), p
}

func (self *configImpl) Lookup(key string) (string, bool, error) {

	var (
		result pstring
		err    error
	)

	if self.mutable {

		called := map[string]bool{key: true}

		var recfun func(k string) (pstring, error)
		recfun = func(k string) (pstring, error) {

			if _, p := called[k]; p {
				return pstring{false, ""}, memfun.CyclicLoopError[string]{
					Stack: []string{k},
				}
			}

			called[key] = true
			defer delete(called, key)

			r, e := resolveValue(self.raws, key, recfun)

			if e != nil {
				if cyclic, ok := e.(memfun.CyclicLoopError[string]); ok {
					return r, cyclic.Append(key)
				}
			}

			return r, e

		}

		result, err = resolveValue(self.raws, key, recfun)

	} else {

		result, err = self.resolved.Get(key)

	}

	if err != nil {
		if cyclic, ok := err.(memfun.CyclicLoopError[string]); ok {

			var b strings.Builder

			b.WriteString("Cyclic loop detected: ")

			fmtElt := func(k string) string {
				v, _ := self.raws.Load(k)
				return fmt.Sprintf("%s: '%v'", k, v)
			}

			b.WriteString(fmtElt(cyclic.Stack[0]))
			for i := 1; i < len(cyclic.Stack); i++ {
				b.WriteString(" -> ")
				b.WriteString(fmtElt(cyclic.Stack[i]))
			}

			return "", false, errors.New(b.String())

		}
	}

	return result.str, result.p, err

}

func (self *configImpl) Get(key string) string {
	if v, _, err := self.Lookup(key); err != nil {
		panic(err)
	} else {
		return v
	}
}

func (self *configImpl) Set(key, value string) {
	self.raws.Store(key, value)
}

func (self *configImpl) Json() json.JsonNode {

	r := make(map[string]string)
	self.raws.Range(func(key, value any) bool {
		r[key.(string)] = value.(string)
		return true
	})

	return json.NewJsonObjectStrings(r)

}

func (self *configImpl) String() string {
	return self.Json().String()
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

	conf := newConfigImpl()

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

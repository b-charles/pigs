package config

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/b-charles/pigs/ioc"
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

func keys(m map[string]string) []string {

	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}

	return keys

}

type MutableConfig map[string]string

func (self MutableConfig) HasKey(key string) bool {
	_, present := self[key]
	return present
}

func (self MutableConfig) Keys() []string {
	return keys(self)
}

func (self MutableConfig) Get(key string) (string, error) {
	return resolveValue(self, key, map[string]bool{}, false)
}

func (self MutableConfig) Lookup(key string) (string, bool, error) {
	_, p := self[key]
	if p {
		v, e := resolveValue(self, key, map[string]bool{}, false)
		return v, true, e
	} else {
		return "", false, nil
	}
}

func (self MutableConfig) Set(key, value string) {
	self[key] = value
}

func (self MutableConfig) String() string {
	return stringify(self)
}

/*
 * Configuration
 */

type Configuration map[string]string

func (self *Configuration) Keys() []string {
	return keys(*self)
}

func (self *Configuration) Get(key string) string {
	return (*self)[key]
}

func (self *Configuration) Lookup(key string) (string, bool) {
	v, p := (*self)[key]
	return v, p
}

func (self Configuration) String() string {
	return stringify(self)
}

/*
 * Factory
 */

func CreateConfiguration(sources []ConfigSource) (Configuration, error) {

	sort.Slice(sources, func(i, j int) bool {
		return sources[i].GetPriority() < sources[j].GetPriority()
	})

	mutable := MutableConfig(map[string]string{})
	for _, source := range sources {
		err := source.LoadEnv(mutable)
		if err != nil {
			return nil, fmt.Errorf("Error during loading configuration from '%v': %w", source, err)
		}
	}

	for key := range mutable {
		if _, err := resolveValue(mutable, key, map[string]bool{}, true); err != nil {
			return nil, err
		}
	}

	return Configuration(mutable), nil

}

func init() {

	ioc.PutFactory(func(sources []ConfigSource) (Configuration, error) {
		return CreateConfiguration(sources)
	})

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
 * Value resolving
 */

var placeholderRegexp *regexp.Regexp = regexp.MustCompile("\\${\\s*([^{}]+)\\s*}")

func resolveValue(env map[string]string, key string, traveled map[string]bool, record bool) (string, *CyclicLoopError) {

	value, exist := env[key]
	if !exist {
		return "", nil
	}

	if t, p := traveled[key]; p && t {
		return "", newCyclicLoopError(key, value)
	}
	traveled[key] = true

	match := placeholderRegexp.FindStringSubmatch(value)
	for match != nil {

		subvalue, err := resolveValue(env, match[1], traveled, record)
		if err != nil {
			return "", err.push(key, value)
		}

		value = strings.Replace(value, match[0], subvalue, -1)

		match = placeholderRegexp.FindStringSubmatch(value)

	}

	traveled[key] = false

	if record {
		env[key] = value
	}

	return value, nil

}

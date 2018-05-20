package config

import (
	"bytes"
	"regexp"
	"sort"
	"strings"
)

/*
 * ConfigSource
 */

type ConfigSource interface {
	GetPriority() int
	LoadEnv() map[string]string
}

/*
 * Configuration
 */

type Configuration struct {
	env map[string]string
}

func (self *Configuration) GetEnv() map[string]string {
	return self.env
}

/*
 * Sort sources by priority
 */

type byPriority []ConfigSource

func (self byPriority) Len() int {
	return len(self)
}

func (self byPriority) Swap(i, j int) {
	self[i], self[j] = self[j], self[i]
}

func (self byPriority) Less(i, j int) bool {
	return self[i].GetPriority() < self[j].GetPriority()
}

/*
 * Factory
 */

func CreateConfiguration(sources []ConfigSource) *Configuration {

	sort.Sort(byPriority(sources))

	env := make(map[string]string)
	for _, source := range sources {
		for key, value := range source.LoadEnv() {
			env[key] = value
		}
	}

	resolved := make(map[string]bool)
	pathMap := make(map[string]bool)

	for key := range env {
		if _, err := resolveValue(env, resolved, key, pathMap); err != nil {
			panic(err)
		}
	}

	return &Configuration{env}

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

	var buffer bytes.Buffer

	buffer.WriteString("Cyclic loop detected: ")

	buffer.WriteString(self.loop[0].key)
	buffer.WriteString(":'")
	buffer.WriteString(self.loop[0].value)
	for i := 1; i < len(self.loop); i++ {
		buffer.WriteString("' -> ")
		buffer.WriteString(self.loop[i].key)
		buffer.WriteString(":'")
		buffer.WriteString(self.loop[i].value)
	}
	buffer.WriteString("'")

	return buffer.String()

}

/*
 * Value resolving
 */

var placeholderRegexp *regexp.Regexp = regexp.MustCompile("\\${\\s*([^{}]+)\\s*}")

func resolveValue(env map[string]string, resolved map[string]bool, key string, pathMap map[string]bool) (string, *CyclicLoopError) {

	value := env[key]
	if _, res := resolved[key]; res {
		return value, nil
	}

	match := placeholderRegexp.FindStringSubmatch(value)
	for match != nil {

		subkey := match[1]
		subvalue := ""
		var err *CyclicLoopError

		if inpath, pres := pathMap[subkey]; pres && inpath {
			return "", newCyclicLoopError(key, value)
		}

		if _, exist := env[subkey]; exist {

			pathMap[subkey] = true
			subvalue, err = resolveValue(env, resolved, subkey, pathMap)
			pathMap[subkey] = false

			if err != nil {
				return "", err.push(key, value)
			}

		}

		value = strings.Replace(value, match[0], subvalue, -1)

		match = placeholderRegexp.FindStringSubmatch(value)

	}

	env[key] = value
	resolved[key] = true

	return value, nil

}

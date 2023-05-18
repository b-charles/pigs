package config

import (
	"sync/atomic"

	"github.com/b-charles/pigs/ioc"
	"github.com/b-charles/pigs/json"
)

/*
 * Tests
 */

var testPriority = &atomic.Int64{}

type TestSourceEntry struct {
	values   map[string]string
	priority int
}

func TestMap(values map[string]string) {
	ioc.TestPut(&TestSourceEntry{
		values:   values,
		priority: int(testPriority.Add(1)),
	}, func(ConfigSource) {})
}

func Test(key, value string) {
	TestMap(map[string]string{key: value})
}

func (self *TestSourceEntry) GetPriority() int {
	return self.priority
}

func (self *TestSourceEntry) LoadEnv(config MutableConfig) error {

	for key, value := range self.values {
		config.Set(key, value)
	}

	return nil

}

func (self *TestSourceEntry) Json() json.JsonNode {
	return json.NewJsonObjectStrings(self.values)
}

func (self *TestSourceEntry) String() string {
	return self.Json().String()
}

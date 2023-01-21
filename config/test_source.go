package config

import (
	"fmt"
	"sync/atomic"

	"github.com/b-charles/pigs/ioc"
)

/*
 * Tests
 */

var testPriority = &atomic.Int64{}

type TestSourceEntry struct {
	values   map[string]string
	priority int
}

func Test(key, value string) {
	ioc.TestPut(&TestSourceEntry{
		values:   map[string]string{key: value},
		priority: int(testPriority.Add(1)),
	}, func(ConfigSource) {})
}

func TestMap(values map[string]string) {
	ioc.TestPut(&TestSourceEntry{
		values:   values,
		priority: int(testPriority.Add(1)),
	}, func(ConfigSource) {})
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

func (self *TestSourceEntry) String() string {
	return fmt.Sprintf("Test entries: %s", stringify(self.values))
}

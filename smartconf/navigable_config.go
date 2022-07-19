package smartconf

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/b-charles/pigs/config"
	"github.com/b-charles/pigs/ioc"
)

type NavConfig struct {
	root     *NavConfig
	path     string
	value    string
	children map[string]*NavConfig
}

func (self *NavConfig) Root() *NavConfig {
	if self.root == nil {
		return self
	} else {
		return self.root
	}
}

func (self *NavConfig) Path() string {
	return self.path
}

func (self *NavConfig) Value() string {
	return self.value
}

func keyLess(a, b string) bool {

	if a == b {
		return false
	}

	int_a, a_err := strconv.Atoi(a)
	int_b, b_err := strconv.Atoi(b)

	if (a_err == nil) && (b_err == nil) {
		return int_a <= int_b
	} else if a_err == nil {
		return true
	} else if b_err == nil {
		return false
	} else {
		return a <= b
	}

}

func (self *NavConfig) Keys() []string {
	keys := make([]string, 0, len(self.children))
	for k := range self.children {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		return keyLess(keys[i], keys[j])
	})
	return keys
}

func path(root, key string) string {
	if root == "" {
		return key
	} else {
		return fmt.Sprintf("%s.%s", root, key)
	}
}

func (self *NavConfig) Child(key string) *NavConfig {
	if child, pres := self.children[key]; pres {
		return child
	}
	child := &NavConfig{self.root, path(self.path, key), "", map[string]*NavConfig{}}
	self.children[key] = child
	return child
}

func (self *NavConfig) Get(key string) *NavConfig {

	conf := self
	for i, k := range strings.Split(key, ".") {

		if i == 0 && k == "" {
			conf = conf.root
		} else {
			conf = conf.Child(k)
		}

	}

	return conf

}

func NewNavMap(config config.Configuration) *NavConfig {

	root := &NavConfig{nil, "", "", map[string]*NavConfig{}}

	for _, key := range config.Keys() {
		keys := strings.Split(key, ".")
		insert(keys, config.Get(key), root)
	}

	return root

}

func insert(keys []string, value string, navKey *NavConfig) {

	if len(keys) == 0 {
		navKey.value = value
	} else {
		insert(keys[1:], value, navKey.Child(keys[0]))
	}

}

func init() {

	ioc.PutFactory(func(config config.Configuration) (*NavConfig, error) {
		return NewNavMap(config), nil
	})

}

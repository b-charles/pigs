package smartconfig

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/b-charles/pigs/config"
	"github.com/b-charles/pigs/ioc"
)

type NavConfig interface {
	Root() NavConfig
	Parent() NavConfig
	Path() string
	Value() string
	Keys() []string
	Child(string) NavConfig
	Get(string) NavConfig
}

type NavConfigImpl struct {
	root     *NavConfigImpl
	parent   *NavConfigImpl
	path     string
	value    string
	children map[string]*NavConfigImpl
}

func (self *NavConfigImpl) Root() NavConfig {
	return self.root
}

func (self *NavConfigImpl) Parent() NavConfig {
	return self.parent
}

func (self *NavConfigImpl) Path() string {
	return self.path
}

func (self *NavConfigImpl) Value() string {
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

func (self *NavConfigImpl) Keys() []string {
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

func (self *NavConfigImpl) child(key string) *NavConfigImpl {
	if child, pres := self.children[key]; pres {
		return child
	}
	child := &NavConfigImpl{
		root:     self.root,
		parent:   self,
		path:     path(self.path, key),
		value:    "",
		children: map[string]*NavConfigImpl{},
	}
	self.children[key] = child
	return child
}

func (self *NavConfigImpl) Child(key string) NavConfig {
	return self.child(key)
}

func (self *NavConfigImpl) Get(key string) NavConfig {

	conf := self
	for i, k := range strings.Split(key, ".") {
		if i == 0 && k == "" {
			conf = conf.root
		} else {
			conf = conf.child(k)
		}
	}

	return conf

}

func NewNavMap(config config.Configuration) (*NavConfigImpl, error) {

	root := &NavConfigImpl{nil, nil, "", "", map[string]*NavConfigImpl{}}
	root.root = root

	for _, key := range config.Keys() {
		keys := strings.Split(key, ".")
		insert(keys, config.Get(key), root)
	}

	return root, nil

}

func insert(keys []string, value string, navKey *NavConfigImpl) {

	if len(keys) == 0 {
		navKey.value = value
	} else {
		insert(keys[1:], value, navKey.child(keys[0]))
	}

}

func init() {
	ioc.PutFactory(NewNavMap, func(NavConfig) {})
}

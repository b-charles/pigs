package smartconfig

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/b-charles/pigs/config"
	"github.com/b-charles/pigs/ioc"
	"github.com/b-charles/pigs/json"
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

type navConfigImpl struct {
	root     *navConfigImpl
	parent   *navConfigImpl
	path     string
	value    string
	children map[string]*navConfigImpl
}

func (self *navConfigImpl) Root() NavConfig {
	return self.root
}

func (self *navConfigImpl) Parent() NavConfig {
	return self.parent
}

func (self *navConfigImpl) Path() string {
	return self.path
}

func (self *navConfigImpl) Value() string {
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

func (self *navConfigImpl) Keys() []string {
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

func (self *navConfigImpl) child(key string) *navConfigImpl {
	if child, pres := self.children[key]; pres {
		return child
	}
	child := &navConfigImpl{
		root:     self.root,
		parent:   self,
		path:     path(self.path, key),
		value:    "",
		children: map[string]*navConfigImpl{},
	}
	self.children[key] = child
	return child
}

func (self *navConfigImpl) Child(key string) NavConfig {
	return self.child(key)
}

func (self *navConfigImpl) Get(key string) NavConfig {

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

func (self *navConfigImpl) Json() json.JsonNode {

	if self == nil {

		return json.JSON_NULL

	} else if len(self.children) > 0 {

		keys := make([]string, 0, len(self.children))
		for k := range self.children {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		b := json.NewJsonBuilder()

		if self.value != "" {
			b.SetString(".", self.value)
		}
		for _, k := range keys {
			b.Set(json.EscapePath(k), self.children[k].Json())
		}

		return b.Build()

	} else if self.value != "" {

		return json.JsonString(self.value)

	} else {

		return json.JSON_NULL

	}

}

func (self *navConfigImpl) String() string {
	return self.Json().String()
}

func NewNavMap(config config.Configuration) (NavConfig, error) {

	root := &navConfigImpl{nil, nil, "", "", map[string]*navConfigImpl{}}
	root.root = root

	for _, key := range config.Keys() {
		keys := strings.Split(key, ".")
		insert(keys, config.Get(key), root)
	}

	return root, nil

}

func insert(keys []string, value string, navKey *navConfigImpl) {

	if len(keys) == 0 {
		navKey.value = value
	} else {
		insert(keys[1:], value, navKey.child(keys[0]))
	}

}

func init() {
	ioc.PutNamedFactory("Navigable configuration",
		NewNavMap,
		func(NavConfig) {})
}

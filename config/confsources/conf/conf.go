package conf

import (
	"fmt"
	"sync"

	. "github.com/l3eegbee/pigs/config/confsources"
	"github.com/l3eegbee/pigs/ioc"
)

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
	ioc.Put(DefaultConfigSourceInstance(), "DefaultConfigSource", "DefaultConfigSources")
}

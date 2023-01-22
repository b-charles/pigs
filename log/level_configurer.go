package log

import (
	"fmt"
	"strings"

	"github.com/b-charles/pigs/ioc"
	"github.com/b-charles/pigs/smartconfig"
)

type LevelConfigurer interface {
	GetLevel(string) Level
}

type levelConfigurerImpl struct {
	config smartconfig.NavConfig
}

var ROOT_CONFIG = "log.level"

func (self *levelConfigurerImpl) GetLevel(name string) Level {

	if strings.HasPrefix(name, ".") {
		panic(fmt.Errorf("Invalid logger name '%s': a name can not starts with '.'.", name))
	}

	rootConfig := self.config.Get(ROOT_CONFIG)

	config := rootConfig.Get(strings.ToLower(name))
	for config.Value() == "" && config != rootConfig {
		config = config.Parent()
	}

	if lvl, err := ParseLevel(config.Value()); err != nil {
		panic(fmt.Errorf("Can not parse log level defined by '%s': %w", config.Path(), err))
	} else {
		return lvl
	}

}

func (self *levelConfigurerImpl) String() string {
	return "LevelConfigurer"
}

func init() {

	ioc.DefaultPutFactory(func(config smartconfig.NavConfig) (*levelConfigurerImpl, error) {
		return &levelConfigurerImpl{config}, nil
	}, func(LevelConfigurer) {})

}

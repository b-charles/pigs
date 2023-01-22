package log

import (
	"github.com/b-charles/pigs/ioc"
	"github.com/b-charles/pigs/json"
)

type LoggerFactory interface {
	NewLogger(name string) Logger
}

type loggerFactoryImpl struct {
	jsons           json.Jsons
	levelConfigurer LevelConfigurer
	contextualizers []Contextualizer
	appenders       []Appender
}

func (self *loggerFactoryImpl) NewLogger(name string) Logger {
	return &loggerImpl{
		jsons:           self.jsons,
		name:            name,
		level:           self.levelConfigurer.GetLevel(name),
		contextualizers: self.contextualizers,
		appenders:       self.appenders,
	}
}

func init() {

	ioc.DefaultPutFactory(func(
		jsons json.Jsons,
		levelConfigurer LevelConfigurer,
		contextualizers []Contextualizer,
		appenders []Appender) (*loggerFactoryImpl, error) {
		return &loggerFactoryImpl{jsons, levelConfigurer, contextualizers, appenders}, nil
	}, func(LoggerFactory) {})

}
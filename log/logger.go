package log

import (
	"fmt"

	"github.com/b-charles/pigs/config"
	"github.com/b-charles/pigs/ioc"
	"github.com/b-charles/pigs/json"
)

type Logger interface {
	GetName() string

	GetLevel() Level
	IsLevelEnabled(level Level) bool

	Builder(level Level) LogBuilder
	Log(level Level, path string, value any)

	IsTraceEnabled() bool
	Trace() LogBuilder
	TraceLog(path string, value any)

	IsDebugEnabled() bool
	Debug() LogBuilder
	DebugLog(path string, value any)

	IsInfoEnabled() bool
	Info() LogBuilder
	InfoLog(path string, value any)

	IsWarnEnabled() bool
	Warn() LogBuilder
	WarnLog(path string, value any)

	IsErrorEnabled() bool
	Error() LogBuilder
	ErrorLog(path string, value any)

	IsFatalEnabled() bool
	Fatal() LogBuilder
	FatalLog(path string, value any)

	AddContextualizer(Contextualizer) Logger
	AddContext(string, any) Logger
}

// Impl

type loggerImpl struct {
	jsons           json.Jsons
	name            string
	level           Level
	contextualizers []Contextualizer
	appenders       []Appender
}

func (self *loggerImpl) GetName() string {
	return self.name
}

func (self *loggerImpl) GetLevel() Level {
	return self.level
}

func (self *loggerImpl) IsLevelEnabled(level Level) bool {
	return self.level <= level
}

func (self *loggerImpl) Builder(level Level) LogBuilder {

	if self.IsLevelEnabled(level) {

		builder := &logBuilderImpl{
			jsons:     self.jsons,
			builder:   json.NewJsonBuilder(),
			appenders: self.appenders,
		}

		for _, contextualizer := range self.contextualizers {
			contextualizer.AddContext(self, level, builder)
		}

		return builder

	} else {
		return logBuilderNullInst
	}

}

func (self *loggerImpl) Log(level Level, path string, value any) {
	self.Builder(level).Set(path, value).Log()
}

func (self *loggerImpl) IsTraceEnabled() bool {
	return self.IsLevelEnabled(Trace)
}

func (self *loggerImpl) Trace() LogBuilder {
	return self.Builder(Trace)
}

func (self *loggerImpl) TraceLog(path string, value any) {
	self.Log(Trace, path, value)
}

func (self *loggerImpl) IsDebugEnabled() bool {
	return self.IsLevelEnabled(Debug)
}

func (self *loggerImpl) Debug() LogBuilder {
	return self.Builder(Debug)
}

func (self *loggerImpl) DebugLog(path string, value any) {
	self.Log(Debug, path, value)
}

func (self *loggerImpl) IsInfoEnabled() bool {
	return self.IsLevelEnabled(Info)
}

func (self *loggerImpl) Info() LogBuilder {
	return self.Builder(Info)
}

func (self *loggerImpl) InfoLog(path string, value any) {
	self.Log(Info, path, value)
}

func (self *loggerImpl) IsWarnEnabled() bool {
	return self.IsLevelEnabled(Warn)
}

func (self *loggerImpl) Warn() LogBuilder {
	return self.Builder(Warn)
}

func (self *loggerImpl) WarnLog(path string, value any) {
	self.Log(Warn, path, value)
}

func (self *loggerImpl) IsErrorEnabled() bool {
	return self.IsLevelEnabled(Error)
}

func (self *loggerImpl) Error() LogBuilder {
	return self.Builder(Error)
}

func (self *loggerImpl) ErrorLog(path string, value any) {
	self.Log(Error, path, value)
}

func (self *loggerImpl) IsFatalEnabled() bool {
	return self.IsLevelEnabled(Fatal)
}

func (self *loggerImpl) Fatal() LogBuilder {
	return self.Builder(Fatal)
}

func (self *loggerImpl) FatalLog(path string, value any) {
	self.Log(Fatal, path, value)
}

func (self *loggerImpl) AddContextualizer(contextualizer Contextualizer) Logger {
	return &loggerImpl{
		jsons:           self.jsons,
		name:            self.name,
		level:           self.level,
		contextualizers: append(self.contextualizers, contextualizer),
		appenders:       self.appenders,
	}
}

func (self *loggerImpl) AddContext(key string, value any) Logger {
	return self.AddContextualizer(NewStaticContextualizer(key, value))
}

var DEFAULT_LOGGER_NAME = "root"

func init() {

	config.Set(ROOT_CONFIG, Info.String())
	config.Set(fmt.Sprintf("%s.%s", ROOT_CONFIG, DEFAULT_LOGGER_NAME), fmt.Sprintf("${%s}", ROOT_CONFIG))

	ioc.DefaultPutNamedFactory(fmt.Sprintf("Logger '%s'", DEFAULT_LOGGER_NAME),
		func(loggerFactory LoggerFactory) (Logger, error) {
			return loggerFactory.NewLogger(DEFAULT_LOGGER_NAME), nil
		}, func(Logger) {})

}

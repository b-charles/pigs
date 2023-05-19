package log

import (
	"time"

	"github.com/b-charles/pigs/ioc"
	"github.com/benbjohnson/clock"
)

// Interface

type Contextualizer interface {
	GetPriority() int
	AddContext(Logger, Level, LogBuilder)
}

// Default

type DateLevelContextualizer Contextualizer

var DATE_LEVEL_CONTEXTUALIZER_PRIORITY = 0

type DateLevelContextualizerImpl struct {
	clock clock.Clock
}

func (self *DateLevelContextualizerImpl) GetPriority() int {
	return DATE_LEVEL_CONTEXTUALIZER_PRIORITY
}

func (self *DateLevelContextualizerImpl) AddContext(logger Logger, level Level, builder LogBuilder) {
	builder.Set("time", self.clock.Now().Format(time.RFC3339Nano))
	builder.Set("level", level)
}

func init() {

	ioc.DefaultPutNamedFactory("Time and level contextualizer (default)",
		func(clock clock.Clock) *DateLevelContextualizerImpl {
			return &DateLevelContextualizerImpl{clock}
		}, func(DateLevelContextualizer) {})

	ioc.PutNamedFactory("Time and level contextualizer (promoter)",
		func(c DateLevelContextualizer) (Contextualizer, error) { return c, nil })

}

// Static

type StaticContextualizer struct {
	priority int
	context  map[string]any
}

func (self *StaticContextualizer) GetPriority() int {
	return self.priority
}

func (self *StaticContextualizer) AddContext(logger Logger, level Level, builder LogBuilder) {
	for k, v := range self.context {
		builder.Set(k, v)
	}
}

func NewStaticContextualizer(priority int, key string, value any) *StaticContextualizer {
	return &StaticContextualizer{priority, map[string]any{key: value}}
}

func NewStaticContextualizerMap(priority int, m map[string]any) *StaticContextualizer {
	r := make(map[string]any, len(m))
	for k, v := range m {
		r[k] = v
	}
	return &StaticContextualizer{priority, r}
}

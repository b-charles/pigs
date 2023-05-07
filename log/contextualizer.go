package log

import (
	"time"

	"github.com/b-charles/pigs/ioc"
	"github.com/benbjohnson/clock"
)

// Interface

type Contextualizer interface {
	AddContext(Logger, Level, LogBuilder)
}

// Default

type DefaultContextualizer struct {
	clock clock.Clock
}

func (self *DefaultContextualizer) AddContext(logger Logger, level Level, builder LogBuilder) {
	builder.Set("time", self.clock.Now().Format(time.RFC3339Nano))
	builder.Set("level", level)
}

func init() {

	ioc.DefaultPutNamedFactory("Time and level contextualizer",
		func(clock clock.Clock) *DefaultContextualizer {
			return &DefaultContextualizer{clock}
		}, func(Contextualizer) {})

}

// Static

type StaticContextualizer map[string]any

func (self StaticContextualizer) AddContext(logger Logger, level Level, builder LogBuilder) {
	for k, v := range self {
		builder.Set(k, v)
	}
}

func newStaticContextualizer(key string, value any) StaticContextualizer {
	return map[string]any{key: value}
}

func newStaticContextualizerMap(m map[string]any) StaticContextualizer {
	r := make(map[string]any, len(m))
	for k, v := range m {
		r[k] = v
	}
	return r
}

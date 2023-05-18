package log

import (
	"fmt"
	"strings"

	"github.com/b-charles/pigs/ioc"
	"github.com/b-charles/pigs/json"
	"github.com/b-charles/pigs/smartconfig"
)

type Level uint

const (
	Trace Level = iota
	Debug
	Info
	Warn
	Error
	Fatal
)

func (self Level) String() string {

	switch self {
	case Trace:
		return "TRACE"
	case Debug:
		return "DEBUG"
	case Info:
		return "INFO"
	case Warn:
		return "WARN"
	case Error:
		return "ERROR"
	case Fatal:
		return "FATAL"
	default:
		panic(fmt.Errorf("Unexpected level %#v", self))
	}

}

func (self Level) Json() json.JsonNode {
	return json.JsonString(self.String())
}

func ParseLevel(value string) (Level, error) {
	trim := strings.ToUpper(strings.TrimSpace(value))
	switch {
	case trim == "":
		return Info, nil
	case trim == "TRACE":
		return Trace, nil
	case trim == "DEBUG":
		return Debug, nil
	case trim == "INFO":
		return Info, nil
	case trim == "WARN":
		return Warn, nil
	case trim == "ERROR":
		return Error, nil
	case trim == "FATAL":
		return Fatal, nil
	default:
		return Info, fmt.Errorf("Value '%s' is not a valid log level.", value)
	}
}

type LevelMarshaller func(Level) (json.JsonNode, error)
type LevelUnmarshaller func(json.JsonNode) (Level, error)
type LevelParser func(string) (Level, error)

func init() {

	ioc.DefaultPutNamed("Log level Json marshaller (default)",
		func(level Level) (json.JsonNode, error) {
			return level.Json(), nil
		}, func(LevelMarshaller) {})

	ioc.PutNamedFactory("Log level Json marshaller (promoter)",
		func(m LevelMarshaller) (json.JsonMarshaller, error) { return m, nil })

	ioc.DefaultPutNamed("Log level Json unmarshaller (default)",
		func(node json.JsonNode) (Level, error) {
			if node.IsString() {
				return ParseLevel(node.AsString())
			} else {
				return Info, fmt.Errorf("Can not parse json %v as a log level.", node)
			}
		}, func(LevelUnmarshaller) {})

	ioc.PutNamedFactory("Log level Json unmarshaller (promoter)",
		func(u LevelUnmarshaller) (json.JsonUnmarshaller, error) { return u, nil })

	ioc.DefaultPutNamed("Log level parser (default)",
		ParseLevel, func(LevelParser) {})

	ioc.PutNamedFactory("Log level parser (promoter)",
		func(p LevelParser) (smartconfig.Parser, error) { return p, nil })

}

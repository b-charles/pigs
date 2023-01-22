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

func init() {

	ioc.Put(func(level Level) (json.JsonNode, error) {
		return json.JsonString(level.String()), nil
	}, func(json.JsonMarshaller) {})
	ioc.Put(func(node json.JsonNode) (Level, error) {
		return ParseLevel(node.AsString())
	}, func(json.JsonUnmarshaller) {})

	ioc.Put(ParseLevel, func(smartconfig.Parser) {})

}
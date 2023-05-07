package smartconfig

import (
	"strconv"

	"github.com/b-charles/pigs/ioc"
)

// Default parsers

type StringParser func(string) (string, error)
type Float64Parser func(string) (float64, error)
type IntParser func(string) (int, error)
type BoolParser func(string) (bool, error)

func init() {

	ioc.DefaultPutNamed("String parser (default)",
		func(value string) (string, error) {
			return value, nil
		}, func(StringParser) {})

	ioc.PutNamedFactory("String parser (promoter)",
		func(p StringParser) (Parser, error) { return p, nil })

	ioc.DefaultPutNamed("Float64 parser (default)",
		func(value string) (float64, error) {
			return strconv.ParseFloat(value, 64)
		}, func(Float64Parser) {})

	ioc.PutNamedFactory("Float64 parser (promoter)",
		func(p Float64Parser) (Parser, error) { return p, nil })

	ioc.DefaultPutNamed("Int parser (default)",
		strconv.Atoi, func(IntParser) {})

	ioc.PutNamedFactory("Int parser (promoter)",
		func(p IntParser) (Parser, error) { return p, nil })

	ioc.DefaultPutNamed("Bool parser (default)",
		strconv.ParseBool, func(BoolParser) {})

	ioc.PutNamedFactory("Bool parser (promoter)",
		func(p BoolParser) (Parser, error) { return p, nil })

}

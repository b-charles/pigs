package config

import (
	"fmt"
	"os"
	"regexp"

	"github.com/b-charles/pigs/ioc"
	"github.com/b-charles/pigs/json"
)

var (
	valueRegexp       *regexp.Regexp = regexp.MustCompile("^--?([^=]+)=(.*)$")
	valueRegexpSimple *regexp.Regexp = regexp.MustCompile("^--?([^=]+)='(.*)'$")
	valueRegexpDouble *regexp.Regexp = regexp.MustCompile("^--?([^=]+)=\"(.*)\"$")
)

func keyvalueArg(arg string) (string, string, bool) {

	match := valueRegexpSimple.FindStringSubmatch(arg)
	if match == nil {
		match = valueRegexpDouble.FindStringSubmatch(arg)
	}
	if match == nil {
		match = valueRegexp.FindStringSubmatch(arg)
	}

	if match == nil {
		return "", "", false
	} else {
		return match[1], match[2], true
	}

}

var (
	boolRegexp   *regexp.Regexp = regexp.MustCompile("^--?([^=]+)$")
	noboolRegexp *regexp.Regexp = regexp.MustCompile("^--?no-([^=]+)$")
)

func keyboolArg(arg string) (string, string, bool) {

	match := noboolRegexp.FindStringSubmatch(arg)
	if match != nil {
		return match[1], "false", true
	}

	match = boolRegexp.FindStringSubmatch(arg)
	if match != nil {
		return match[1], "true", true
	}

	return "", "", false

}

func ParseArgs(args []string) (map[string]string, error) {

	env := make(map[string]string)

	for _, arg := range args {

		if key, value, ok := keyvalueArg(arg); ok {
			env[key] = value
			continue
		}

		if key, value, ok := keyboolArg(arg); ok {
			env[key] = value
			continue
		}

		return env, fmt.Errorf("Can't parse argument '%s': unknown pattern.", arg)

	}

	return env, nil

}

type ArgsConfigSource ConfigSource

var CONFIG_SOURCE_PRIORITY_ARGS = 100

type ArgsConfigSourceImpl struct {
	source map[string]string
}

func (self *ArgsConfigSourceImpl) GetPriority() int {
	return CONFIG_SOURCE_PRIORITY_ARGS
}

func (self *ArgsConfigSourceImpl) LoadEnv(config MutableConfig) error {
	for k, v := range self.source {
		config.Set(k, v)
	}
	return nil
}

func (self *ArgsConfigSourceImpl) Json() json.JsonNode {
	return json.NewJsonObjectStrings(self.source)
}

func (self *ArgsConfigSourceImpl) String() string {
	return self.Json().String()
}

func init() {

	ioc.DefaultPutNamedFactory("Args config source (default)",
		func() (*ArgsConfigSourceImpl, error) {
			m, err := ParseArgs(os.Args[1:])
			return &ArgsConfigSourceImpl{m}, err
		}, func(ArgsConfigSource) {})

	ioc.PutNamedFactory("Args config source (promoter)",
		func(v ArgsConfigSource) (ConfigSource, error) { return v, nil })

}

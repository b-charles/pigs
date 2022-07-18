package config

import (
	"fmt"
	"os"
	"regexp"

	"github.com/b-charles/pigs/ioc"
)

var valueRegexp *regexp.Regexp = regexp.MustCompile("^--([^=]+)=(.*)$")
var valueRegexpSimple *regexp.Regexp = regexp.MustCompile("^--([^=]+)='(.*)'$")
var valueRegexpDouble *regexp.Regexp = regexp.MustCompile("^--([^=]+)=\"(.*)\"$")

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

var boolRegexp *regexp.Regexp = regexp.MustCompile("^--([^=]+)$")
var noboolRegexp *regexp.Regexp = regexp.MustCompile("^--no-([^=]+)$")

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

		return nil, fmt.Errorf("Can't parse argument '%s': unknown pattern.", arg)

	}

	return env, nil

}

type ArgsConfigSource map[string]string

func (self ArgsConfigSource) GetPriority() int {
	return CONFIG_SOURCE_PRIORITY_ARGS
}

func (self ArgsConfigSource) LoadEnv(config MutableConfig) error {
	for k, v := range self {
		config.Set(k, v)
	}
	return nil
}

func (self ArgsConfigSource) String() string {
	return stringify(self)
}

func NewArgsConfigSource() (ArgsConfigSource, error) {
	return ParseArgs(os.Args[1:])
}

func init() {
	ioc.PutFactory(NewArgsConfigSource, func(ConfigSource) {})
}

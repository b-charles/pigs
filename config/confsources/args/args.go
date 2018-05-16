package args

import (
	"os"
	"regexp"

	. "github.com/l3eegbee/pigs/config/confsources"
)

var valueRegexp *regexp.Regexp = regexp.MustCompile("^--([^=]+)=(.*)$")
var valueRegexpSimple *regexp.Regexp = regexp.MustCompile("^--([^=]+)='(.*)'$")
var valueRegexpDouble *regexp.Regexp = regexp.MustCompile("^--([^=]+)=\"(.*)\"$")

var boolRegexp *regexp.Regexp = regexp.MustCompile("^--([^=]+)$")
var noboolRegexp *regexp.Regexp = regexp.MustCompile("^--no-([^=]+)$")

func ParseArgs(args []string) map[string]string {

	env := make(map[string]string)

	for _, arg := range args {

		var match []string

		match = valueRegexpSimple.FindStringSubmatch(arg)
		if match == nil {
			match = valueRegexpDouble.FindStringSubmatch(arg)
		}
		if match == nil {
			match = valueRegexp.FindStringSubmatch(arg)
		}
		if match != nil {
			env[match[1]] = match[2]
			continue
		}

		match = noboolRegexp.FindStringSubmatch(arg)
		if match != nil {
			env[match[1]] = "false"
		}

		match = boolRegexp.FindStringSubmatch(arg)
		if match != nil {
			env[match[1]] = "true"
		}

	}

	return env

}

func NewArgsConfigSource() *SimpleConfigSource {

	env := ParseArgs(os.Args[1:])

	return &SimpleConfigSource{
		Priority: CONFIG_SOURCE_PRIORITY_ARGS,
		Env:      env,
	}

}

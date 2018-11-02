package init

import (
	"fmt"
	"regexp"

	"github.com/juju/loggo"
	"github.com/l3eegbee/pigs/ioc"
)

var (
	loglevelroot  string         = "log.level"
	loglevelreg   *regexp.Regexp = regexp.MustCompile(fmt.Sprintf("^%s\\.(.*)$", regexp.QuoteMeta(loglevelroot)))
	writernamereg *regexp.Regexp = regexp.MustCompile("^(.*)LogWriter$")
)

func init() {

	ioc.PutFactory(func(injected struct {
		Configuration map[string]string
		LogWriters    map[string]loggo.Writer
	}) *loggo.Context {

		context := loggo.NewContext(loggo.WARNING)

		// add writers
		for name, writer := range injected.LogWriters {
			if match := writernamereg.FindStringSubmatch(name); match != nil {
				name = match[1]
			}
			context.AddWriter(name, writer)
		}

		// set up log levels
		for key, value := range injected.Configuration {
			if match := loglevelreg.FindStringSubmatch(key); match != nil {
				if level, ok := loggo.ParseLevel(value); ok {

					name := match[1]
					if name == "root" {
						name = ""
					}

					context.GetLogger(name).SetLogLevel(level)

				}
			}
		}

		return context

	}, "LogContext")

}

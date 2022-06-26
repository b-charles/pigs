package init

import (
	"os"
	"strings"

	"github.com/juju/loggo"
	"github.com/juju/loggo/loggocolor"
	"github.com/b-charles/pigs/ioc"
)

var (
	configOut   string = "log.writer.default.out"
	configColor string = "log.writer.default.color"
	configLevel string = "log.writer.default.level"
)

func init() {

	ioc.PutFactory(func(injected struct {
		Configuration map[string]string
	}) loggo.Writer {

		out := os.Stdout
		if value, ok := injected.Configuration[configOut]; ok && strings.EqualFold(value, "stderr") {
			out = os.Stderr
		}

		color := true
		if value, ok := injected.Configuration[configColor]; ok && strings.EqualFold(value, "false") {
			color = false
		}

		level := loggo.UNSPECIFIED
		if value, ok := injected.Configuration[configLevel]; ok {
			level, _ = loggo.ParseLevel(value)
		}

		var writer loggo.Writer
		if color {
			writer = loggocolor.NewWriter(out)
		} else {
			writer = loggo.NewSimpleWriter(out, loggo.DefaultFormatter)
		}

		return loggo.NewMinimumLevelWriter(writer, level)

	}, "DefaultLogWriter", "LogWriters")

}

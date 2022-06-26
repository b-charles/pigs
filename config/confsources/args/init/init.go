package init

import (
	. "github.com/b-charles/pigs/config/confsources/args"
	"github.com/b-charles/pigs/ioc"
)

func init() {
	ioc.Put(NewArgsConfigSource(), "ArgsConfigSource", "ConfigSources")
}

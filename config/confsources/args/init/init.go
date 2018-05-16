package init

import (
	. "github.com/l3eegbee/pigs/config/confsources/args"
	"github.com/l3eegbee/pigs/ioc"
)

func init() {
	ioc.Put(NewArgsConfigSource(), "ArgsConfigSource", "ConfigSources")
}

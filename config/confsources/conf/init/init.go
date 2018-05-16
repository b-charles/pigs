package init

import (
	. "github.com/l3eegbee/pigs/config/confsources/conf"
	"github.com/l3eegbee/pigs/ioc"
)

func init() {
	ioc.Put(DefaultConfigSourceInstance(), "DefaultConfigSource", "DefaultConfigSources")
}

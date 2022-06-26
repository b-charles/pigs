package init

import (
	. "github.com/b-charles/pigs/config/confsources/conf"
	"github.com/b-charles/pigs/ioc"
)

func init() {
	ioc.Put(DefaultConfigSourceInstance(), "DefaultConfigSource", "DefaultConfigSources")
}

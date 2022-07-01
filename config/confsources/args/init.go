package args

import (
	"github.com/b-charles/pigs/ioc"
)

func init() {
	ioc.Put(NewArgsConfigSource(), "ArgsConfigSource", "ConfigSource")
}

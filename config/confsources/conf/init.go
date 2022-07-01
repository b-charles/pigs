package conf

import (
	"github.com/b-charles/pigs/ioc"
)

func init() {
	ioc.Put(DefaultConfigSourceInstance(), "DefaultConfigSource", "ConfigSource")
}

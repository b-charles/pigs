package init

import (
	"github.com/l3eegbee/pigs/ioc"
	. "github.com/l3eegbee/pigs/smartconf"
)

type EnvLoader interface {
	GetEnv() map[string]string
}

func init() {

	ioc.PutFactory(
		func(envLoader EnvLoader, parsers []interface{}) *SmartConf {
			return NewSmartConf(envLoader.GetEnv(), parsers)
		}, []string{"Configuration", "StringParsers"}, "SmartConfiguration")

}

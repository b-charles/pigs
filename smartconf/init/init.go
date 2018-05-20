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
		func(injected struct {
			Configuration map[string]string
			StringParsers []interface{}
		}) *SmartConf {
			return NewSmartConf(injected.Configuration, injected.StringParsers)
		}, "SmartConfiguration")

}

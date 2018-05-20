package smart

import (
	"github.com/l3eegbee/pigs/ioc"
	. "github.com/l3eegbee/pigs/smartconf"
	"github.com/pkg/errors"
)

func Inject(conf interface{}, root string, name string, aliases ...string) {

	ioc.PutFactory(
		func(injected struct {
			SmartConfiguration *SmartConf
		}) interface{} {

			if conf, err := injected.SmartConfiguration.Configure(root, conf); err == nil {
				return conf
			} else {
				panic(errors.Wrapf(err, "Error during smart configuration of %v", conf))
			}

		}, name, aliases...)

}

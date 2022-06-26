package smart

import (
	"fmt"

	"github.com/b-charles/pigs/ioc"
	. "github.com/b-charles/pigs/smartconf"
)

func Inject(conf interface{}, root string, name string, aliases ...string) {

	ioc.PutFactory(
		func(injected struct {
			SmartConfiguration *SmartConf
		}) interface{} {

			if conf, err := injected.SmartConfiguration.Configure(root, conf); err == nil {
				return conf
			} else {
				panic(fmt.Errorf("Error during smart configuration of %v: %w", conf, err))
			}

		}, name, aliases...)

}

package smartconfig

import (
	"strconv"

	"github.com/b-charles/pigs/ioc"
)

func init() {

	ioc.Put(func(value string) (string, error) {
		return value, nil
	}, func(Parser) {})
	ioc.Put(strconv.ParseBool, func(Parser) {})
	ioc.Put(strconv.Atoi, func(Parser) {})

}

package properties

import (
	enc "github.com/magiconair/properties"
)

func ParsePropertiesToEnv(content string) map[string]string {

	p := enc.MustLoadString(content)
	p.Prefix = ""
	p.Postfix = ""

	return p.Map()

}

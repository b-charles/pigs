package ioc

import (
	"github.com/benbjohnson/clock"
	"github.com/spf13/afero"
)

func init() {

	// File system
	DefaultPutFactory(afero.NewOsFs, func(afero.Fs) {})

	// Clock
	DefaultPutFactory(clock.New, func(clock.Clock) {})

}

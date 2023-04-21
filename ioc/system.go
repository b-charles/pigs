package ioc

import (
	"github.com/benbjohnson/clock"
	"github.com/spf13/afero"
)

func init() {

	// File system
	DefaultPutNamedFactory("Native file system", afero.NewOsFs, func(afero.Fs) {})

	// Clock
	DefaultPutNamedFactory("Native clock", clock.New, func(clock.Clock) {})

}

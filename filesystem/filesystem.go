package filesystem

import (
	"flag"

	"github.com/l3eegbee/pigs/ioc"
	"github.com/spf13/afero"
)

func init() {

	ioc.PutFactory(func() *afero.Afero {

		var fs afero.Fs

		if flag.Lookup("test.v") == nil {
			fs = afero.NewOsFs()
		} else {
			fs = afero.NewMemMapFs()
		}

		return &afero.Afero{Fs: fs}

	}, []string{}, "Filesystem")

}

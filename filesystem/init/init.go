package init

import (
	"github.com/l3eegbee/pigs/ioc"
	"github.com/spf13/afero"
)

func init() {

	ioc.PutFactory(func() *afero.Afero {

		return &afero.Afero{Fs: afero.NewOsFs()}

	}, []string{}, "Filesystem")

}

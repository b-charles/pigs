package confsources

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"

	"github.com/l3eegbee/pigs/config"
	"github.com/l3eegbee/pigs/ioc"

	"github.com/spf13/afero"

	_ "github.com/l3eegbee/pigs/filesystem"
)

func convertObjectInEnv(env map[string]string, root string, object interface{}) {

	if object == nil {
		return
	}

	switch value := object.(type) {

	case bool:
		env[root] = strconv.FormatBool(value)

	case string:
		env[root] = value

	case int:
		env[root] = strconv.FormatInt(int64(value), 10)

	case int64:
		env[root] = strconv.FormatInt(value, 10)

	case float64:
		env[root] = strconv.FormatFloat(value, 'G', -1, 64)

	case []interface{}:

		for idx, elt := range value {
			convertObjectInEnv(env, fmt.Sprintf("%s[%d]", root, idx), elt)
		}

	case map[string]interface{}:

		rootWithPoint := root
		if len(root) > 0 {
			rootWithPoint = rootWithPoint + "."
		}

		for k, v := range value {
			convertObjectInEnv(env, rootWithPoint+k, v)
		}

	case map[interface{}]interface{}:

		rootWithPoint := root
		if len(root) > 0 {
			rootWithPoint = rootWithPoint + "."
		}

		for k, v := range value {
			convertObjectInEnv(env, fmt.Sprintf("%s%v", rootWithPoint, k), v)
		}

	default:
		panic(fmt.Errorf("Unexpected type: %v (%v)", value, reflect.TypeOf(value)))

	}

}

func RegisterFileConfig(
	priority int,
	ext string,
	parser func(string) map[string]string,
	name string,
	aliases ...string) {

	ioc.PutFactory(func(filesystem *afero.Afero) *config.SimpleConfigSource {

		app, err := filepath.Abs(os.Args[0])
		app = strings.TrimSuffix(app, filepath.Ext(app))
		file := app + ext

		content, err := filesystem.ReadFile(file)
		if err != nil {
			return nil
		}

		return &config.SimpleConfigSource{
			Priority: priority,
			Env:      parser(string(content)),
		}

	}, []string{"Filesystem"}, name, append(aliases, "ConfigSources")...)

}

func RegisterFormatedEnvVarConfig(
	priority int,
	prefix string,
	parser func(string) map[string]string,
	name string,
	aliases ...string) {

	ioc.PutFactory(func(envvar config.ConfigSource) *config.SimpleConfigSource {

		app := filepath.Base(os.Args[0])
		app = strings.TrimSuffix(app, filepath.Ext(app))
		app = ConvertEnvVarKey(app + prefix)

		env := envvar.LoadEnv()
		if content, ok := env[app]; ok {
			return &config.SimpleConfigSource{
				Priority: priority,
				Env:      parser(content),
			}
		}

		return nil

	}, []string{"EnvVarConfigSource"}, name, append(aliases, "ConfigSources")...)

}

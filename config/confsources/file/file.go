package file

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"

	. "github.com/l3eegbee/pigs/config/confsources"
	"github.com/l3eegbee/pigs/ioc"
)

func ConvertObjectInEnv(env map[string]string, root string, object interface{}) {

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
			ConvertObjectInEnv(env, fmt.Sprintf("%s.%d", root, idx), elt)
		}

	case map[string]interface{}:

		rootWithPoint := root
		if len(root) > 0 {
			rootWithPoint = rootWithPoint + "."
		}

		for k, v := range value {
			ConvertObjectInEnv(env, rootWithPoint+k, v)
		}

	case map[interface{}]interface{}:

		rootWithPoint := root
		if len(root) > 0 {
			rootWithPoint = rootWithPoint + "."
		}

		for k, v := range value {
			ConvertObjectInEnv(env, fmt.Sprintf("%s%v", rootWithPoint, k), v)
		}

	default:
		panic(fmt.Errorf("Unexpected type: %v (%v)", value, reflect.TypeOf(value)))

	}

}

type Filesystem interface {
	ReadFile(file string) (string, error)
}

func RegisterFileConfig(
	priority int,
	ext string,
	parser func(string) map[string]string,
	name string,
	aliases ...string) {

	ioc.PutFactory(func(injected struct {
		Filesystem Filesystem
	}) *SimpleConfigSource {

		app, err := filepath.Abs(os.Args[0])
		app = strings.TrimSuffix(app, filepath.Ext(app))
		file := app + ext

		content, err := injected.Filesystem.ReadFile(file)
		if err != nil {
			return nil
		}

		return &SimpleConfigSource{
			Priority: priority,
			Env:      parser(string(content)),
		}

	}, name, append(aliases, "ConfigSources")...)

}

type EnvVar interface {
	LoadEnv() map[string]string
}

func RegisterFormatedEnvVarConfig(
	priority int,
	prefix string,
	parser func(string) map[string]string,
	name string,
	aliases ...string) {

	ioc.PutFactory(func(injected struct {
		EnvVar EnvVar
	}) *SimpleConfigSource {

		app := filepath.Base(os.Args[0])
		app = strings.TrimSuffix(app, filepath.Ext(app))
		app = ConvertEnvVarKey(app + prefix)

		env := injected.EnvVar.LoadEnv()
		if content, ok := env[app]; ok {
			return &SimpleConfigSource{
				Priority: priority,
				Env:      parser(content),
			}
		}

		return nil

	}, name, append(aliases, "ConfigSources")...)

}

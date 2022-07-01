package conf

import (
	. "github.com/b-charles/pigs/config/confsources"
	"github.com/b-charles/pigs/ioc"
)

func envForTests(defaultConfigSource *DefaultConfigSource, priority int, env map[string]string) *SimpleConfigSource {

	res := make(map[string]string)
	for k, v := range defaultConfigSource.Env {
		res[k] = v
	}
	for k, v := range env {
		res[k] = v
	}

	return &SimpleConfigSource{
		Priority: priority,
		Env:      res,
	}

}

func SetEnvForTestsWithPriority(priority int, env map[string]string) {

	ioc.TestPutFactory(func(injected struct {
		DefaultConfigSource *DefaultConfigSource
	}) *SimpleConfigSource {
		return envForTests(injected.DefaultConfigSource, priority, env)
	}, "DefaultTestConfigSource", "ConfigSource")

}

func SetEnvForTests(env map[string]string) {
	SetEnvForTestsWithPriority(CONFIG_SOURCE_PRIORITY_TESTS, env)
}

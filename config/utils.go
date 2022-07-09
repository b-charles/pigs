package config

// Priorities

var CONFIG_SOURCE_PRIORITY_DEFAULT = -999
var CONFIG_SOURCE_PRIORITY_ENV_VAR = 0
var CONFIG_SOURCE_PRIORITY_ARGS = 100
var CONFIG_SOURCE_PRIORITY_TESTS = 999

// Simple impl

type SimpleConfigSource struct {
	Priority int
	Env      map[string]string
}

func (self *SimpleConfigSource) GetPriority() int {
	return self.Priority
}

func (self *SimpleConfigSource) LoadEnv() map[string]string {
	return self.Env
}

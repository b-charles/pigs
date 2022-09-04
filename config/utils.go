package config

import (
	"sort"
	"strings"
)

// Priorities

var (
	CONFIG_SOURCE_PRIORITY_DEFAULT    = -999
	CONFIG_SOURCE_PRIORITY_ENV_VAR    = 0
	CONFIG_SOURCE_PRIORITY_ARGS       = 100
	CONFIG_SOURCE_PRIORITY_JSON_FILES = 200
	CONFIG_SOURCE_PRIORITY_TESTS      = 999
)

// stringify

func stringify(m map[string]string) string {

	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})

	var builder strings.Builder

	for i, k := range keys {
		if i > 0 {
			builder.WriteString(" ")
		}
		builder.WriteString(k)
		builder.WriteString(":\"")
		builder.WriteString(m[k])
		builder.WriteString("\"")
	}

	return builder.String()

}

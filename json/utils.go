package json

import (
	"fmt"
	"reflect"
	"sort"
	"strconv"
)

// types

var errorType = reflect.TypeOf(func(error) {}).In(0)
var stringType = reflect.TypeOf(func(string) {}).In(0)
var jsonType = reflect.TypeOf(func(JsonNode) {}).In(0)
var jsonerType = reflect.TypeOf(func(Jsoner) {}).In(0)

// [T] to Json mappers

func jsonToJson(v JsonNode) JsonNode { return v }
func stringToJson(v string) JsonNode { return JsonString(v) }
func floatToJson(v float64) JsonNode { return JsonFloat(v) }
func intToJson(v int) JsonNode       { return JsonInt(v) }
func boolToJson(v bool) JsonNode {
	if v {
		return JSON_TRUE
	} else {
		return JSON_FALSE
	}
}

// Reflect type slice to Json mapper

func reflectTypeInfos(t reflect.Type) (sortable string, packagePath string, displayName string) {

	kind := t.Kind()

	if kind == reflect.Array {

		s, p, n := reflectTypeInfos(t.Elem())
		return fmt.Sprintf("%s[.]", s), p, fmt.Sprintf("[%s]", n)

	} else if kind == reflect.Chan {

		s, p, n := reflectTypeInfos(t.Elem())
		return fmt.Sprintf("%s$", s), p, fmt.Sprintf("$%s", n)

	} else if kind == reflect.Map {

		ks, kp, kn := reflectTypeInfos(t.Key())
		vs, _, _ := reflectTypeInfos(t.Elem())
		return fmt.Sprintf("%s{%s}", ks, vs), kp, fmt.Sprintf("{%s: %s}", kn, t.Elem().Name())

	} else if kind == reflect.Pointer {

		s, p, n := reflectTypeInfos(t.Elem())
		return fmt.Sprintf("%s*", s), p, fmt.Sprintf("*%s", n)

	} else if kind == reflect.Slice {

		s, p, n := reflectTypeInfos(t.Elem())
		return fmt.Sprintf("%s[]", s), p, fmt.Sprintf("[]%s", n)

	} else {

		packagePath := t.PkgPath()
		displayName := t.Name()
		return fmt.Sprintf("%s.%s", packagePath, displayName), packagePath, displayName

	}

}

func ReflectTypeSliceToJson(types []reflect.Type) JsonNode {

	// packagePath -> sortable -> displayName
	allInfos := make(map[string]map[string]string)

	for _, k := range types {
		s, p, n := reflectTypeInfos(k)
		if m, ok := allInfos[p]; ok {
			m[s] = n
		} else {
			nm := make(map[string]string)
			nm[s] = n
			allInfos[p] = nm
		}
	}

	allPackageNames := make([]string, 0, len(allInfos))
	for k := range allInfos {
		allPackageNames = append(allPackageNames, k)
	}
	sort.Strings(allPackageNames)

	jsonInfosByPackage := make(map[string]JsonNode)
	for _, p := range allPackageNames {

		infos := allInfos[p]

		sortables := make([]string, 0, len(infos))
		for k := range infos {
			sortables = append(sortables, k)
		}
		sort.Strings(sortables)

		nodes := make([]JsonNode, 0, len(sortables))
		for _, s := range sortables {
			nodes = append(nodes, JsonString(infos[s]))
		}

		jsonInfosByPackage[p] = NewJsonArray(nodes)

	}

	return newJsonObject(jsonInfosByPackage, allPackageNames)

}

// string comparaison

func stringsLess(a, b string) bool {

	if a == b {
		return false
	}

	int_a, a_err := strconv.Atoi(a)
	int_b, b_err := strconv.Atoi(b)

	if (a_err == nil) && (b_err == nil) {
		return int_a < int_b
	} else if a_err == nil {
		return true
	} else if b_err == nil {
		return false
	} else {
		return a < b
	}

}

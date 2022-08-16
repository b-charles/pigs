package core

import (
	"regexp"
	"strconv"
)

type jsonBuilderNode struct {
	value    JsonNode
	children *sortedMap[*jsonBuilderNode]
}

func newJsonBuilderNode() *jsonBuilderNode {
	return &jsonBuilderNode{nil, newEmptySortedMap[*jsonBuilderNode]()}
}

func (self *jsonBuilderNode) put(path []string, value JsonNode) {

	if len(path) == 0 {
		self.value = value
		self.children = newEmptySortedMap[*jsonBuilderNode]()
	} else {

		self.value = nil

		child, ok := self.children.get(path[0])
		if !ok {
			child = newJsonBuilderNode()
			self.children.put(path[0], child)
		}

		child.put(path[1:], value)

	}

}

func (self *jsonBuilderNode) json() JsonNode {

	if self.value != nil {

		return self.value

	} else if self.children.len() == 0 {

		return JSON_NULL

	} else if array, parsed, max := self.arrayCheck(); array {

		elements := make([]JsonNode, max+1)
		for k, e := range self.children.m {
			elements[parsed[k]] = e.json()
		}

		return &JsonArray{elements}

	} else {

		n := self.children.len()
		members := &sortedMap[JsonNode]{make(map[string]JsonNode, n), make([]string, 0, n)}
		for _, k := range self.children.s {
			members.put(k, self.children.m[k].json())
		}

		return &JsonObject{members}

	}

}

func (self *jsonBuilderNode) arrayCheck() (bool, map[string]int, int) {

	parsed := make(map[string]int)
	max := 0

	for _, k := range self.children.s {
		if k[0] != '[' {
			return false, nil, 0
		} else if v, err := strconv.Atoi(k[1:]); err != nil {
			return false, nil, 0
		} else {
			parsed[k] = v
			if max < v {
				max = v
			}
		}
	}

	return true, parsed, max

}

type JsonBuilder struct {
	root *jsonBuilderNode
}

func NewJsonBuilder() *JsonBuilder {
	return &JsonBuilder{newJsonBuilderNode()}
}

func (self *JsonBuilder) Build() JsonNode {
	return self.root.json()
}

var pathregexp = regexp.MustCompile(`\.?([^.\[]+)|(\[\d+)\]`)

func parsePath(path string) []string {

	parsed := pathregexp.FindAllStringSubmatch(path, -1)

	paths := make([]string, 0, len(parsed))
	for _, p := range parsed {
		if p[2] != "" {
			paths = append(paths, p[2])
		} else {
			paths = append(paths, p[1])
		}
	}

	return paths

}

func (self *JsonBuilder) Set(path string, value JsonNode) {
	self.root.put(parsePath(path), value)
}

func (self *JsonBuilder) SetString(path string, value string) *JsonBuilder {
	self.Set(path, JsonString(value))
	return self
}

func (self *JsonBuilder) SetFloat(path string, value float64) *JsonBuilder {
	self.Set(path, JsonFloat(value))
	return self
}

func (self *JsonBuilder) SetInt(path string, value int) *JsonBuilder {
	self.Set(path, JsonInt(value))
	return self
}

func (self *JsonBuilder) SetBool(path string, value bool) *JsonBuilder {
	if value {
		self.Set(path, JSON_TRUE)
	} else {
		self.Set(path, JSON_FALSE)
	}
	return self
}

func (self *JsonBuilder) SetEmptyObject(path string) *JsonBuilder {
	self.Set(path, JSON_EMPTY_OBJECT)
	return self
}

func (self *JsonBuilder) SetEmptyArray(path string) *JsonBuilder {
	self.Set(path, JSON_EMPTY_ARRAY)
	return self
}

func (self *JsonBuilder) SetNull(path string) *JsonBuilder {
	self.Set(path, JSON_NULL)
	return self
}

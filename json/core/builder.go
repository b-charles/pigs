package core

import (
	"fmt"
	"regexp"
	"strconv"
)

type JsonBuilder struct {
	root JsonNode
}

func NewJsonBuilder() *JsonBuilder {
	return &JsonBuilder{JSON_NULL}
}

func (self *JsonBuilder) Build() JsonNode {
	return self.root
}

var pathregexp = regexp.MustCompile(`\.?([^.\[]+)|\[(\d+)\]`)

func getKey(parsedPathPart []string) (bool, string, int) {
	if parsedPathPart[2] != "" {
		if index, err := strconv.Atoi(parsedPathPart[2]); err != nil {
			panic(err)
		} else {
			return false, "", index
		}
	} else {
		return true, parsedPathPart[1], 0
	}
}

func (self *JsonBuilder) set(path string, value JsonNode) {
	parsedPath := pathregexp.FindAllStringSubmatch(path, -1)
	if newRoot, err := recursiveSet(self.root, parsedPath, value); err != nil {
		panic(err)
	} else {
		self.root = newRoot
	}
}

func recursiveSet(node JsonNode, parsedPath [][]string, value JsonNode) (JsonNode, error) {

	if len(parsedPath) == 0 {

		return value, nil

	} else {

		if isObj, key, index := getKey(parsedPath[0]); isObj {

			if node.IsNull() {

				if newValue, err := recursiveSet(JSON_NULL, parsedPath[1:], value); err != nil {
					return node, fmt.Errorf("Can not add member '%v': %w", key, err)
				} else {
					newNode := newJsonObject()
					newNode.set(key, newValue)
					return newNode, nil
				}

			} else if !node.IsObject() {

				return node, fmt.Errorf("Can not add member '%v' in %v.", key, node)

			} else {

				casted := node.(*JsonObject)
				if sub := casted.GetMember(key); sub.IsNull() {

					if newMember, err := recursiveSet(JSON_NULL, parsedPath[1:], value); err != nil {
						return node, fmt.Errorf("Can not add member '%v': %w", key, err)
					} else {
						casted.set(key, newMember)
					}

				} else {

					if _, err := recursiveSet(sub, parsedPath[1:], value); err != nil {
						return node, fmt.Errorf("Can not add member '%v': %w", key, err)
					}

				}

				return node, nil

			}

		} else {

			if node.IsNull() {

				if newValue, err := recursiveSet(JSON_NULL, parsedPath[1:], value); err != nil {
					return node, fmt.Errorf("Can not add element at %v: %w", index, err)
				} else {
					newNode := newJsonArray()
					newNode.set(index, newValue)
					return newNode, nil
				}

			} else if !node.IsArray() {

				return node, fmt.Errorf("Can not add element at %v in %v.", index, node)

			} else {

				casted := node.(*JsonArray)
				if sub := casted.GetElement(index); sub.IsNull() {

					if newElement, err := recursiveSet(JSON_NULL, parsedPath[1:], value); err != nil {
						return node, fmt.Errorf("Can not add element at %v: %w", key, err)
					} else {
						casted.set(index, newElement)
					}

				} else {

					if _, err := recursiveSet(sub, parsedPath[1:], value); err != nil {
						return node, fmt.Errorf("Can not add member '%v': %w", key, err)
					}

				}

				return node, nil

			}

		}

	}

}

func (self *JsonBuilder) SetString(path string, value string) *JsonBuilder {
	self.set(path, JsonString(value))
	return self
}

func (self *JsonBuilder) SetFloat(path string, value float64) *JsonBuilder {
	self.set(path, JsonFloat(value))
	return self
}

func (self *JsonBuilder) SetInt(path string, value int) *JsonBuilder {
	self.set(path, JsonInt(value))
	return self
}

func (self *JsonBuilder) SetBool(path string, value bool) *JsonBuilder {
	if value {
		self.set(path, JSON_TRUE)
	} else {
		self.set(path, JSON_FALSE)
	}
	return self
}

func (self *JsonBuilder) SetNull(path string) *JsonBuilder {
	self.set(path, JSON_NULL)
	return self
}

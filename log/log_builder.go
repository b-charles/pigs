package log

import "github.com/b-charles/pigs/json"

type LogBuilder interface {
	Set(path string, value any) LogBuilder
	SetString(path string, value string) LogBuilder
	SetFloat(path string, value float64) LogBuilder
	SetInt(path string, value int) LogBuilder
	SetBool(path string, value bool) LogBuilder
	SetEmptyObject(path string) LogBuilder
	SetEmptyArray(path string) LogBuilder
	SetNull(path string) LogBuilder
	Log()
}

// Default impl

type logBuilderImpl struct {
	jsons     json.Jsons
	builder   *json.JsonBuilder
	appenders []Appender
}

func (self *logBuilderImpl) Set(path string, value any) LogBuilder {
	if node, err := self.jsons.Marshal(value); err != nil {
		panic(err)
	} else {
		self.builder.Set(path, node)
	}
	return self
}

func (self *logBuilderImpl) SetString(path string, value string) LogBuilder {
	self.builder.SetString(path, value)
	return self
}

func (self *logBuilderImpl) SetFloat(path string, value float64) LogBuilder {
	self.builder.SetFloat(path, value)
	return self
}

func (self *logBuilderImpl) SetInt(path string, value int) LogBuilder {
	self.builder.SetInt(path, value)
	return self
}

func (self *logBuilderImpl) SetBool(path string, value bool) LogBuilder {
	self.builder.SetBool(path, value)
	return self
}

func (self *logBuilderImpl) SetEmptyObject(path string) LogBuilder {
	self.builder.SetEmptyObject(path)
	return self
}

func (self *logBuilderImpl) SetEmptyArray(path string) LogBuilder {
	self.builder.SetEmptyArray(path)
	return self
}

func (self *logBuilderImpl) SetNull(path string) LogBuilder {
	self.builder.SetNull(path)
	return self
}

func (self *logBuilderImpl) Log() {
	node := self.builder.Build()
	for _, appender := range self.appenders {
		appender.Append(node)
	}
}

// Null impl

type logBuilderNull struct{}

func (self *logBuilderNull) Set(path string, value any) LogBuilder {
	return self
}

func (self *logBuilderNull) SetString(path string, value string) LogBuilder {
	return self
}

func (self *logBuilderNull) SetFloat(path string, value float64) LogBuilder {
	return self
}

func (self *logBuilderNull) SetInt(path string, value int) LogBuilder {
	return self
}

func (self *logBuilderNull) SetBool(path string, value bool) LogBuilder {
	return self
}

func (self *logBuilderNull) SetEmptyObject(path string) LogBuilder {
	return self
}

func (self *logBuilderNull) SetEmptyArray(path string) LogBuilder {
	return self
}

func (self *logBuilderNull) SetNull(path string) LogBuilder {
	return self
}

func (self *logBuilderNull) Log() {}

var logBuilderNullInst *logBuilderNull = &logBuilderNull{}

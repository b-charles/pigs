package ioc

import "time"

type ContainerInfo interface {
	TestMode() bool
	CreationTime() time.Time
	StartingTime() time.Time
	ClosingTime() time.Time
}

type containerInfoImpl struct {
	testMode     bool
	creationTime time.Time
	startingTime time.Time
	closingTime  time.Time
}

func (self *containerInfoImpl) start(container *Container, testMode bool) {
	self.testMode = testMode
	self.startingTime = time.Now()
}

func (self *containerInfoImpl) close(container *Container) {
	self.closingTime = time.Now()
}

func (self *containerInfoImpl) TestMode() bool {
	return self.testMode
}

func (self *containerInfoImpl) CreationTime() time.Time {
	return self.creationTime
}

func (self *containerInfoImpl) StartingTime() time.Time {
	return self.startingTime
}

func (self *containerInfoImpl) ClosingTime() time.Time {
	return self.closingTime
}

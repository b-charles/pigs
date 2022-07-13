package ioc

import "reflect"

// preCall

type PreCallAwared interface {
	Precall(method reflect.Value) error
}

type noopPreCallAwaredHandler struct{}

func (self noopPreCallAwaredHandler) Precall(reflect.Value) error {
	return nil
}

// postInst

type PostInstAwared interface {
	Postinst(method reflect.Value, args []reflect.Value) error
}

type noopPostInstAwaredHandler struct{}

func (self noopPostInstAwaredHandler) Postinst(reflect.Value, []reflect.Value) error {
	return nil
}

// preClose

type PreCloseAwared interface {
	Preclose()
}

type noopPreCloseAwaredHandler struct{}

func (self noopPreCloseAwaredHandler) Preclose() {}

// postClose

type PostCloseAwared interface {
	Postclose()
}

type noopPostCloseAwaredHandler struct{}

func (self noopPostCloseAwaredHandler) Postclose() {}

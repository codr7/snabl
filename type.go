package snabl

import (
	"io"
)

type Type interface {
	Init(name string)
	Name() string
	Dump(val V, out io.Writer) error
}

type BasicType struct {
	name string
}

func (self *BasicType) Init(name string) {
	self.name = name
}

func (self *BasicType) Name() string {
	return self.name
}

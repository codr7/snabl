package snabl

import (
	"io"
	"strings"
)

type V struct {
	t VT
	d any
}

func (self *V) Init(t VT, d any) {
	self.t = t
	self.d = d
}

func (self *V) Type() VT {
	return self.t
}

func (self *V) Data() any {
	return self.d
}

func (self *V) Dump(out io.Writer) error {
	return self.t.Dump(*self, out)
}

func (self *V) String() string {
	var out strings.Builder
	self.Dump(&out)
	return out.String()
}

type VT interface {
	Name() string
	Dump(val V, out io.Writer) error
}

type BasicVT struct {
	name string
}

func (self *BasicVT) Init(name string) {
	self.name = name
}

func (self *BasicVT) Name() string {
	return self.name
}

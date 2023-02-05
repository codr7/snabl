package snabl

import (
	"io"
	"strings"
)

type V struct {
	t Type
	d any
}

func (self *V) Init(t Type, d any) {
	self.t = t
	self.d = d
}

func (self *V) Type() Type {
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

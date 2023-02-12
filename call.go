package snabl

import (
	"fmt"
	"strings"
)

type Call struct {
	parent *Call
	pos *Pos
	fun *Fun
	args []V
	retPc Pc
}

func (self *Call) Init(parent *Call, pos *Pos, fun *Fun, args []V, retPc Pc) *Call {
	self.parent = parent
	self.pos = pos
	self.fun = fun
	self.args = args
	self.retPc = retPc
	return self
}

func (self Call) String() string {
	var out strings.Builder

	if self.pos != nil {
		fmt.Fprintf(&out, "%v ", *self.pos)
	}

	out.WriteString(self.fun.String())
	return out.String()
}

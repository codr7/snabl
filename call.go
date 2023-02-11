package snabl

import (
	"fmt"
	"strings"
)

type Call struct {
	pos *Pos
	fun *Fun
	args []V
	retPc Pc
}

func (self Call) String() string {
	var out strings.Builder

	if self.pos != nil {
		fmt.Fprintf(&out, "%v ", *self.pos)
	}

	out.WriteString(self.fun.String())
	return out.String()
}

package snabl

import (
	"fmt"
	"strings"
)

type E struct {
	pos *Pos
	msg string
}

func NewE(pos *Pos, spec string, args...interface{}) *E {
	return new(E).Init(pos, spec, args...)
}

func (self *E) Init(pos *Pos, spec string, args...interface{}) *E {
	self.pos = pos
	self.msg = fmt.Sprintf(spec, args...)
	return self
}

func (self *E) Pos() *Pos {
	return self.pos
}

func (self *E) Msg() string {
	return self.msg
}

func (self *E) Error() string {
	var msg strings.Builder

	if self.pos != nil {
		fmt.Fprintf(&msg, "%v ", self.pos)
	}
	
	fmt.Fprintf(&msg, "Error: %v", self.msg)
	return msg.String()
}

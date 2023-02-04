package snabl

import (
	"fmt"
)

type Pos struct {
	source string
	line, column uint
}

func NewPos(source string) *Pos {
	return new(Pos).Init(source)
}

func (self *Pos) Init(source string) *Pos {
	self.source = source
	self.line = 1
	self.column = 1
	return self
}

func (self *Pos) String() string {
	return fmt.Sprintf("%v@%v:%v", self.source, self.line, self.column)
}

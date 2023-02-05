package snabl

import (
	"fmt"
)

type Pos struct {
	source string
	line, column int
}

func NewPos(source string, line, column int) *Pos {
	return new(Pos).Init(source, line, column)
}

func (self *Pos) Init(source string, line, column int) *Pos {
	self.source = source
	self.line = line
	self.column = column
	return self
}

func (self *Pos) Source() string {
	return self.source
}

func (self *Pos) Line() int {
	return self.line
}

func (self *Pos) Column() int {
	return self.column
}

func (self Pos) String() string {
	return fmt.Sprintf("%v@%v:%v", self.source, self.line, self.column)
}

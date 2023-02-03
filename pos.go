package snabl

var NilPos = NewPos("n/a")

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

package snabl

import (
)

type PrimBody = func(self *Prim, vm *Vm, pos *Pos) error
type PrimArity = uint

type Prim struct {
	name string
	arity uint
	body PrimBody
}

func NewPrim(name string, arity uint, body PrimBody) *Prim {
	return new(Prim).Init(name, arity, body)
}

func (self *Prim) Init(name string, arity uint, body PrimBody) *Prim {
	self.name = name
	self.arity = arity
	self.body = body
	return self
}

func (self *Prim) Call(vm *Vm, pos *Pos) error {
	return self.body(self, vm, pos)
}

func (self *Prim) String() string {
	return self.name
}

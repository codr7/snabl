package snabl

import (
)

type PrimBody = func(self *Prim, vm *Vm, pos *Pos) error

type Prim struct {
	tag Tag
	name string
	arity int
	body PrimBody
}

func NewPrim(vm *Vm, name string, arity int, body PrimBody) *Prim {
	return new(Prim).Init(vm, name, arity, body)
}

func (self *Prim) Init(vm *Vm, name string, arity int, body PrimBody) *Prim {
	self.tag = vm.Tag(&vm.AbcLib.PrimType, self)
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

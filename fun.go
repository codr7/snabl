package snabl

import (
)

type Fun struct {
	tag Tag
	name string
	pc Pc
	args []string
}

func NewFun(vm *Vm, name string, pc Pc, args...string) *Fun {
	return new(Fun).Init(vm, name, pc, args...)
}

func (self *Fun) Init(vm *Vm, name string, pc Pc, args...string) *Fun {
	self.tag = vm.Tag(V{t: &vm.AbcLib.FunType, d: self})
	self.name = name
	self.pc = pc
	self.args = args
	return self
}

func (self *Fun) Arity() int {
	return len(self.args)
}

func (self *Fun) String() string {
	return self.name
}

func (self *Fun) ArgIndex(arg string) int {
	for i, a := range self.args {
		if a == arg {
			return i
		}
	}

	return -1
}

package snabl

import (
)

type Fun struct {
	name string
	pc Pc
	args []string
	argOffsTag Tag
}

func NewFun(vm *Vm, name string, pc Pc, args...string) *Fun {
	return new(Fun).Init(vm, name, pc, args...)
}

func (self *Fun) Init(vm *Vm, name string, pc Pc, args...string) *Fun {
	self.name = name
	self.pc = pc
	self.args = args

	if len(args) > 0 {
		self.argOffsTag = vm.Tag(&vm.AbcLib.IntType, -1)
	}

	return self
}

func (self *Fun) Pc() Pc {
	return self.pc
}

func (self *Fun) ArgOffsTag() Tag {
	return self.argOffsTag
}

func (self *Fun) String() string {
	return self.name
}

func (self *Fun) ArgIndex(arg string) int {
	for i, a := range self.args {
		if a == arg {
			return len(self.args) - i - 1;
		}
	}

	return -1
}

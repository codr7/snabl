package snabl

import (
)

type MacroBody = func(self *Macro, args *Forms, vm *Vm, env Env, pos Pos) error

type Macro struct {
	name string
	arity int
	body MacroBody
}

func NewMacro(name string, arity int, body MacroBody) *Macro {
	return new(Macro).Init(name, arity, body)
}

func (self *Macro) Init(name string, arity int, body MacroBody) *Macro {
	self.name = name
	self.arity = arity
	self.body = body
	return self
}

func (self *Macro) Emit(args *Forms, vm *Vm, env Env, pos Pos) error {
	if args.Len() < self.arity {
		return vm.E(&pos, "Not enough arguments")
	}
	
	if vm.Debug {
		vm.EmitPos(pos)
	}
	
	return self.body(self, args, vm, env, pos)
}

func (self *Macro) String() string {
	return self.name
}

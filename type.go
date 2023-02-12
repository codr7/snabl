package snabl

import (
	"io"
)

type Type interface {
	Init(name string)
	Name() string
	Bool(val V) bool
	Eq(left, right V) bool
	Dump(val V, out io.Writer) error
	Write(val V, out io.Writer) error
	Emit(val V, args *Forms, vm *Vm, env Env, pos Pos) error
}

type SeqType interface {
	Type
	Len(val V) int
}

type BasicType struct {
	name string
}

func (self *BasicType) Init(name string) {
	self.name = name
}

func (self *BasicType) Name() string {
	return self.name
}

func (self *BasicType) Bool(val V) bool {
	return true
}

func (Self *BasicType) Eq(left, right V) bool {
	return left.d == right.d
}

func (self *BasicType) Emit(val V, args *Forms, vm *Vm, env Env, pos Pos) error {	
	vm.EmitTag(val.t, val.d)
	return nil
}


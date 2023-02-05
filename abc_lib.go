package snabl

import (
	"fmt"
	"io"
)

var Abc AbcLib

func init () {
	Abc.Init()
}

type IntType struct {
	BasicType
}

func (self *IntType) Dump(val V, out io.Writer) error {
	_, err := fmt.Fprintf(out, "%v", val.d.(int))
	return err
}

type FunType struct {
	BasicType
}

func (self *FunType) Dump(val V, out io.Writer) error {
	_, err := io.WriteString(out, val.d.(*Fun).String())
	return err
}

type PosType struct {
	BasicType
}

func (self *PosType) Dump(val V, out io.Writer) error {
	_, err := fmt.Fprintf(out, "%v", val.d.(Pos))
	return err
}

type PrimType struct {
	BasicType
}

func (self *PrimType) Dump(val V, out io.Writer) error {
	_, err := io.WriteString(out, val.d.(*Prim).String())
	return err
}

type StringType struct {
	BasicType
}

func (self *StringType) Dump(val V, out io.Writer) error {
	_, err := io.WriteString(out, val.d.(string))
	return err
}

type AbcLib struct {
	BasicLib
	FunType FunType
	IntType IntType
	PosType PosType
	PrimType PrimType
	StringType StringType

	AddPrim, DumpPrim, FailPrim Prim
}

func (self *AbcLib) Init() {
	self.BasicLib.Init("abc")
	
	self.BindType(&self.FunType, "Fun")
	self.BindType(&self.IntType, "Int")
	self.BindType(&self.PosType, "Pos")
	self.BindType(&self.PrimType, "Prim")
	self.BindType(&self.StringType, "String")

	self.BindPrim(&self.AddPrim, "+", 2, func(self *Prim, vm *Vm, pos *Pos) error {
		b := vm.Stack.Pop().d.(int)
		a := vm.Stack.Top()
		a.Init(&Abc.IntType, a.d.(int) + b)
		return nil
	})
	
	 self.BindPrim(&self.DumpPrim, "dump", 1, func(self *Prim, vm *Vm, pos *Pos) error {
		vm.Stack.Pop().Dump(vm.Stdout)
		return nil
	})

	 self.BindPrim(&self.FailPrim, "fail", 1, func(self *Prim, vm *Vm, pos *Pos) error {
		return vm.E(pos, vm.Stack.Pop().String())
	})
}

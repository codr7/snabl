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
	BasicVT
}

func (self *IntType) Dump(val V, out io.Writer) error {
	_, err := fmt.Fprintf(out, "%v", val.d.(int))
	return err
}

type FunType struct {
	BasicVT
}

func (self *FunType) Dump(val V, out io.Writer) error {
	_, err := fmt.Fprintf(out, "%v", val.d.(*Fun))
	return err
}

type PosType struct {
	BasicVT
}

func (self *PosType) Dump(val V, out io.Writer) error {
	_, err := fmt.Fprintf(out, "%v", val.d.(Pos))
	return err
}

type PrimType struct {
	BasicVT
}

func (self *PrimType) Dump(val V, out io.Writer) error {
	_, err := fmt.Fprintf(out, "%v", val.d.(*Prim))
	return err
}

type AbcLib struct {
	BasicLib
	FunType FunType
	IntType IntType
	PosType PosType
	PrimType PrimType
}

func (self *AbcLib) Init() {
	self.BasicLib.Init("abc")
	self.FunType.Init("Fun")
	self.IntType.Init("Int")
	self.PosType.Init("Pos")
	self.PrimType.Init("Prim")

	self.BindPrim("dump", 1, func(vm *Vm, pos Pos) error {
		vm.Stack.Pop().Dump(vm.Stdout)
		return nil
	})	
}

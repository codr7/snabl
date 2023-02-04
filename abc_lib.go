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

type AbcLib struct {
	BasicLib
	FunType FunType
	IntType IntType
	PosType PosType
}

func (self *AbcLib) Init() {
	self.BasicLib.Init("abc")
	self.FunType.Init("Fun")
	self.IntType.Init("Int")
	self.PosType.Init("Pos")
}

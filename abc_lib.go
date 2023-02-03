package snabl

import (
	"fmt"
	"io"
)

type IntType struct {
	BasicVT
}

func (self *IntType) Dump(val V, out io.Writer) error {
	_, err := fmt.Fprintf(out, "%v", val.d.(int))
	return err
}


type AbcLib struct {
	BasicLib
	IntType IntType
}

func (self *AbcLib) Init() {
	self.BasicLib.Init("abc")
	self.IntType.Init("Int")
}

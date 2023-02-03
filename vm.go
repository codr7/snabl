package snabl

import (
	"fmt"
	"io"
	"os"
)

type Pc = uint64

type Vm struct {
	AbcLib AbcLib
	Stdin io.Reader
	Stdout io.Writer
	Code []Op
	Stack Stack
}

func (self *Vm) Init() {
	self.AbcLib.Init()
	self.Stdin = os.Stdin
	self.Stdout = os.Stdout
}

func (self *Vm) Emit(n Pc) Pc {
	pc := self.EmitPc()
	
	for i := Pc(0); i < n; i++ {
		self.Code = append(self.Code, 0)
	}

	return pc
}

func (self *Vm) EmitPc() Pc {
	return Pc(len(self.Code))
}

func (self *Vm) Eval(pc *Pc) error {
	for {
		op := self.Code[*pc]

		OpTrace(*pc, op, self.Stdout)
		
		switch id := OpId(op); id {
		case ADD_OP:
			b := self.Stack.Pop()
			a := self.Stack.Top()
			a.Init(&self.AbcLib.IntType, a.Data().(int) + b.Data().(int))
			*pc++;
		case PUSH_INT_OP:
			self.Stack.Push(&self.AbcLib.IntType, PushIntVal(op))
			*pc++
		case STOP_OP:
			*pc++
			return nil
		default:
			panic(fmt.Sprintf("Invalid op id: %v", id))
		}
	}
}

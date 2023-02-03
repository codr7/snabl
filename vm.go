package snabl

import (
	"fmt"
	"io"
	"os"
)

type Pc = uint

type Vm struct {
	Debug bool
	Trace bool
	
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

func (self *Vm) E(pos *Pos, spec string, args...interface{}) error {
	err := NewE(pos, spec, args...)
	
	if self.Debug {
		panic(err.Error())
	}

	return err
}

func (self *Vm) EmitNoTrace() Pc {
	pc := self.EmitPc()
	self.Code = append(self.Code, 0)
	return pc
}

func (self *Vm) Emit() Pc {
	if self.Trace {
		self.Code[self.EmitNoTrace()] = TraceOp()
	}
	
	return self.EmitNoTrace()
}

func (self *Vm) EmitPc() Pc {
	return Pc(len(self.Code))
}

func (self *Vm) Eval(pc *Pc) error {
	for {
		op := self.Code[*pc]
		
		switch id := op.Id(); id {
		case ADD_OP:
			b := self.Stack.Pop()
			a := self.Stack.Top()
			a.Init(&self.AbcLib.IntType, a.Data().(int) + b.Data().(int))
			*pc++;
		case PUSH_INT_OP:
			self.Stack.Push(&self.AbcLib.IntType, op.PushIntVal())
			*pc++
		case STOP_OP:
			*pc++
			return nil
		case TRACE_OP:
			*pc++
			self.Code[*pc].Trace(*pc, self.Stdout)
			*pc++
		default:
			panic(fmt.Sprintf("Invalid op id: %v", id))
		}
	}
}

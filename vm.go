package snabl

import (
	"fmt"
	"io"
	"os"
)

type Pc = int
type Tag = int

type Vm struct {
	Debug bool
	Trace bool
	
	Stdin io.Reader
	Stdout io.Writer
	
	Tags []V
	Env BasicEnv
	
	Code []Op
	Stack Stack
}

func (self *Vm) Init() {
	self.Stdin = os.Stdin
	self.Stdout = os.Stdout
	self.Env.Init()
}

func (self *Vm) Tag(t VT, d any) Tag {
	i := len(self.Tags)
	self.Tags = append(self.Tags, V{t: t, d: d})
	return i
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

func (self *Vm) EmitPos(pos Pos) {
	tag := self.Tag(&Abc.PosType, pos)
	self.Code[self.EmitNoTrace()] = PosOp(tag)
}

func (self *Vm) EmitPc() Pc {
	return Pc(len(self.Code))
}

func (self *Vm) Eval(pc *Pc) error {
	var pos *Pos
	
	for {
		op := self.Code[*pc]
		
		switch id := op.Id(); id {
		case ADD_OP:
			b := self.Stack.Pop()
			a := self.Stack.Top()
			a.Init(&Abc.IntType, a.Data().(int) + b.Data().(int))
			*pc++;
		case ARG_OP:
			v := self.Stack.items[self.Tags[op.ArgTag()].d.(ArgOffs) + op.ArgIndex()]
			self.Stack.Push(v.t, v.d);
			*pc++
		case ARG_OFFS_OP:
			self.Tags[op.ArgOffsTag()].d = self.Stack.Len()
			*pc++
		case POS_OP:
			p := self.Tags[op.PosTag()].Data().(Pos)
			pos = &p
			*pc++
		case PUSH_INT_OP:
			self.Stack.Push(&Abc.IntType, op.PushIntVal())
			*pc++
		case STOP_OP:
			*pc++
			return nil
		case TRACE_OP:
			*pc++
			self.Code[*pc].Trace(*pc, pos, self.Stdout)
		default:
			panic(fmt.Sprintf("Invalid op id: %v", id))
		}
	}
}

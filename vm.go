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

	AbcLib AbcLib
	
	Tags []V
	
	Code []Op
	Calls []Call
	Stack Stack
	
	env Env
	fun *Fun
}

func (self *Vm) Init() {
	self.Stdin = os.Stdin
	self.Stdout = os.Stdout
	self.AbcLib.Init(self)
	self.env = NewEnv()
}

func (self *Vm) Tag(t Type, d any) Tag {
	i := len(self.Tags)
	self.Tags = append(self.Tags, V{t: t, d: d})
	return i
}

func (self *Vm) E(pos *Pos, spec string, args...interface{}) *E {
	err := NewE(pos, spec, args...)
	
	if self.Debug {
		panic(err.Error())
	}

	return err
}

func (self *Vm) Env() Env {
	return self.env
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
	tag := self.Tag(&self.AbcLib.PosType, pos)
	self.Code[self.EmitNoTrace()] = PosOp(tag)
}

func (self *Vm) EmitString(str string) {
	tag := self.Tag(&self.AbcLib.StringType, str)
	self.Code[self.Emit()] = PushOp(tag) 
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
			a.Init(&self.AbcLib.IntType, a.Data().(int) + b.Data().(int))
			*pc++;
		case ARG_OP:
			v := self.Stack.items[self.Tags[op.ArgTag()].d.(int) - op.ArgIndex() - 1]
			self.Stack.Push(v.t, v.d);
			*pc++
		case ARG_OFFS_OP:
			self.Tags[op.ArgOffsTag()].d = self.Stack.Len()
			*pc++
		case CALL_FUN_OP:
			f := self.Tags[op.CallFunTag()].d.(*Fun)
			self.Calls = append(self.Calls, Call{pos: pos, fun: f, retPc: *pc+1})
			*pc = f.pc
		case CALL_PRIM_OP:
			if err := self.Tags[op.ArgOffsTag()].d.(*Prim).Call(self, pos); err != nil {
				return err
			}

			*pc++
		case GOTO_OP:
			*pc = op.GotoPc()
		case POS_OP:
			p := self.Tags[op.PosTag()].Data().(Pos)
			pos = &p
			*pc++
		case PUSH_OP:
			v := self.Tags[op.PushTag()]
			self.Stack.Push(v.t, v.d)
			*pc++
		case PUSH_INT_OP:
			self.Stack.Push(&self.AbcLib.IntType, op.PushIntVal())
			*pc++
		case STOP_OP:
			*pc++
			return nil
		case TRACE_OP:
			*pc++
			self.Code[*pc].Trace(self, *pc, pos, self.Stdout)
		case RET_OP:
			i := len(self.Calls)
			c := self.Calls[i-1]
			self.Calls = self.Calls[:i-1]
			*pc = c.retPc
		default:
			panic(fmt.Sprintf("Invalid op id: %v", id))
		}
	}
}

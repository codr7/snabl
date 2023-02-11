package snabl

import (
	"fmt"
	"io"
	"os"
	"time"
)

const (
	VERSION = 1
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
	Stack Deque[V]
	Calls Deque[Call]
	
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
		case ARG_OP:
			v := self.Calls.Top(0).args[op.ArgIndex()]
			self.Stack.Push(v)
			*pc++
		case BENCH_OP:
			*pc++
			startPc := *pc
			startTime := time.Now()
			
			for i := 0; i < op.BenchReps(); i++ {
				*pc = startPc
				
				if err := self.Eval(pc); err != nil {
					return err
				}

				self.Stack.Clear()
			}
			
			self.Stack.Push(V{t: &self.AbcLib.IntType, d: time.Now().Sub(startTime)})
		case CALL_FUN_OP:
			f := self.Tags[op.CallFunTag()].d.(*Fun)
			
			self.Calls.Push(Call{
				pos: pos, fun: f, args: self.Stack.Drop(f.Arity()), retPc: *pc+1})
			
			*pc = f.pc
		case CALL_PRIM_OP:
			p := self.Tags[op.CallPrimTag()].d.(*Prim)
			
			if err := p.Call(self, pos); err != nil {
				return err
			}

			*pc++
		case GOTO_OP:
			*pc = op.GotoPc()
		case IF_OP:
			if self.Stack.Pop().Bool() {
				*pc++
			} else {
				*pc = op.IfElsePc()
			}
		case POS_OP:
			p := self.Tags[op.PosTag()].Data().(Pos)
			pos = &p
			*pc++
		case PUSH_OP:
			self.Stack.Push(self.Tags[op.PushTag()])
			*pc++
		case PUSH_BOOL_OP:
			self.Stack.Push(V{t: &self.AbcLib.BoolType, d: op.PushBoolVal()})
			*pc++
		case PUSH_INT_OP:
			self.Stack.Push(V{t: &self.AbcLib.IntType, d: op.PushIntVal()})
			*pc++
		case STOP_OP:
			*pc++
			return nil
		case TRACE_OP:
			*pc++
			self.Code[*pc].Trace(self, *pc, pos, self.Stdout)
		case RET_OP:
			*pc = self.Calls.Pop().retPc
		default:
			panic(fmt.Sprintf("Invalid op id: %v", id))
		}
	}
}

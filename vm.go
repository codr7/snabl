package snabl

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

const (
	VERSION = 2
)

type Pc = int
type Tag = int

type Vm struct {
	Debug bool
	Trace bool
	
	Path string
	Stdin io.Reader
	Stdout io.Writer

	AbcLib AbcLib
	
	Tags []V
	
	Code []Op
	Stack Slice[V]
	Calls Slice[Call]
	
	env Env
	fun *Fun
}

func (self *Vm) Init() {
	self.Stdin = os.Stdin
	self.Stdout = os.Stdout
	self.AbcLib.Init(self)
	self.env = NewEnv()
}

func (self *Vm) Load(path string, eval bool) error {
	var p string

	if filepath.IsAbs(path) {
		p = path
	} else {
		p = filepath.Join(self.Path, path)
	}
	
	f, err := os.Open(p)

	if err != nil {
		return err
	}

	defer f.Close()
	var forms Forms

	pos := NewPos(p, 1, 1)
	
	if err := ReadForms(self, pos, bufio.NewReader(f), &forms); err != nil {
		return err
	}
	
	pc := self.EmitPc()

	if err := forms.Emit(self, self.env); err != nil {
		return err
	}
	
	if !eval {
		return nil
	}

	self.Code[self.Emit()] = StopOp()
	prevPath := self.Path

	defer func() {
		self.Path = prevPath
	}()
	
	self.Path = filepath.Dir(p)

	if err := self.Eval(&pc); err != nil {
		return err
	}

	return nil
}

func (self *Vm) Tag(val V) Tag {
	i := len(self.Tags)
	self.Tags = append(self.Tags, val)
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
	tag := self.Tag(V{t: &self.AbcLib.PosType, d: pos})
	self.Code[self.EmitNoTrace()] = PosOp(tag)
}

func (self *Vm) EmitString(str string) {
	tag := self.Tag(V{t: &self.AbcLib.StringType, d: str})
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
			
			self.Stack.Push(V{t: &self.AbcLib.TimeType, d: time.Now().Sub(startTime)})
		case CALL_FUN_OP:
			f := self.Tags[op.CallFunTag()].d.(*Fun)
			
			self.Calls.Push(Call{
				pos: pos, fun: f, args: self.Stack.Tail(f.Arity()), retPc: *pc+1})
			
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
			self.Stack.Push(self.Tags[op.PushVal()])
			*pc++
		case PUSH_BOOL_OP:
			self.Stack.Push(V{t: &self.AbcLib.BoolType, d: op.PushBoolVal()})
			*pc++
		case PUSH_INT_OP:
			self.Stack.Push(V{t: &self.AbcLib.IntType, d: op.PushIntVal()})
			*pc++
		case PUSH_NIL_OP:
			self.Stack.Push(V{t: &self.AbcLib.NilType, d: nil})
			*pc++
		case PUSH_TIME_OP:
			self.Stack.Push(V{t: &self.AbcLib.TimeType, d: op.PushTimeVal()})
			*pc++
		case STOP_OP:
			*pc++
			return nil
		case TEST_OP:
			expected := self.Stack.Pop()
			fmt.Fprintf(self.Stdout, "Testing %v...", expected.String())
			*pc++

			if err := self.Eval(pc); err != nil {
				return err
			}

			if actual := self.Stack.Pop(); actual.Eq(*expected) {
				fmt.Fprintln(self.Stdout, "OK")
			} else {
				fmt.Fprintln(self.Stdout, "FAIL")
			}
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

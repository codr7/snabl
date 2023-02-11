package snabl

import (
	"fmt"
	"time"
)

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
			startTime := time.Now()
			var benchPc Pc

			for i := 0; i < op.BenchReps(); i++ {
				benchPc = *pc+1
				
				if err := self.Eval(&benchPc); err != nil {
					return err
				}

				self.Stack.Clear()
			}

			self.Stack.Push(V{t: &self.AbcLib.TimeType, d: time.Now().Sub(startTime)})
			*pc = benchPc
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
		case CLEAR_OP:
			self.Stack.Clear()
			*pc++
		case DEC_OP:
			v := self.Stack.Top(0)
			v.Init(v.t, v.d.(int) - op.IncDelta())
			*pc++
		case GOTO_OP:
			*pc = op.GotoPc()
		case IF_OP:
			if self.Stack.Pop().Bool() {
				*pc++
			} else {
				*pc = op.IfElsePc()
			}
		case INC_OP:
			v := self.Stack.Top(0)
			v.Init(v.t, v.d.(int) + op.IncDelta())
			*pc++
		case NOP:
			*pc++
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

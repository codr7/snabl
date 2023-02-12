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
		case ADD_INT_OP:
			v := self.Stack.Top(0)
			v.Init(v.t, v.d.(int) + op.AddInt())
			*pc++
		case ARG_OP:
			v := self.call.args[op.ArgIndex()]
			self.Stack.Push(v.t, v.d)
			*pc++
		case BENCH_OP:
			reps := self.Stack.Pop().d.(int)
			startTime := time.Now()
			var benchPc Pc

			for i := 0; i < reps; i++ {
				benchPc = *pc+1
				
				if err := self.Eval(&benchPc); err != nil {
					return err
				}

				self.Stack.Clear()
			}

			self.Stack.Push(&self.AbcLib.TimeType, time.Now().Sub(startTime))
			*pc = benchPc
		case CALL_FUN_OP:
			f := self.Tags[op.CallFun()].d.(*Fun)
			self.call = new(Call).Init(self.call, pos, f, self.Stack.Tail(f.Arity()), *pc+1)
			*pc = f.pc
		case CALL_PRIM_OP:
			p := self.Tags[op.CallPrim()].d.(*Prim)
			
			if err := p.Call(self, pos); err != nil {
				return err
			}

			*pc++
		case EQ_INT_OP:
			v := self.Stack.Top(0)
			v.Init(&self.AbcLib.BoolType, v.d.(int) == op.EqInt())
			*pc++
		case GOTO_OP:
			*pc = op.GotoPc()
		case GT_INT_OP:
			v := self.Stack.Top(0)
			v.Init(&self.AbcLib.BoolType, v.d.(int) > op.GtInt())
			*pc++
		case IF_OP:
			if self.Stack.Pop().Bool() {
				*pc++
			} else {
				*pc = op.IfElsePc()
			}
		case NOP:
			*pc++
		case POS_OP:
			p := self.Tags[op.Pos()].Data().(Pos)
			pos = &p
			*pc++
		case PUSH_OP:
			v := self.Tags[op.PushVal()]
			self.Stack.Push(v.t, v.d)
			*pc++
		case PUSH_BOOL_OP:
			self.Stack.Push(&self.AbcLib.BoolType, op.PushBool())
			*pc++
		case PUSH_INT_OP:
			self.Stack.Push(&self.AbcLib.IntType, op.PushInt())
			*pc++
		case PUSH_NIL_OP:
			self.Stack.Push(&self.AbcLib.NilType, nil)
			*pc++
		case PUSH_TIME_OP:
			self.Stack.Push(&self.AbcLib.TimeType, op.PushTime())
			*pc++
		case REC_OP:
			c := self.call
			f := self.Tags[op.RecFun()].d.(*Fun)
			c.Init(c.parent, pos, f, self.Stack.Tail(f.Arity()), c.retPc)
			*pc = f.pc
		case SLICE_OP:
			prevStack := self.Stack
			self.Stack = new(Stack)
			*pc++
			
			if err := self.Eval(pc); err != nil {
				return err
			}

			prevStack.Push(&self.AbcLib.SliceType, &self.Stack.Slice)
			self.Stack = prevStack
		case STOP_OP:
			*pc++
			return nil
		case SUB_INT_OP:
			v := self.Stack.Top(0)
			v.Init(v.t, v.d.(int) - op.SubInt())
			*pc++
		case TEST_OP:
			expected := self.Stack.Pop()
			expected.Dump(self.Stdout)
			fs := self.Tags[op.TestForm()].d.(*GroupForm)
			
			for _, f := range fs.items {
				fmt.Fprintf(self.Stdout, " %v", f.String())
			}
			
			*pc++

			if err := self.Eval(pc); err != nil {
				return err
			}

			fmt.Fprintln(self.Stdout, "")

			if actual := self.Stack.Pop(); !actual.Eq(*expected) {
				fmt.Fprintf(self.Stdout, "TEST FAILED %v ", fs.pos)
				actual.Dump(self.Stdout)
				fmt.Fprintln(self.Stdout, "")
			}
		case TRACE_OP:
			*pc++
			self.Code[*pc].Trace(self, *pc, pos, true, self.Stdout)
		case RET_OP:
			c := self.call
			self.call = c.parent
			*pc = c.retPc
		default:
			panic(fmt.Sprintf("Invalid op id: %v", id))
		}
	}
}

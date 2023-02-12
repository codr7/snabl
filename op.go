package snabl

import (
	"fmt"
	"io"
	"time"
	"unsafe"
)

type Op uint64
type OpId = int

const (
	OP_WIDTH = int(unsafe.Sizeof(Op(0)) * 8)
	OP_ID_WIDTH = 6

	ARG_INDEX = OP_ID_WIDTH
	ARG_INDEX_WIDTH = OP_WIDTH - ARG_INDEX

	BENCH_REPS = OP_ID_WIDTH
	BENCH_REPS_WIDTH = OP_WIDTH - BENCH_REPS

	CALL_FUN = OP_ID_WIDTH
	CALL_FUN_WIDTH = OP_WIDTH - CALL_FUN

	CALL_PRIM = OP_ID_WIDTH
	CALL_PRIM_WIDTH = OP_WIDTH - CALL_PRIM

	DEC_DELTA = OP_ID_WIDTH
	DEC_DELTA_WIDTH = OP_WIDTH - DEC_DELTA

	GOTO_PC = OP_ID_WIDTH
	GOTO_PC_WIDTH = OP_WIDTH - GOTO_PC

	IF_ELSE_PC = OP_ID_WIDTH
	IF_ELSE_PC_WIDTH = OP_WIDTH - IF_ELSE_PC
	
	INC_DELTA = OP_ID_WIDTH
	INC_DELTA_WIDTH = OP_WIDTH - INC_DELTA

	POS_TAG = OP_ID_WIDTH
	POS_TAG_WIDTH = OP_WIDTH - POS_TAG

	PUSH_VAL = OP_ID_WIDTH
	PUSH_VAL_WIDTH = OP_WIDTH - PUSH_VAL

	PUSH_BOOL_VAL = OP_ID_WIDTH

	PUSH_INT_VAL = OP_ID_WIDTH
	PUSH_INT_VAL_WIDTH = OP_WIDTH - PUSH_INT_VAL

	PUSH_TIME_VAL = OP_ID_WIDTH
	PUSH_TIME_VAL_WIDTH  = OP_WIDTH - PUSH_TIME_VAL

	REC_FUN = OP_ID_WIDTH
	REC_FUN_WIDTH = OP_WIDTH - REC_FUN
	
	ARG_OP = iota
	BENCH_OP
	CALL_FUN_OP
	CALL_PRIM_OP
	CLEAR_OP
	DEC_OP
	GOTO_OP
	IF_OP
	INC_OP
	NOP
	POS_OP
	PUSH_OP
	PUSH_BOOL_OP
	PUSH_INT_OP
	PUSH_NIL_OP
	PUSH_TIME_OP
	REC_OP
	RET_OP
	STOP_OP
	TEST_OP
	TRACE_OP
)

type OpArgType interface {
	int | time.Duration
}

func OpArg[T OpArgType](op Op, pos, width int) T {
	return (T(op) >> pos) & ((T(1) << width) - 1)
}

func (self Op) Id() OpId {
	return OpArg[OpId](self, 0, OP_ID_WIDTH)
}

func (self Op) Trace(vm *Vm, pc Pc, pos *Pos, out io.Writer) {	
	fmt.Fprintf(out, "%v ", pc) 
	
	switch id := self.Id(); id {
	case ARG_OP:
		fmt.Fprintf(out, "ARG %v", self.ArgIndex())
	case BENCH_OP:
		fmt.Fprintf(out, "BENCH %v", self.BenchReps())
	case CALL_FUN_OP:
		fmt.Fprintf(out, "CALL_FUN %v", vm.Tags[self.CallFun()].String())
	case CALL_PRIM_OP:
		fmt.Fprintf(out, "CALL_PRIM %v", vm.Tags[self.CallPrim()].String())
	case CLEAR_OP:
		io.WriteString(out, "CLEAR")
	case DEC_OP:
		fmt.Fprintf(out, "DEC %v", self.DecDelta())
	case GOTO_OP:
		fmt.Fprintf(out, "GOTO %v", self.GotoPc())
	case IF_OP:
		fmt.Fprintf(out, "IF %v", self.IfElsePc())
	case INC_OP:
		fmt.Fprintf(out, "INC %v", self.IncDelta())
	case NOP:
		io.WriteString(out, "NOP")
	case PUSH_OP:
		fmt.Fprintf(out, "PUSH %v", vm.Tags[self.PushVal()].String())
	case PUSH_BOOL_OP:
		fmt.Fprintf(out, "PUSH_BOOL %v", self.PushBoolVal())
	case PUSH_INT_OP:
		fmt.Fprintf(out, "PUSH_INT %v", self.PushIntVal())
	case PUSH_NIL_OP:
		io.WriteString(out, "PUSH_NIL")
	case PUSH_TIME_OP:
		fmt.Fprintf(out, "PUSH_TIME %v", self.PushTimeVal())
	case REC_OP:
		fmt.Fprintf(out, "REC %v", vm.Tags[self.RecFun()].String())
	case RET_OP:
		io.WriteString(out, "RET")
	case STOP_OP:
		io.WriteString(out, "STOP")
	case TEST_OP:
		fmt.Fprintf(out, "TEST")
	default:
		panic(fmt.Sprintf("Invalid op id: %v", id))
	}

	if pos != nil {
		fmt.Fprintf(out, " (%v)", *pos)
	}

	fmt.Fprintf(out, " %v", vm.Stack.String())
	fmt.Fprintln(out, "")
}

func ArgOp(index int) Op {
	return Op(ARG_OP) + Op(index << ARG_INDEX)
}

func (self Op) ArgIndex() int {
	return OpArg[int](self, ARG_INDEX, ARG_INDEX_WIDTH)
}

func BenchOp(reps int) Op {
	return Op(BENCH_OP) + Op(reps << BENCH_REPS)
}

func (self Op) BenchReps() int {
	return OpArg[int](self, BENCH_REPS, BENCH_REPS_WIDTH)
}

func CallFunOp(fun *Fun) Op {
	return Op(CALL_FUN_OP) + Op(fun.tag << CALL_FUN)
}

func (self Op) CallFun() Tag {
	return OpArg[Tag](self, CALL_FUN, CALL_FUN_WIDTH)
}

func CallPrimOp(prim *Prim) Op {
	return Op(CALL_PRIM_OP) + Op(prim.tag << CALL_PRIM)
}

func (self Op) CallPrim() Tag {
	return OpArg[Tag](self, CALL_PRIM, CALL_PRIM_WIDTH)
}

func ClearOp() Op {
	return Op(CLEAR_OP)
}

func DecOp(delta int) Op {
	return Op(DEC_OP) + Op(delta << DEC_DELTA)
}

func (self Op) DecDelta() int {
	return OpArg[int](self, DEC_DELTA, DEC_DELTA_WIDTH)
}

func GotoOp(pc Pc) Op {
	return Op(GOTO_OP) + Op(pc << GOTO_PC)
}

func (self Op) GotoPc() Pc {
	return OpArg[Pc](self, GOTO_PC, GOTO_PC_WIDTH)
}

func IfOp(elsePc Pc) Op {
	return Op(IF_OP) + Op(elsePc << IF_ELSE_PC)
}

func IncOp(delta int) Op {
	return Op(INC_OP) + Op(delta << INC_DELTA)
}

func (self Op) IncDelta() int {
	return OpArg[int](self, INC_DELTA, INC_DELTA_WIDTH)
}

func (self Op) IfElsePc() Pc {
	return OpArg[Pc](self, IF_ELSE_PC, IF_ELSE_PC_WIDTH)
}

func NOp() Op {
	return Op(NOP)
}

func PosOp(tag Tag) Op {
	return Op(POS_OP) + Op(tag << POS_TAG)
}

func (self Op) PosTag() Tag {
	return OpArg[Tag](self, POS_TAG, POS_TAG_WIDTH)
}

func PushOp(val Tag) Op {
	return Op(PUSH_OP) + Op(val << PUSH_VAL)
}

func (self Op) PushVal() Tag {
	return OpArg[Tag](self, PUSH_VAL, PUSH_VAL_WIDTH)
}

func PushBoolOp(val bool) Op {
	var v int
	
	if val {
		v = 1
	} else {
		v = 0
	}
	
	return Op(PUSH_BOOL_OP) + Op(v << PUSH_BOOL_VAL)
}

func (self Op) PushBoolVal() bool {
	return OpArg[int](self, PUSH_BOOL_VAL, 1) == 1
}

func PushIntOp(val int) Op {
	return Op(PUSH_INT_OP) + Op(val << PUSH_INT_VAL)
}

func (self Op) PushIntVal() int {
	return OpArg[int](self, PUSH_INT_VAL, PUSH_INT_VAL_WIDTH)
}

func PushNilOp() Op {
	return Op(PUSH_NIL_OP)
}

func PushTimeOp(val time.Duration) Op {
	return Op(PUSH_TIME_OP) + Op(val << PUSH_TIME_VAL)
}

func (self Op) PushTimeVal() time.Duration {
	return OpArg[time.Duration](self, PUSH_TIME_VAL, PUSH_TIME_VAL_WIDTH)
}

func RecOp(fun *Fun) Op {
	return Op(REC_OP) + Op(fun.tag << REC_FUN)
}

func (self Op) RecFun() Tag {
	return OpArg[Tag](self, REC_FUN, REC_FUN_WIDTH)
}

func RetOp() Op {
	return Op(RET_OP)
}

func StopOp() Op {
	return Op(STOP_OP)
}

func TestOp() Op {
	return Op(TEST_OP)
}

func TraceOp() Op {
	return Op(TRACE_OP)
}


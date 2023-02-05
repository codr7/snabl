package snabl

import (
	"fmt"
	"io"
	"unsafe"
)

type Op uint64
type OpId = uint8

const (
	OP_WIDTH = uint8((unsafe.Sizeof(Op(0)) * 8))
	OP_ID_WIDTH = 6

	ARG_TAG = OP_ID_WIDTH
	ARG_INDEX_WIDTH = 3
	ARG_TAG_WIDTH = OP_WIDTH - ARG_TAG - ARG_INDEX_WIDTH
	ARG_INDEX = ARG_TAG + ARG_TAG_WIDTH

	ARG_OFFS_TAG = OP_ID_WIDTH
	ARG_OFFS_TAG_WIDTH = OP_WIDTH - ARG_OFFS_TAG

	CALL_PRIM_TAG = OP_ID_WIDTH
	CALL_PRIM_TAG_WIDTH = OP_WIDTH - CALL_PRIM_TAG

	POS_TAG = OP_ID_WIDTH
	POS_TAG_WIDTH = OP_WIDTH - POS_TAG

	PUSH_INT_VAL = OP_ID_WIDTH
	PUSH_INT_VAL_WIDTH = OP_WIDTH - PUSH_INT_VAL

	PUSH_TAG = OP_ID_WIDTH
	PUSH_TAG_WIDTH = OP_WIDTH - PUSH_TAG

	ADD_OP = iota
	ARG_OP
	ARG_OFFS_OP
	CALL_PRIM_OP
	POS_OP
	PUSH_OP
	PUSH_INT_OP
	STOP_OP
	TRACE_OP
)

type OpArgType interface {
	uint8 | uint | int 
}

func OpArg[T OpArgType](op Op, pos, width uint8) T {
	return (T(op) >> pos) & ((T(1) << width) - 1)
}

func (self Op) Id() OpId {
	return OpArg[OpId](self, 0, OP_ID_WIDTH)
}

func (self Op) Trace(vm *Vm, pc Pc, pos *Pos, out io.Writer) {	
	fmt.Fprintf(out, "%v ", pc) 
	
	switch id := self.Id(); id {
	case ADD_OP:
		io.WriteString(out, "ADD")
	case ARG_OP:
		fmt.Fprintf(out, "ARG %v %v", self.ArgTag(), self.ArgIndex())
	case ARG_OFFS_OP:
		fmt.Fprintf(out, "ARG_OFFS %v", self.ArgOffsTag())
	case CALL_PRIM_OP:
		fmt.Fprintf(out, "CALL_PRIM %v", vm.Tags[self.CallPrimTag()].String())
	case PUSH_OP:
		fmt.Fprintf(out, "PUSH %v", vm.Tags[self.PushTag()].String())
	case PUSH_INT_OP:
		fmt.Fprintf(out, "PUSH_INT %v", self.PushIntVal())
	case STOP_OP:
		io.WriteString(out, "STOP")
	default:
		panic(fmt.Sprintf("Invalid op id: %v", id))
	}

	if pos != nil {
		fmt.Fprintf(out, " (%v)", *pos)
	}

	fmt.Fprintln(out, "")
}

func AddOp() Op {
	return Op(ADD_OP)
}

func ArgOp(tag Tag, index ArgIndex) Op {
	return Op(ARG_OP) + Op(tag << ARG_TAG) + Op(index << ARG_INDEX)
}

func (self Op) ArgTag() Tag {
	return OpArg[Tag](self, ARG_TAG, ARG_TAG_WIDTH)
}

func (self Op) ArgIndex() ArgIndex {
	return OpArg[ArgIndex](self, ARG_INDEX, ARG_INDEX_WIDTH)
}

func ArgOffsOp(tag Tag) Op {
	return Op(ARG_OFFS_OP) + Op(tag << ARG_OFFS_TAG)
}

func (self Op) ArgOffsTag() Tag {
	return OpArg[Tag](self, ARG_OFFS_TAG, ARG_OFFS_TAG_WIDTH)
}

func CallPrimOp(tag Tag) Op {
	return Op(CALL_PRIM_OP) + Op(tag << CALL_PRIM_TAG)
}

func (self Op) CallPrimTag() Tag {
	return OpArg[Tag](self, CALL_PRIM_TAG, CALL_PRIM_TAG_WIDTH)
}

func PosOp(tag Tag) Op {
	return Op(POS_OP) + Op(tag << POS_TAG)
}

func (self Op) PosTag() Tag {
	return OpArg[Tag](self, POS_TAG, POS_TAG_WIDTH)
}

func PushOp(tag Tag) Op {
	return Op(PUSH_OP) + Op(tag << PUSH_TAG)
}

func (self Op) PushTag() Tag {
	return OpArg[Tag](self, PUSH_TAG, PUSH_TAG_WIDTH)
}

func PushIntOp(val int) Op {
	return Op(PUSH_INT_OP) + Op(val << PUSH_INT_VAL)
}

func (self Op) PushIntVal() int {
	return OpArg[int](self, PUSH_INT_VAL, PUSH_INT_VAL_WIDTH)
}

func StopOp() Op {
	return Op(STOP_OP)
}

func TraceOp() Op {
	return Op(TRACE_OP)
}


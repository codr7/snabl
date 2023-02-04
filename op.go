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

	POS_TAG = OP_ID_WIDTH
	POS_TAG_WIDTH = OP_WIDTH - POS_TAG

	PUSH_INT_VAL = OP_ID_WIDTH
	PUSH_INT_VAL_WIDTH = OP_WIDTH - PUSH_INT_VAL
		
	ADD_OP = iota
	POS_OP
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

func (self Op) Trace(pc Pc, pos *Pos, out io.Writer) {	
	fmt.Fprintf(out, "%v ", pc) 
	
	switch id := self.Id(); id {
	case ADD_OP:
		io.WriteString(out, "ADD")
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

func PosOp(tag Tag) Op {
	return Op(POS_OP) + Op(tag << POS_TAG)
}

func (self Op) PosTag() Tag {
	return OpArg[Tag](self, POS_TAG, POS_TAG_WIDTH)
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


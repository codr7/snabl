package snabl

import (
	"fmt"
	"io"
	"unsafe"
)

type Op uint64
type OpId = uint8

const (
	OP_ID_WIDTH = 6

	PUSH_INT_VAL_POS = OP_ID_WIDTH
	PUSH_INT_VAL_WIDTH = uint8(unsafe.Sizeof(Op(0)) * 8)
		
	ADD_OP = iota
	PUSH_INT_OP
	STOP_OP
	TRACE_OP
)

type OpArgType interface {
	uint8 | int 
}

func OpArg[T OpArgType](op Op, pos, width uint8) T {
	return (T(op) >> pos) & ((T(1) << width) - 1)
}

func (self Op) Id() OpId {
	return OpArg[OpId](self, 0, OP_ID_WIDTH)
}

func (self Op) Trace(pc Pc, out io.Writer) {
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

	fmt.Fprintln(out, "")
}

func AddOp() Op {
	return Op(ADD_OP)
}

func PushIntOp(val int) Op {
	return Op(PUSH_INT_OP) + Op(val << PUSH_INT_VAL_POS)
}

func (self Op) PushIntVal() int {
	return OpArg[int](self, PUSH_INT_VAL_POS, PUSH_INT_VAL_WIDTH)
}

func StopOp() Op {
	return Op(STOP_OP)
}

func TraceOp() Op {
	return Op(TRACE_OP)
}


package snabl

import (
	"fmt"
	"io"
	"unsafe"
)

type Op = uint64

type OpIdT = uint8

const (
	OP_ID_WIDTH = 6

	PUSH_INT_VAL_POS = OP_ID_WIDTH
	PUSH_INT_VAL_WIDTH = uint8(unsafe.Sizeof(Op(0)) * 8)
		
	ADD_OP = iota
	PUSH_INT_OP
	STOP_OP
)

type OpArgT interface {
	uint8 | int 
}

func OpArg[T OpArgT](op Op, pos, width uint8) T {
	return (T(op) >> pos) & ((T(1) << width) - 1)
}

func OpId(op Op) OpIdT {
	return OpArg[OpIdT](op, 0, OP_ID_WIDTH)
}

func OpTrace(pc Pc, op Op, out io.Writer) {
	fmt.Fprintf(out, "%v ", pc) 

	switch id := OpId(op); id {
	case ADD_OP:
		io.WriteString(out, "ADD")
	case PUSH_INT_OP:
		io.WriteString(out, "PUSH_INT")
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

func PushIntVal(op Op) int {
	return OpArg[int](op, PUSH_INT_VAL_POS, PUSH_INT_VAL_WIDTH)
}

func StopOp() Op {
	return Op(STOP_OP)
}


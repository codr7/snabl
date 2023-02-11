package snabl

import (
	"fmt"
)

type Fuse = func (pc Pc) int

func (self *Vm) Fuse(startPc Pc) {
	for self.FuseDec(startPc) > 0 {}
}

func (self *Vm) FuseDec(startPc Pc) int {
	count := 0
	var prevOp *Op
	
	for pc := startPc; pc < len(self.Code); pc++ {
		op := &self.Code[pc]

		if prevOp != nil &&
			prevOp.Id() == PUSH_INT_OP &&
			op.Id() == CALL_PRIM_OP &&
			op.CallPrimTag() == self.AbcLib.SubPrim.tag {
			fmt.Fprintf(self.Stdout, "%v Fusing PUSH_INT CALL_PRIM -> DEC\n", pc);
			*prevOp = DecOp(prevOp.PushIntVal())
			*op = NOp()
			count++
		}
		
		prevOp = op
	}

	return count
}

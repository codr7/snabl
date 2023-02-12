package snabl

import (
	"fmt"
)

func (self *Vm) Fuse(startPc Pc) {
	for self.FuseAddInt(startPc, nil) +
		self.FuseArg2(startPc, nil) +
		self.FuseArg3(startPc, nil) +
		self.FuseEqInt(startPc, nil) +		
		self.FuseGoto(startPc, nil) +
		self.FuseGtInt(startPc, nil) +
		self.FuseNop(startPc, nil) +
		self.FuseRec(startPc, nil) +
		self.FuseSubInt(startPc, nil) > 0 {}
}

func (self *Vm) FuseAddInt(startPc Pc, prevOp *Op) int {
	count := 0
	
	for pc := startPc; pc < len(self.Code); {
		op := &self.Code[pc]

		switch op.Id() {
		case GOTO_OP:
			count += self.FuseAddInt(pc+1, prevOp)
			pc = op.GotoPc()
			continue
		case NOP, POS_OP, TRACE_OP:
			pc++
			continue
		}

		if prevOp != nil &&
			prevOp.Id() == PUSH_INT_OP &&
			op.Id() == CALL_PRIM_OP &&
			op.CallPrim() == self.AbcLib.AddPrim.tag {
			fmt.Fprintf(self.Stdout, "Fusing %v ADD_INT\n", pc);
			*op = AddIntOp(prevOp.PushInt())
			*prevOp = NOp()
			count++
		}
		
		prevOp = op
		pc++
	}

	return count
}

func (self *Vm) FuseArg2(startPc Pc, prevOp *Op) int {
	count := 0
	
	for pc := startPc; pc < len(self.Code); {
		op := &self.Code[pc]

		switch op.Id() {
		case NOP, POS_OP, TRACE_OP:
			pc++
			continue
		}

		if prevOp != nil && prevOp.Id() == ARG_OP && op.Id() == ARG_OP {
			fmt.Fprintf(self.Stdout, "Fusing %v ARG2\n", pc);
			*op = Arg2Op(prevOp.ArgIndex(), op.ArgIndex())
			*prevOp = NOp()
			count++
		}

		prevOp = op
		pc++
	}

	return count
}

func (self *Vm) FuseArg3(startPc Pc, prevOp *Op) int {
	count := 0
	
	for pc := startPc; pc < len(self.Code); {
		op := &self.Code[pc]

		switch op.Id() {
		case NOP, POS_OP, TRACE_OP:
			pc++
			continue
		}

		if prevOp != nil && prevOp.Id() == ARG2_OP && op.Id() == ARG_OP {
			fmt.Fprintf(self.Stdout, "Fusing %v ARG3\n", pc);
			*op = Arg3Op(prevOp.Arg2Index1(), prevOp.Arg2Index2(), op.ArgIndex())
			*prevOp = NOp()
			count++
		}

		prevOp = op
		pc++
	}

	return count
}

func (self *Vm) FuseEqInt(startPc Pc, prevOp *Op) int {
	count := 0
	
	for pc := startPc; pc < len(self.Code); {
		op := &self.Code[pc]

		switch op.Id() {
		case GOTO_OP:
			count += self.FuseEqInt(pc+1, prevOp)
			pc = op.GotoPc()
			continue
		case NOP, POS_OP, TRACE_OP:
			pc++
			continue
		}

		if prevOp != nil &&
			prevOp.Id() == PUSH_INT_OP &&
			op.Id() == CALL_PRIM_OP &&
			op.CallPrim() == self.AbcLib.EqPrim.tag {
			fmt.Fprintf(self.Stdout, "Fusing %v EQ_INT\n", pc);
			*op = EqIntOp(prevOp.PushInt())
			*prevOp = NOp()
			count++
		}
		
		prevOp = op
		pc++
	}

	return count
}

func (self *Vm) FuseGoto(startPc Pc, prevOp *Op) int {
	count := 0
	
	for pc := startPc; pc < len(self.Code); {
		op := &self.Code[pc]

		switch op.Id() {
		case NOP, POS_OP, TRACE_OP:
			pc++
			continue
		}

		if prevOp != nil &&
			prevOp.Id() == GOTO_OP &&
			(op.Id() == GOTO_OP || op.Id() == RET_OP || op.Id() == STOP_OP) {
			fmt.Fprintf(self.Stdout, "Fusing %v GOTO\n", pc);
			*prevOp = *op
			count++
		}
		
		prevOp = op

		if op.Id() == GOTO_OP {
			count += self.FuseGoto(pc+1, nil)
			pc = op.GotoPc()
		} else {
			pc++
		}
	}

	return count
}

func (self *Vm) FuseGtInt(startPc Pc, prevOp *Op) int {
	count := 0
	
	for pc := startPc; pc < len(self.Code); {
		op := &self.Code[pc]

		switch op.Id() {
		case GOTO_OP:
			count += self.FuseGtInt(pc+1, prevOp)
			pc = op.GotoPc()
			continue
		case NOP, POS_OP, TRACE_OP:
			pc++
			continue
		}

		if prevOp != nil &&
			prevOp.Id() == PUSH_INT_OP &&
			op.Id() == CALL_PRIM_OP &&
			op.CallPrim() == self.AbcLib.GtPrim.tag {
			fmt.Fprintf(self.Stdout, "Fusing %v GT_INT\n", pc);
			*op = GtIntOp(prevOp.PushInt())
			*prevOp = NOp()
			count++
		}
		
		prevOp = op
		pc++
	}

	return count
}

func (self *Vm) FuseNop(startPc Pc, prevOp *Op) int {
	count := 0
	
	for pc := startPc; pc < len(self.Code); pc++ {
		op := &self.Code[pc]

		switch op.Id() {
		case POS_OP, TRACE_OP:
			pc++
			continue
		}

		if prevOp != nil &&
			((prevOp.Id() == GOTO_OP && prevOp.GotoPc() < pc+1) || prevOp.Id() == NOP) &&
			op.Id() == NOP {
			fmt.Fprintf(self.Stdout, "Fusing %v NOP\n", pc);
			*prevOp = GotoOp(pc+1)
			count++
		}

		prevOp = op
	}

	return count
}

func (self *Vm) FuseRec(startPc Pc, prevOp *Op) int {
	count := 0
	
	for pc := startPc; pc < len(self.Code); {
		op := &self.Code[pc]

		switch op.Id() {
		case GOTO_OP:
			count += self.FuseRec(pc+1, prevOp)
			pc = op.GotoPc()
			continue
		case NOP, POS_OP, TRACE_OP:
			pc++
			continue
		}

		if prevOp != nil && prevOp.Id() == CALL_FUN_OP && op.Id() == RET_OP {
			fmt.Fprintf(self.Stdout, "Fusing %v CALL_FUN RET\n", pc);
			*prevOp = RecOp(self.Tags[prevOp.CallFun()].d.(*Fun))
			count++
		}
		
		prevOp = op
		pc++
	}

	return count
}

func (self *Vm) FuseSubInt(startPc Pc, prevOp *Op) int {
	count := 0
	
	for pc := startPc; pc < len(self.Code); {
		op := &self.Code[pc]

		switch op.Id() {
		case GOTO_OP:
			count += self.FuseSubInt(pc+1, prevOp)
			pc = op.GotoPc()
			continue
		case NOP, POS_OP, TRACE_OP:
			pc++
			continue
		}

		if prevOp != nil &&
			prevOp.Id() == PUSH_INT_OP &&
			op.Id() == CALL_PRIM_OP &&
			op.CallPrim() == self.AbcLib.SubPrim.tag {
			fmt.Fprintf(self.Stdout, "Fusing %v SUB_INT\n", pc);
			*op = SubIntOp(prevOp.PushInt())
			*prevOp = NOp()
			count++
		}
		
		prevOp = op
		pc++
	}

	return count
}

package snabl

import (
	"fmt"
)

func (self *Vm) Fuse(startPc Pc) {
	for self.FuseAddInt(startPc, nil) +
		self.FuseGoto(startPc, nil) +
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
		case NOP:
			pc++
			continue
		}

		if prevOp != nil &&
			prevOp.Id() == PUSH_INT_OP &&
			op.Id() == CALL_PRIM_OP &&
			op.CallPrim() == self.AbcLib.AddPrim.tag {
			fmt.Fprintf(self.Stdout, "%v Fusing PUSH_INT CALL_PRIM(+)\n", pc);
			*op = AddIntOp(prevOp.PushInt())
			*prevOp = NOp()
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
		case NOP:
			pc++
			continue
		}

		if prevOp != nil &&
			prevOp.Id() == PUSH_INT_OP &&
			op.Id() == CALL_PRIM_OP &&
			op.CallPrim() == self.AbcLib.SubPrim.tag {
			fmt.Fprintf(self.Stdout, "%v Fusing PUSH_INT CALL_PRIM(-)\n", pc);
			*op = SubIntOp(prevOp.PushInt())
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

		if  op.Id() == NOP {
			pc++
			continue
		}

		if prevOp != nil &&
			prevOp.Id() == GOTO_OP &&
			(op.Id() == GOTO_OP || op.Id() == RET_OP || op.Id() == STOP_OP) {
			fmt.Fprintf(self.Stdout, "%v Fusing GOTO\n", pc);
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

func (self *Vm) FuseNop(startPc Pc, prevOp *Op) int {
	count := 0
	
	for pc := startPc; pc < len(self.Code); pc++ {
		op := &self.Code[pc]

		if prevOp != nil && (prevOp.Id() == GOTO_OP || prevOp.Id() == NOP) && op.Id() == NOP {
			fmt.Fprintf(self.Stdout, "%v Fusing NOP\n", pc);
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
		case NOP:
			pc++
			continue
		}

		if prevOp != nil && prevOp.Id() == CALL_FUN_OP && op.Id() == RET_OP {
			fmt.Fprintf(self.Stdout, "%v Fusing CALL_FUN RET\n", pc);
			*prevOp = RecOp(self.Tags[prevOp.CallFun()].d.(*Fun))
			count++
		}
		
		prevOp = op
		pc++
	}

	return count
}

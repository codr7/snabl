package tests

import (
	"testing"
	"github.com/codr7/snabl"
)

func TestEval(t *testing.T) {
	var vm snabl.Vm
	vm.Init()
	
	pc := vm.EmitPc()
	vm.Code[vm.Emit(1)] = snabl.PushIntOp(35) 
	vm.Code[vm.Emit(1)] = snabl.PushIntOp(7) 
	vm.Code[vm.Emit(1)] = snabl.AddOp() 
	vm.Code[vm.Emit(1)] = snabl.StopOp()
	
	if err := vm.Eval(&pc); err != nil {
		t.Fatal(err)
	}

	if vm.Stack.Len() != 1 {
		t.Fatalf("Expected one item: %v",  vm.Stack.String())
	}

	if v := vm.Stack.Pop(); v.Type() != &vm.AbcLib.IntType || v.Data().(int) != 42 {
		t.Errorf("Expected 42: %v",  v.String())
	}
}

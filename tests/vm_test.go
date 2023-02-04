package tests

import (
	"testing"
	"github.com/codr7/snabl"
)

func NewVm() *snabl.Vm {
	var vm snabl.Vm
	vm.Init()
	vm.Debug = true
	vm.Trace = true
	return &vm
}

func TestEval(t *testing.T) {
	vm := NewVm()
	pc := vm.EmitPc()
	vm.Code[vm.Emit()] = snabl.PushIntOp(35) 
	vm.Code[vm.Emit()] = snabl.PushIntOp(7) 
	vm.Code[vm.Emit()] = snabl.AddOp() 
	vm.Code[vm.Emit()] = snabl.StopOp()
	
	if err := vm.Eval(&pc); err != nil {
		t.Fatal(err)
	}

	if vm.Stack.Len() != 1 {
		t.Fatalf("Expected one item: %v",  vm.Stack.String())
	}

	if v := vm.Stack.Pop(); v.Type() != &snabl.Abc.IntType || v.Data().(int) != 42 {
		t.Errorf("Expected 42: %v",  v.String())
	}
}

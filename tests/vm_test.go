package tests

import (
	"testing"
	"github.com/codr7/snabl"
)

func NewVm() *snabl.Vm {
	var vm snabl.Vm
	vm.Init()
	//vm.Debug = true
	vm.Trace = true
	return &vm
}

func TestAdd(t *testing.T) {
	vm := NewVm()
	pc := vm.EmitPc()
	vm.Code[vm.Emit()] = snabl.PushIntOp(35) 
	vm.Code[vm.Emit()] = snabl.PushIntOp(7)
	vm.EmitPrim(&snabl.Abc.AddPrim)
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

func TestFail(t *testing.T) {
	vm := NewVm()
	pc := vm.EmitPc()
	pos := snabl.NewPos("test", 7, 14)
	vm.EmitPos(*pos)
	msg := "failing"
	vm.EmitString(msg)
	vm.EmitPrim(&snabl.Abc.FailPrim)
	vm.Code[vm.Emit()] = snabl.StopOp()
	
	if err := vm.Eval(&pc); err == nil {
		t.Fatal("Should fail with error")
	} else if e := err.(*snabl.E);
	e.Pos().Source() != pos.Source() || e.Pos().Line() != pos.Line() || e.Msg() != msg {
		t.Fatalf("Wrong information in error: %v", e.Error())
	}
	
	if vm.Stack.Len() != 0 {
		t.Fatalf("Expected empty stack: %v",  vm.Stack.String())
	}
}

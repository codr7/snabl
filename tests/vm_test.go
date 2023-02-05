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

func PopInt(t *testing.T, vm *snabl.Vm, expected int) {
	if actual := vm.Stack.Pop(); actual.Type() != &snabl.Abc.IntType || actual.Data().(int) != expected {
		t.Errorf("Expected %v: %v", expected, actual.String())
	}
}

func TestAdd(t *testing.T) {
	vm := NewVm()
	pc := vm.EmitPc()
	vm.Code[vm.Emit()] = snabl.PushIntOp(7) 
	vm.Code[vm.Emit()] = snabl.PushIntOp(14)
	vm.EmitPrim(&snabl.Abc.AddPrim)
	vm.Code[vm.Emit()] = snabl.StopOp()
	
	if err := vm.Eval(&pc); err != nil {
		t.Fatal(err)
	}

	if vm.Stack.Len() != 1 {
		t.Fatalf("Expected [21]: %v",  vm.Stack.String())
	}

	PopInt(t, vm, 21)
}

func TestArgs(t *testing.T) {
	vm := NewVm()
	pc := vm.EmitPc()
	vm.Code[vm.Emit()] = snabl.PushIntOp(7)
	vm.Code[vm.Emit()] = snabl.PushIntOp(14)
	argOffsTag := vm.Tag(&snabl.Abc.IntType, 0)
	vm.Code[vm.Emit()] = snabl.ArgOffsOp(argOffsTag) 
	vm.Code[vm.Emit()] = snabl.ArgOp(argOffsTag, 0) 
	vm.Code[vm.Emit()] = snabl.ArgOp(argOffsTag, 1) 
	vm.Code[vm.Emit()] = snabl.StopOp()
	
	if err := vm.Eval(&pc); err != nil {
		t.Fatal(err)
	}

	if v := vm.Tags[argOffsTag].Data().(int); v != 2 {
		t.Fatalf("Expected arg offset 2: %v", v) 
	}

	if vm.Stack.Len() != 4 {
		t.Fatalf("Expected [7, 14, 7, 14]: %v",  vm.Stack.String())
	}

	PopInt(t, vm, 7)
	PopInt(t, vm, 14)
}

func TestFail(t *testing.T) {
	vm := NewVm()
	vm.Debug = false
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
		t.Fatalf("Expected []: %v",  vm.Stack.String())
	}
}

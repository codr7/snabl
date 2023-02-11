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

func EmitPrim(vm *snabl.Vm, prim *snabl.Prim) {
	tag := vm.Tag(&vm.AbcLib.PrimType, prim)
	vm.Code[vm.Emit()] = snabl.CallPrimOp(tag) 
}

func PopInt(t *testing.T, vm *snabl.Vm, expected int) {
	if actual := vm.Stack.Pop(); actual.Type() != &vm.AbcLib.IntType || actual.Data().(int) != expected {
		t.Errorf("Expected %v: %v", expected, actual.String())
	}
}

func TestAdd(t *testing.T) {
	vm := NewVm()
	pc := vm.EmitPc()
	vm.Code[vm.Emit()] = snabl.PushIntOp(7) 
	vm.Code[vm.Emit()] = snabl.PushIntOp(14)
	EmitPrim(vm, &vm.AbcLib.AddPrim)
	vm.Code[vm.Emit()] = snabl.StopOp()
	
	if err := vm.Eval(&pc); err != nil {
		t.Fatal(err)
	}

	if vm.Stack.Len() != 1 {
		t.Fatalf("Expected [21]: %v",  vm.Stack.String())
	}

	PopInt(t, vm, 21)
}

func TestFail(t *testing.T) {
	vm := NewVm()
	vm.Debug = false
	pc := vm.EmitPc()
	pos := snabl.NewPos("test", 7, 14)
	vm.EmitPos(*pos)
	msg := "failing"
	vm.EmitString(msg)
	EmitPrim(vm, &vm.AbcLib.FailPrim)
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

func TestFun(t *testing.T) {
	vm := NewVm()
	pc := vm.EmitPc()

	pos := snabl.NewPos("TestCall", 1, 1)
	var args snabl.Forms
	fun := snabl.NewIdForm(*pos, "foo")
	arg := snabl.NewIdForm(*pos, "x")
	args.Push(fun)
	args.Push(snabl.NewGroupForm(*pos, arg))
	args.Push(arg)

	if err := vm.AbcLib.FunMacro.Emit(&args, vm, vm.Env(), *pos); err != nil {
		t.Fatal(err)
	}

	args.Push(snabl.NewLitForm(*pos, &vm.AbcLib.IntType, 42))
	
	if err := fun.Emit(&args, vm, vm.Env()); err != nil {
		t.Fatal(err)
	}

	vm.Code[vm.Emit()] = snabl.StopOp()
	
	if err := vm.Eval(&pc); err != nil {
		t.Fatal(err)
	}

	PopInt(t, vm, 42)
}

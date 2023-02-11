package snabl

import (
	"fmt"
	"io"
)

type AbcLib struct {
	BasicLib
	BoolType BoolType
	FunType FunType
	IntType IntType
	MacroType MacroType
	MetaType MetaType
	PosType PosType
	PrimType PrimType
	StringType StringType

	BenchMacro, FunMacro Macro
	
	AddPrim, FailPrim, GtPrim, SubPrim Prim
}

func (self *AbcLib) Init(vm *Vm) {
	self.BasicLib.Init(vm, "abc")
	
	self.BindType(&self.BoolType, "Bool")
	self.BindType(&self.FunType, "Fun")
	self.BindType(&self.IntType, "Int")
	self.BindType(&self.MacroType, "Macro")
	self.BindType(&self.PosType, "Pos")
	self.BindType(&self.PrimType, "Prim")
	self.BindType(&self.StringType, "String")

	self.BindMacro(&self.BenchMacro, "bench", 2,
		func(self *Macro, args *Forms, vm *Vm, env Env, pos Pos) error {
			reps := args.Pop().(*LitForm).value.d.(int)
			vm.Code[vm.Emit()] = BenchOp(reps)

			if err := args.Pop().Emit(args, vm, env); err != nil {
				return err
			}

			vm.Code[vm.Emit()] = StopOp()			
			return nil
		})
	
	self.BindMacro(&self.FunMacro, "fun", 3,
		func(self *Macro, args *Forms, vm *Vm, env Env, pos Pos) error {
			name := args.Pop().(*IdForm).name
			var funArgs []string

			for _, f := range args.Pop().(*GroupForm).items {
				funArgs = append(funArgs, f.(*IdForm).name)
			}

			gotoPc := vm.Emit()
			prevFun := vm.fun

			defer func () {
				vm.fun = prevFun
			}()

			vm.fun = NewFun(vm, name, vm.EmitPc(), funArgs...)
			vm.env.Bind(name, &vm.AbcLib.FunType, vm.fun)
			
			if err := args.Pop().Emit(args, vm, env); err != nil {
				return err
			}
			
			vm.Code[vm.Emit()] = RetOp()
			vm.Code[gotoPc] = GotoOp(vm.EmitPc())
			return nil
		})

	self.BindMacro(&self.BenchMacro, "if", 2,
		func(self *Macro, args *Forms, vm *Vm, env Env, pos Pos) error {
			if err := args.Pop().Emit(args, vm, env); err != nil {
				return err
			}

			ifPc := vm.Emit()
			
			if err := args.Pop().Emit(args, vm, env); err != nil {
				return err
			}

			elsePc := vm.EmitPc()
			
			if f, ok := args.Top().(*IdForm); f != nil && ok && f.name == "else" {
				args.Pop()
				gotoPc := vm.Emit()
				elsePc = vm.EmitPc()
				
				if err := args.Pop().Emit(args, vm, env); err != nil {
					return err
				}

				vm.Code[gotoPc] = GotoOp(vm.EmitPc())
			}
			
			vm.Code[ifPc] = IfOp(elsePc)
			return nil
		})
	
	self.BindPrim(&self.AddPrim, "+", 2, func(self *Prim, vm *Vm, pos *Pos) error {
		b := vm.Stack.Pop().d.(int)
		a := vm.Stack.Pop().d.(int)
		vm.Stack.Push(V{&vm.AbcLib.IntType, a + b})
		return nil
	})
	
	 self.BindPrim(&self.FailPrim, "fail", 1, func(self *Prim, vm *Vm, pos *Pos) error {
		return vm.E(pos, vm.Stack.Pop().String())
	})

	self.BindPrim(&self.GtPrim, ">", 2, func(self *Prim, vm *Vm, pos *Pos) error {
		b := vm.Stack.Pop().d.(int)
		a := vm.Stack.Pop().d.(int)
		vm.Stack.Push(V{&vm.AbcLib.BoolType, a > b})
		return nil
	})

	self.BindPrim(&self.SubPrim, "-", 2, func(self *Prim, vm *Vm, pos *Pos) error {
		b := vm.Stack.Pop().d.(int)
		a := vm.Stack.Pop().d.(int)
		vm.Stack.Push(V{&vm.AbcLib.IntType, a - b})
		return nil
	})
}

type BoolType struct {
	BasicType
}

func (self *BoolType) Emit(val V, args *Forms, vm *Vm, env Env, pos Pos) error {	
	vm.Code[vm.Emit()] = PushBoolOp(val.d.(bool))
	return nil
}

func (self *BoolType) Bool(val V) bool {
	return val.d.(bool)
}

func (self *BoolType) Dump(val V, out io.Writer) error {
	var err error
	
	if val.d.(bool) {
		_, err = fmt.Fprint(out, "T")
	} else {
		_, err = fmt.Fprint(out, "F")
	}
	
	return err
}

type FunType struct {
	BasicType
}

func (self *FunType) Emit(val V, args *Forms, vm *Vm, env Env, pos Pos) error {	
	vm.Code[vm.Emit()] = PushOp(vm.Tag(self, val.d.(*Fun)))
	return nil
}

func (self *FunType) Dump(val V, out io.Writer) error {
	_, err := io.WriteString(out, val.d.(*Fun).String())
	return err
}

type IntType struct {
	BasicType
}

func (self *IntType) Emit(val V, args *Forms, vm *Vm, env Env, pos Pos) error {	
	vm.Code[vm.Emit()] = PushIntOp(val.d.(int))
	return nil
}

func (self *IntType) Dump(val V, out io.Writer) error {
	_, err := fmt.Fprintf(out, "%v", val.d.(int))
	return err
}

type MacroType struct {
	BasicType
}

func (self *MacroType) Emit(val V, args *Forms, vm *Vm, env Env, pos Pos) error {	
	vm.Code[vm.Emit()] = PushOp(vm.Tag(self, val.d.(*Macro)))
	return nil
}

func (self *MacroType) Dump(val V, out io.Writer) error {
	_, err := io.WriteString(out, val.d.(*Macro).String())
	return err
}

type MetaType struct {
	BasicType
}

func (self *MetaType) Emit(val V, args *Forms, vm *Vm, env Env, pos Pos) error {	
	vm.Code[vm.Emit()] = PushOp(vm.Tag(self, val.d.(Type)))
	return nil
}

func (self *MetaType) Dump(val V, out io.Writer) error {
	_, err := io.WriteString(out, val.d.(Type).Name())
	return err
}

type PosType struct {
	BasicType
}

func (self *PosType) Emit(val V, args *Forms, vm *Vm, env Env, pos Pos) error {	
	vm.EmitPos(val.d.(Pos))
	return nil
}

func (self *PosType) Dump(val V, out io.Writer) error {
	_, err := fmt.Fprintf(out, "%v", val.d.(Pos))
	return err
}

type PrimType struct {
	BasicType
}

func (self *PrimType) Emit(val V, args *Forms, vm *Vm, env Env, pos Pos) error {	
	vm.Code[vm.Emit()] = PushOp(vm.Tag(self, val.d.(*Prim)))
	return nil
}

func (self *PrimType) Dump(val V, out io.Writer) error {
	_, err := io.WriteString(out, val.d.(*Prim).String())
	return err
}

type StringType struct {
	BasicType
}

func (self *StringType) Emit(val V, args *Forms, vm *Vm, env Env, pos Pos) error {	
	vm.EmitString(val.d.(string))
	return nil
}

func (self *StringType) Dump(val V, out io.Writer) error {
	_, err := io.WriteString(out, val.d.(string))
	return err
}

package snabl

import (
	"fmt"
	"io"
	"time"
)

type AbcLib struct {
	BasicLib
	BoolType BoolType
	FormType FormType
	FunType FunType
	IntType IntType
	MacroType MacroType
	MetaType MetaType
	NilType NilType
	PosType PosType
	PrimType PrimType
	SliceType SliceType
	StringType StringType
	TimeType TimeType

	BenchMacro, DebugMacro, DefMacro, DefunMacro, IfMacro, PosMacro, TestMacro, TraceMacro Macro
	
	AddPrim, EqPrim, FailPrim, GtPrim, HoursPrim, LenPrim, LoadPrim, LtPrim, MinsPrim, MsecsPrim, SayPrim,
	SecsPrim, SleepPrim, SubPrim Prim
}

func (self *AbcLib) Init(vm *Vm) {
	self.BasicLib.Init(vm, nil, "abc")
	
	self.BindType(&self.MetaType, "Meta")
	
	self.BindType(&self.BoolType, "Bool")
	self.BindType(&self.FormType, "Form")
	self.BindType(&self.FunType, "Fun")
	self.BindType(&self.IntType, "Int")
	self.BindType(&self.MacroType, "Macro")
	self.BindType(&self.PosType, "Pos")
	self.BindType(&self.PrimType, "Prim")
	self.BindType(&self.SliceType, "Slice")
	self.BindType(&self.StringType, "String")
	self.BindType(&self.TimeType, "Time")

	self.BindMacro(&self.BenchMacro, "bench", 2,
		func(self *Macro, args *Forms, vm *Vm, env Env, pos Pos) error {
			if err := args.Pop().Emit(args, vm, env); err != nil {
				return err
			}
			
			vm.Code[vm.Emit()] = BenchOp()
			
			body := args.Pop()
			
			if err := body.Emit(args, vm, env); err != nil {
				return err
			}

			vm.Code[vm.Emit()] = StopOp()			
			return nil
		})
	
	self.BindMacro(&self.DebugMacro, "debug", 0,
		func(self *Macro, args *Forms, vm *Vm, env Env, pos Pos) error {
			vm.Debug = !vm.Debug
			return nil
		})

	self.BindMacro(&self.DefMacro, "def", 2,
		func(self *Macro, args *Forms, vm *Vm, env Env, pos Pos) error {
			name := args.Pop().(*IdForm).name

			skipPc := vm.Emit()
			pc := vm.EmitPc()
			
			if err := args.Pop().Emit(args, vm, env); err != nil {
				return err
			}

			vm.Code[vm.Emit()] = StopOp()
			vm.Code[skipPc] = GotoOp(vm.EmitPc())

			if err := vm.Eval(&pc); err != nil {
				return err
			}

			v := vm.Stack.Pop()
			env.Bind(name, v.t, v.d)
			return nil
		})

	self.BindMacro(&self.DefunMacro, "defun", 3,
		func(self *Macro, args *Forms, vm *Vm, env Env, pos Pos) error {
			name := args.Pop().(*IdForm).name
			var funArgs []string

			for _, f := range args.Pop().(*GroupForm).items {
				funArgs = append(funArgs, f.(*IdForm).name)
			}

			skipPc := vm.Emit()
			prevFun := vm.fun

			defer func () {
				vm.fun = prevFun
			}()

			vm.fun = NewFun(vm, name, vm.EmitPc(), funArgs...)
			env.Bind(name, &vm.AbcLib.FunType, vm.fun)
			
			if err := args.Pop().Emit(args, vm, env); err != nil {
				return err
			}
			
			vm.Code[vm.Emit()] = RetOp()
			vm.Code[skipPc] = GotoOp(vm.EmitPc())
			return nil
		})

	self.BindMacro(&self.IfMacro, "if", 2,
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

	self.BindMacro(&self.PosMacro, "pos", 0,
		func(self *Macro, args *Forms, vm *Vm, env Env, pos Pos) error {
			vm.EmitTag(&vm.AbcLib.PosType, pos)
			return nil
		})

	self.BindMacro(&self.TestMacro, "test", 2,
		func(self *Macro, args *Forms, vm *Vm, env Env, pos Pos) error {
			if err := args.Pop().Emit(args, vm, env); err != nil {
				return err
			}

			pc := vm.Emit()
			forms := *args
			
			if err := args.Pop().Emit(args, vm, env); err != nil {
				return err
			}

			f := NewGroupForm(forms.Top().Pos(), forms.items[:forms.Len()-args.Len()]...)
			vm.Code[pc] = TestOp(vm.Tag(&vm.AbcLib.FormType, f))
			vm.Code[vm.Emit()] = StopOp()			
			return nil
		})

	self.BindMacro(&self.TraceMacro, "trace", 0,
		func(self *Macro, args *Forms, vm *Vm, env Env, pos Pos) error {
			vm.Trace = !vm.Trace
			return nil
		})

	self.BindPrim(&self.AddPrim, "+", 2, func(self *Prim, vm *Vm, pos *Pos) error {
		b := vm.Stack.Pop().d.(int)
		a := vm.Stack.Top(0)
		a.Init(&vm.AbcLib.IntType, a.d.(int) + b)
		return nil
	})

	self.BindPrim(&self.EqPrim, "=", 2, func(self *Prim, vm *Vm, pos *Pos) error {
		b := vm.Stack.Pop()
		a := vm.Stack.Top(0)
		a.Init(&vm.AbcLib.BoolType, a.Eq(*b))
		return nil
	})

	self.BindPrim(&self.FailPrim, "fail", 1, func(self *Prim, vm *Vm, pos *Pos) error {
		return vm.E(pos, vm.Stack.Pop().String())
	})

	self.BindPrim(&self.GtPrim, ">", 2, func(self *Prim, vm *Vm, pos *Pos) error {
		b := vm.Stack.Pop()
		a := vm.Stack.Top(0)

		if a.t != b.t {
			return vm.E(pos, "Type mismtch: %v/%v", a.t.String(), b.t.String())
		}

		t, ok := a.t.(CmpType)

		if !ok {
			return vm.E(pos, "> not supported: %v", t.String())
		}
		
		a.Init(&vm.AbcLib.BoolType, t.Gt(*a, *b))
		return nil
	})

	self.BindPrim(&self.HoursPrim, "hours", 1, func(self *Prim, vm *Vm, pos *Pos) error {
		v := vm.Stack.Top(0)

		switch v.t {
		case &vm.AbcLib.IntType:
			v.Init(&vm.AbcLib.TimeType, time.Duration(v.d.(int)) * time.Hour)
		case &vm.AbcLib.TimeType:
			v.Init(&vm.AbcLib.IntType, int(v.d.(time.Duration).Hours()))
		default:
			return vm.E(pos, "hours not supported: %v", v.String())
		}

		return nil
	})

	self.BindPrim(&self.LenPrim, "len", 1, func(self *Prim, vm *Vm, pos *Pos) error {
		v := vm.Stack.Pop()
		t, ok := v.t.(LenType)

		if !ok {
			return vm.E(pos, "'len' not supported: %v", v.String())
		}

		vm.Stack.Push(&vm.AbcLib.IntType, t.Len(*v))
		return nil
	})

	self.BindPrim(&self.LoadPrim, "load", 1, func(self *Prim, vm *Vm, pos *Pos) error {
		p := vm.Stack.Pop().d.(string)
		return vm.Load(p, true)
	})

	self.BindPrim(&self.LtPrim, "<", 2, func(self *Prim, vm *Vm, pos *Pos) error {
		b := vm.Stack.Pop()
		a := vm.Stack.Top(0)

		if a.t != b.t {
			return vm.E(pos, "Type mismtch: %v/%v", a.t.String(), b.t.String())
		}

		t, ok := a.t.(CmpType)

		if !ok {
			return vm.E(pos, "> not supported: %v", t.String())
		}
		
		a.Init(&vm.AbcLib.BoolType, t.Lt(*a, *b))
		return nil
	})

	self.BindPrim(&self.MinsPrim, "mins", 1, func(self *Prim, vm *Vm, pos *Pos) error {
		v := vm.Stack.Top(0)

		switch v.t {
		case &vm.AbcLib.IntType:
			v.Init(&vm.AbcLib.TimeType, time.Duration(v.d.(int)) * time.Minute)
		case &vm.AbcLib.TimeType:
			v.Init(&vm.AbcLib.IntType, int(v.d.(time.Duration).Minutes()))
		default:
			return vm.E(pos, "mins not supported: %v", v.String())
		}

		return nil
	})

	self.BindPrim(&self.MsecsPrim, "msecs", 1, func(self *Prim, vm *Vm, pos *Pos) error {
		v := vm.Stack.Top(0)

		switch v.t {
		case &vm.AbcLib.IntType:
			v.Init(&vm.AbcLib.TimeType, time.Duration(v.d.(int)) * time.Millisecond)
		case &vm.AbcLib.TimeType:
			v.Init(&vm.AbcLib.IntType, int(v.d.(time.Duration).Milliseconds()))
		default:
			return vm.E(pos, "msecs not supported: %v", v.String())
		}

		return nil
	})

	self.BindPrim(&self.SayPrim, "say", 1, func(self *Prim, vm *Vm, pos *Pos) error {
		if err := vm.Stack.Pop().Write(vm.Stdout); err != nil {
			return err
		}

		if _, err := fmt.Fprintln(vm.Stdout, ""); err != nil {
			return err
		}

		return nil
	})

	self.BindPrim(&self.SecsPrim, "secs", 1, func(self *Prim, vm *Vm, pos *Pos) error {
		v := vm.Stack.Top(0)

		switch v.t {
		case &vm.AbcLib.IntType:
			v.Init(&vm.AbcLib.TimeType, time.Duration(v.d.(int)) * time.Second)
		case &vm.AbcLib.TimeType:
			v.Init(&vm.AbcLib.IntType, int(v.d.(time.Duration).Seconds()))
		default:
			return vm.E(pos, "secs not supported: %v", v.String())
		}

		return nil
	})

	self.BindPrim(&self.SleepPrim, "sleep", 1, func(self *Prim, vm *Vm, pos *Pos) error {
		time.Sleep(vm.Stack.Pop().d.(time.Duration))
		return nil
	})
	
	self.BindPrim(&self.SubPrim, "-", 2, func(self *Prim, vm *Vm, pos *Pos) error {
		b := vm.Stack.Pop().d.(int)
		a := vm.Stack.Top(0)
		a.Init(&vm.AbcLib.IntType, a.d.(int) - b)
		return nil
	})

	self.Bind("T", &self.BoolType, true)
	self.Bind("F", &self.BoolType, false)
	self.Bind("NIL", &self.NilType, nil)
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

func (self *BoolType) Write(val V, out io.Writer) error {
	return self.Dump(val, out)
}

type FormType struct {
	BasicType
}

func (self *FormType) Dump(val V, out io.Writer) error {
	_, err := io.WriteString(out, val.d.(Form).String())
	return err
}

func (self *FormType) Write(val V, out io.Writer) error {
	return self.Dump(val, out)
}

type FunType struct {
	BasicType
}

func (self *FunType) Dump(val V, out io.Writer) error {
	_, err := io.WriteString(out, val.d.(*Fun).String())
	return err
}

func (self *FunType) Write(val V, out io.Writer) error {
	return self.Dump(val, out)
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

func (self *IntType) Write(val V, out io.Writer) error {
	return self.Dump(val, out)
}

func (self *IntType) Gt(left, right V) bool {
	return left.d.(int) > right.d.(int)
}

func (self *IntType) Lt(left, right V) bool {
	return left.d.(int) < right.d.(int)
}

type MacroType struct {
	BasicType
}

func (self *MacroType) Dump(val V, out io.Writer) error {
	_, err := io.WriteString(out, val.d.(*Macro).String())
	return err
}

func (self *MacroType) Write(val V, out io.Writer) error {
	return self.Dump(val, out)
}

type MetaType struct {
	BasicType
}

func (self *MetaType) Dump(val V, out io.Writer) error {
	_, err := io.WriteString(out, val.d.(Type).String())
	return err
}

func (self *MetaType) Write(val V, out io.Writer) error {
	return self.Dump(val, out)
}

type NilType struct {
	BasicType
}

func (self *NilType) Emit(val V, args *Forms, vm *Vm, env Env, pos Pos) error {	
	vm.Code[vm.Emit()] = PushNilOp()
	return nil
}

func (self *NilType) Bool(val V) bool {
	return false
}

func (self *NilType) Dump(val V, out io.Writer) error {
	_, err := io.WriteString(out, "NIL")
	return err
}

func (self *NilType) Write(val V, out io.Writer) error {
	return self.Dump(val, out)
}

type PosType struct {
	BasicType
}

func (self *PosType) Dump(val V, out io.Writer) error {
	_, err := fmt.Fprintf(out, "%v", val.d.(Pos))
	return err
}

func (self *PosType) Write(val V, out io.Writer) error {
	return self.Dump(val, out)
}

type PrimType struct {
	BasicType
}

func (self *PrimType) Dump(val V, out io.Writer) error {
	_, err := io.WriteString(out, val.d.(*Prim).String())
	return err
}

func (self *PrimType) Write(val V, out io.Writer) error {
	return self.Dump(val, out)
}

type SliceType struct {
	BasicType
}

func (self *SliceType) Dump(val V, out io.Writer) error {
	return val.d.(*Slice).Dump(out)
}

func (self *SliceType) Write(val V, out io.Writer) error {
	return self.Dump(val, out)
}

func (self *SliceType) Len(val V) int {
	return val.d.(*Slice).Len()
}

type StringType struct {
	BasicType
}

func (self *StringType) Dump(val V, out io.Writer) error {
	_, err := fmt.Fprintf(out, "\"%v\"", val.d.(string))
	return err
}

func (self *StringType) Write(val V, out io.Writer) error {
	_, err := io.WriteString(out, val.d.(string))
	return err
}

func (self *StringType) Len(val V) int {
	return len(val.d.(string))
}

func (self *StringType) Gt(left, right V) bool {
	return left.d.(string) > right.d.(string)
}

func (self *StringType) Lt(left, right V) bool {
	return left.d.(string) < right.d.(string)
}

type TimeType struct {
	BasicType
}

func (self *TimeType) Emit(val V, args *Forms, vm *Vm, env Env, pos Pos) error {	
	vm.Code[vm.Emit()] = PushTimeOp(val.d.(time.Duration))
	return nil
}

func (self *TimeType) Dump(val V, out io.Writer) error {
	_, err := fmt.Fprintf(out, "%v", val.d.(time.Duration))
	return err
}

func (self *TimeType) Write(val V, out io.Writer) error {
	return self.Dump(val, out)
}

func (self *TimeType) Gt(left, right V) bool {
	return left.d.(time.Duration) > right.d.(time.Duration)
}

func (self *TimeType) Lt(left, right V) bool {
	return left.d.(time.Duration) < right.d.(time.Duration)
}

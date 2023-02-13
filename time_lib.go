package snabl

import (
	"fmt"
	"io"
	"time"
)

type TimeLib struct {
	BasicLib
	TimeType TimeType

	HoursPrim, MinsPrim, MsecsPrim, SecsPrim, SleepPrim Prim
}

func (self *TimeLib) Init(vm *Vm) {
	self.BasicLib.Init(vm, nil, "time")
	self.BindType(&self.TimeType, "Time")

	self.BindPrim(&self.HoursPrim, "hours", 1, func(self *Prim, vm *Vm, pos *Pos) error {
		v := vm.Stack.Top(0)

		switch v.t {
		case &vm.AbcLib.IntType:
			v.Init(&vm.TimeLib.TimeType, time.Duration(v.d.(int)) * time.Hour)
		case &vm.TimeLib.TimeType:
			v.Init(&vm.AbcLib.IntType, int(v.d.(time.Duration).Hours()))
		default:
			return vm.E(pos, "hours not supported: %v", v.String())
		}

		return nil
	})

	self.BindPrim(&self.MinsPrim, "mins", 1, func(self *Prim, vm *Vm, pos *Pos) error {
		v := vm.Stack.Top(0)

		switch v.t {
		case &vm.AbcLib.IntType:
			v.Init(&vm.TimeLib.TimeType, time.Duration(v.d.(int)) * time.Minute)
		case &vm.TimeLib.TimeType:
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
			v.Init(&vm.TimeLib.TimeType, time.Duration(v.d.(int)) * time.Millisecond)
		case &vm.TimeLib.TimeType:
			v.Init(&vm.AbcLib.IntType, int(v.d.(time.Duration).Milliseconds()))
		default:
			return vm.E(pos, "msecs not supported: %v", v.String())
		}

		return nil
	})

	self.BindPrim(&self.SecsPrim, "secs", 1, func(self *Prim, vm *Vm, pos *Pos) error {
		v := vm.Stack.Top(0)

		switch v.t {
		case &vm.AbcLib.IntType:
			v.Init(&vm.TimeLib.TimeType, time.Duration(v.d.(int)) * time.Second)
		case &vm.TimeLib.TimeType:
			v.Init(&vm.AbcLib.IntType, int(v.d.(time.Duration).Seconds()))
		default:
			return vm.E(pos, "secs not supported: %v", v.String())
		}

		return nil
	})
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

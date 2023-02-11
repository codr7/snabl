package snabl

import (
	//"fmt"
	"io"
	"strings"
)

type Form interface {
	Emit(args *Forms, vm *Vm, env Env) error
	String() string
}

type BasicForm struct {
	pos Pos
}

func (self *BasicForm) Init(pos Pos) {
	self.pos = pos
}

func (self *BasicForm) Emit(args *Forms, vm *Vm, env Env) error {
	if vm.Debug {
		vm.EmitPos(self.pos)
	}

	return nil
}

type GroupForm struct {
	BasicForm
	items []Form
}

func NewGroupForm(pos Pos, items...Form) *GroupForm {
	return new(GroupForm).Init(pos, items...)
}

func (self *GroupForm) Init(pos Pos, items...Form) *GroupForm {
	self.BasicForm.Init(pos)
	self.items = items
	return self
}

func (self *GroupForm) Emit(args *Forms, vm *Vm, env Env) error {
	if err := self.BasicForm.Emit(args, vm, env); err != nil {
		return err
	}

	for _, f := range self.items {
		if err := f.Emit(args, vm, env); err != nil {
			return err
		}
	}

	return nil
}

func (self *GroupForm) String() string {
	var out strings.Builder
	io.WriteString(&out, "(")

	for i, f := range self.items {
		if i > 0 {
			io.WriteString(&out, " ")
		}
		
		io.WriteString(&out, f.String())
	}

	io.WriteString(&out, ")")
	return out.String()
}

type IdForm struct {
	BasicForm
	name string
}

func NewIdForm(pos Pos, name string) *IdForm {
	return new(IdForm).Init(pos, name)
}

func (self *IdForm) Init(pos Pos, name string) *IdForm {
	self.BasicForm.Init(pos)
	self.name = name
	return self
}

func (self *IdForm) Emit(args *Forms, vm *Vm, env Env) error {
	if err := self.BasicForm.Emit(args, vm, env); err != nil {
		return err
	}

	if vm.fun != nil {
		i := vm.fun.ArgIndex(self.name)

		if i > -1 {
			vm.Code[vm.Emit()] = ArgOp(i)
			return nil
		}
	}
	
	found := env.Find(self.name)

	if found == nil {
		return vm.E(&self.pos, "%v?", self.name)
	}
	
	if found.t == &vm.AbcLib.FunType {
		fun := found.d.(*Fun)

		for i := 0; i < fun.Arity(); i++ {
			if err := args.Pop().Emit(args, vm, env); err != nil {
				return err
			}
		}
		
		vm.Code[vm.Emit()] = CallFunOp(fun)
	} else if found.t == &vm.AbcLib.MacroType {
		return found.d.(*Macro).Emit(args, vm, env, self.pos)
	} else if found.t == &vm.AbcLib.PrimType {
		prim := found.d.(*Prim)

		for i := 0; i < prim.arity; i++ {
			if err := args.Pop().Emit(args, vm, env); err != nil {
				return err
			}
		}

		vm.Code[vm.Emit()] = CallPrimOp(prim)
	} else {
		return found.Emit(args, vm, env, self.pos)
	}

	return nil
}

func (self *IdForm) String() string {
	return self.name
}

type LitForm struct {
	BasicForm
	value V
}

func NewLitForm(pos Pos, t Type, d any) *LitForm {
	return new(LitForm).Init(pos, t, d)
}

func (self *LitForm) Init(pos Pos, t Type, d any) *LitForm {
	self.BasicForm.Init(pos)
	self.value.Init(t, d)
	return self
}

func (self *LitForm) Emit(args *Forms, vm *Vm, env Env) error {
	if err := self.BasicForm.Emit(args, vm, env); err != nil {
		return err
	}
	
	return self.value.Emit(args, vm, env, self.pos)
}

func (self *LitForm) String() string {
	return self.value.String()
}

type Forms struct {
	items []Form
}

func (self *Forms) Top() Form {
	i := len(self.items)
	
	if i == 0 {
		return nil
	}

	return self.items[0]
}

func (self *Forms) Pop() Form {
	i := len(self.items)
	
	if i == 0 {
		return nil
	}

	f := self.items[0]
	self.items = self.items[1:]
	return f
}

func (self *Forms) Init(items []Form) *Forms {
	self.items = items
	return self
}

func (self *Forms) Len() int {
	return len(self.items)
}

func (self *Forms) Push(form Form) {
	self.items = append(self.items, form)
}

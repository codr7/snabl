package snabl

import (
	"io"
)

type Form interface {
	Emit(args *Forms, vm *Vm, env Env) error
	Dump(out io.Writer) error
}

type BasicForm struct {
	pos Pos
}

func (self *BasicForm) Init(pos Pos) {
	self.pos = pos
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
	for _, f := range self.items {
		if err := f.Emit(args, vm, env); err != nil {
			return err
		}
	}

	return nil
}

func (self *GroupForm) Dump(out io.Writer) error {
	if _, err := io.WriteString(out, "("); err != nil {
		return err
	}

	for i, f := range self.items {
		if i > 0 {
			if _, err := io.WriteString(out, " "); err != nil {
				return err
			}
		}
		
		if err := f.Dump(out); err != nil {
			return err
		}
	}

	if _, err := io.WriteString(out, ")"); err != nil {
		return err
	}
	
	return nil
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
	if vm.fun != nil {
		i := vm.fun.ArgIndex(self.name)

		if i > -1 {
			vm.Code[vm.Emit()] = ArgOp(vm.fun.argOffsTag, i)
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
		
		tag := vm.Tag(&vm.AbcLib.FunType, fun)
		vm.Code[vm.Emit()] = CallFunOp(tag)
	} else if found.t == &vm.AbcLib.MacroType {
		return found.d.(*Macro).Emit(args, vm, env, &self.pos)
	} else {
		return found.Emit(args, vm, env, self.pos)
	}

	return nil
}

func (self *IdForm) Dump(out io.Writer) error {
	_, err := io.WriteString(out, self.name)
	return err
}

type LitForm struct {
	BasicForm
	val V
}

func NewLitForm(pos Pos, t Type, d any) *LitForm {
	return new(LitForm).Init(pos, t, d)
}

func (self *LitForm) Init(pos Pos, t Type, d any) *LitForm {
	self.BasicForm.Init(pos)
	self.val.Init(t, d)
	return self
}

func (self *LitForm) Emit(args *Forms, vm *Vm, env Env) error {
	return self.val.Emit(args, vm, env, self.pos)
}

func (self *LitForm) Dump(out io.Writer) error {
	return self.val.Dump(out)
}

type Forms struct {
	items []Form
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

func (self *Forms) Push(form Form) {
	self.items = append(self.items, form)
}

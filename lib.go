package snabl

type Lib interface {
	Name() string
	BindPrim(name string, body PrimBody)
}

type BasicLib struct {
	BasicEnv
	vm *Vm
	name string
}

func (self *BasicLib) Init(vm *Vm, parent Env, name string) {
	self.BasicEnv.Init(parent)
	self.vm = vm
	self.name = name
}

func (self *BasicLib) Name() string {
	return self.name
}

func (self *BasicLib) BindMacro(m *Macro, name string, arity int, body MacroBody) {
	m.Init(name, arity, body)
	self.Bind(name, &self.vm.AbcLib.MacroType, m)
}

func (self *BasicLib) BindPrim(p *Prim, name string, arity int, body PrimBody) {
	p.Init(self.vm, name, arity, body)
	self.Bind(name, &self.vm.AbcLib.PrimType, p)
}

func (self *BasicLib) BindType(t Type, name string) {
	t.Init(name)
	self.Bind(name, &self.vm.AbcLib.MetaType, t)
}

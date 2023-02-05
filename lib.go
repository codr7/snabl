package snabl

type Lib interface {
	Name() string
	BindPrim(name string, body PrimBody)
}

type BasicLib struct {
	Env
	name string
}

func (self *BasicLib) Init(name string) {
	self.Env.Init()
	self.name = name
}

func (self *BasicLib) Name() string {
	return self.name
}

func (self *BasicLib) BindPrim(p *Prim, name string, arity uint, body PrimBody) {
	p.Init(name, arity, body)
	self.Bind(name, &Abc.PrimType, p)
}

func (self *BasicLib) BindType(t Type, name string) {
	t.Init(name)
	self.Bind(name, &Abc.PrimType, t)
}

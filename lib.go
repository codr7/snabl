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

func (self *BasicLib) BindPrim(name string, arity uint, body PrimBody) *Prim {
	p := NewPrim(name, arity, body)
	self.Bind(name, &Abc.PrimType, p)
	return p
}

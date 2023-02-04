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

func (self *BasicLib) BindPrim(name string, arity uint, body PrimBody) {
	self.Bind(name, &Abc.PrimType, NewPrim(name, arity, body))
}

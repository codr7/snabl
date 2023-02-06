package snabl

type Env interface {
	Bind(id string, t Type, d any)
	Find(id string) *V
}

type BasicEnv struct {
	bindings map[string]V
}

func NewEnv() *BasicEnv {
	return new(BasicEnv).Init()
}

func (self *BasicEnv) Init() *BasicEnv {
	self.bindings = make(map[string]V)
	return self
}

func (self *BasicEnv) Bind(id string, t Type, d any) {
	self.bindings[id] = V{t: t, d: d}
}

func (self *BasicEnv) Find(id string) *V {
	v, ok := self.bindings[id]

	if !ok {
		return nil
	}

	return &v
}

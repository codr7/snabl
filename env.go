package snabl

type Env interface {
	Init()
	Bind(id string, t VT, d any)
	Find(id string) *V
}

type BasicEnv struct {
	bindings map[string]V
}

func (self *BasicEnv) Init() {
	self.bindings = make(map[string]V)
}

func (self *BasicEnv) Bind(id string, t VT, d any) {
	self.bindings[id] = V{t: t, d: d}
}

func (self *BasicEnv) Find(id string) *V {
	v, ok := self.bindings[id]

	if !ok {
		return nil
	}

	return &v
}

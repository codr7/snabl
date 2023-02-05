package snabl

type Env struct {
	bindings map[string]V
}

func (self *Env) Init() {
	self.bindings = make(map[string]V)
}

func (self *Env) Bind(id string, t Type, d any) {
	self.bindings[id] = V{t: t, d: d}
}

func (self *Env) Find(id string) *V {
	v, ok := self.bindings[id]

	if !ok {
		return nil
	}

	return &v
}

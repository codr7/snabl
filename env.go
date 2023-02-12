package snabl

type Env interface {
	Bind(id string, t Type, d any)
	Find(id string) *V
	Each(cb func (k string, v V))
	Import(source Env, vm *Vm, pos *Pos, names...string) error
}

type BasicEnv struct {
	parent Env
	bindings map[string]V
}

func NewEnv(parent Env) *BasicEnv {
	return new(BasicEnv).Init(parent)
}

func (self *BasicEnv) Init(parent Env) *BasicEnv {
	self.parent = parent
	self.bindings = make(map[string]V)
	return self
}

func (self *BasicEnv) Bind(id string, t Type, d any) {
	self.bindings[id] = V{t: t, d: d}
}

func (self *BasicEnv) Find(id string) *V {
	v, ok := self.bindings[id]

	if !ok {
		if self.parent != nil {
			return self.parent.Find(id)
		}
		
		return nil
	}

	return &v
}

func (self *BasicEnv) Each(cb func (k string, v V)) {
	for k, v := range self.bindings {
		cb(k, v)
	}
}

func (self *BasicEnv) Import(source Env, vm *Vm, pos *Pos, names...string) error {
	if len(names) == 0 {
		source.Each(func (k string, v V) {
			self.Bind(k, v.t, v.d)
		})
	} else {
		for _, k := range names {
			v := source.Find(k)

			if v == nil {
				return vm.E(pos, "%v?", k) 
			}

			self.Bind(k, v.t, v.d)
		}
	}

	return nil
}

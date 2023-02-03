package snabl

type Lib interface {
	Name() string
}

type BasicLib struct {
	name string
}

func (self *BasicLib) Init(name string) {
	self.name = name
}

func (self *BasicLib) Name() string {
	return self.name
}

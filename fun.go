package snabl

import (
	"fmt"
)

type ArgOffs = uint
type ArgIndex = uint

type Fun struct {
	name string
	pc Pc
	args []string
}

func (self *Fun) Init(name string, pc Pc, args...string) {
	self.name = name
	self.pc = pc
	self.args = args
}

func (self *Fun) String() string {
	return self.name
}

func (self *Fun) ArgIndex(arg string) ArgIndex {
	for i, a := range self.args {
		if a == arg {
			return ArgIndex(i);
		}
	}

	panic(fmt.Sprintf("Arg not found in %v: %v", self.name, arg))
}

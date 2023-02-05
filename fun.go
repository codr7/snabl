package snabl

import (
	"fmt"
)

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

func (self *Fun) ArgIndex(arg string) int {
	for i, a := range self.args {
		if a == arg {
			return len(self.args) - i - 1;
		}
	}

	panic(fmt.Sprintf("Arg not found in %v: %v", self.name, arg))
}

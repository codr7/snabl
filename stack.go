package snabl

import (
	"io"
	"strings"
)

type Stack struct {
	items []V
}

func NewStack(items []V) *Stack {
	return new(Stack).Init(items)
}

func (self *Stack) Init(items []V) *Stack {
	self.items = items
	return self
}

func (self *Stack) Clear() {
	self.items = nil
}

func (self Stack) Dump(out io.Writer) error {
	if _, err := io.WriteString(out, "["); err != nil {
		return err
	}

	for i, v := range self.items {
		if i > 0 {
			io.WriteString(out, " ");
		}
		
		v.Dump(out)
	}
	
	if _, err := io.WriteString(out, "]"); err != nil {
		return err
	}
	
	return nil
}

func (self Stack) Len() int {
	return len(self.items)
}

func (self Stack) Top() *V {
	i := len(self.items)
	
	if i == 0 {
		return nil
	}

	return &self.items[i-1]
}

func (self *Stack) Pop() *V {
	i := len(self.items)
	
	if i == 0 {
		return nil
	}

	v := self.items[i-1]
	self.items = self.items[:i-1]
	return &v
}

func (self *Stack) Push(t VT, d any) {
	self.items = append(self.items, V{t: t, d: d})
}

func (self *Stack) String() string {
	var out strings.Builder
	self.Dump(&out)
	return out.String()
}

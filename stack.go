package snabl

import (
	"io"
	"strings"
)

type Stack struct {
	parent *Stack
	items []V
}

func (self *Stack) Init(parent *Stack) *Stack {
	self.parent = parent
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
		
		io.WriteString(out, v.String())
	}
	
	if _, err := io.WriteString(out, "]"); err != nil {
		return err
	}
	
	return nil
}

func (self Stack) Len() int {
	return len(self.items)
}

func (self Stack) Top(i int) *V {
	n := len(self.items) - 1
	
	if n < i {
		return nil
	}

	return &self.items[n-i]
}

func (self *Stack) Pop() *V {
	i := len(self.items)-1
	
	if i < 0 {
		return nil
	}

	v := self.items[i]
	self.items = self.items[:i]
	return &v
}

func (self *Stack) Tail(n int) []V {
	l := len(self.items)
	
	if l < n {
		return nil
	}

	out := make([]V, n)
	copy(out, self.items[l-n:])
	self.items = self.items[:l-n]
	return out
}

func (self *Stack) Push(t Type, d any) {
	self.items = append(self.items, V{t: t, d: d})
}

func (self Stack) String() string {
	var out strings.Builder
	self.Dump(&out)
	return out.String()
}

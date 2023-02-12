package snabl

type Stack struct {
	Slice
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

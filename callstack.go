package snabl

import (
	"fmt"
	"io"
	"strings"
)

type CallStack [T fmt.Stringer] struct {
	items []T
}

func NewCallStack[T fmt.Stringer](items []T) *CallStack[T] {
	return new(CallStack[T]).Init(items)
}

func (self *CallStack[T]) Init(items []T) *CallStack[T] {
	self.items = items
	return self
}

func (self *CallStack[T]) Clear() {
	self.items = nil
}

func (self CallStack[T]) Dump(out io.Writer) error {
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

func (self CallStack[T]) Len() int {
	return len(self.items)
}

func (self CallStack[T]) Top(i int) *T {
	n := len(self.items) - 1
	
	if n < i {
		return nil
	}

	return &self.items[n-i]
}

func (self *CallStack[T]) Pop() *T {
	i := len(self.items)
	
	if i == 0 {
		return nil
	}

	v := self.items[i-1]
	self.items = self.items[:i-1]
	return &v
}

func (self *CallStack[T]) Push(val T) {
	self.items = append(self.items, val)
}

func (self *CallStack[T]) String() string {
	var out strings.Builder
	self.Dump(&out)
	return out.String()
}

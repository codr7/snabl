package snabl

import (
	"fmt"
	"io"
	"strings"
)

type Slice [T fmt.Stringer] struct {
	items []T
}

func NewSlice[T fmt.Stringer](items []T) *Slice[T] {
	return new(Slice[T]).Init(items)
}

func (self *Slice[T]) Init(items []T) *Slice[T] {
	self.items = items
	return self
}

func (self *Slice[T]) Clear() {
	self.items = nil
}

func (self Slice[T]) Dump(out io.Writer) error {
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

func (self Slice[T]) Len() int {
	return len(self.items)
}

func (self Slice[T]) Top(i int) *T {
	n := len(self.items) - 1
	
	if n < i {
		return nil
	}

	return &self.items[n-i]
}

func (self *Slice[T]) Pop() *T {
	i := len(self.items)-1
	
	if i < 0 {
		return nil
	}

	v := self.items[i]
	self.items = self.items[:i]
	return &v
}

func (self *Slice[T]) Drop(n int) []T {
	l := len(self.items)
	
	if l < n {
		return nil
	}

	out := make([]T, n)
	copy(out, self.items[l-n:])
	self.items = self.items[:l-n]
	return out
}

func (self *Slice[T]) Push(val T) {
	self.items = append(self.items, val)
}

func (self *Slice[T]) String() string {
	var out strings.Builder
	self.Dump(&out)
	return out.String()
}

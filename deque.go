package snabl

import (
	"fmt"
	"io"
	"strings"
)

type Deque [T fmt.Stringer] struct {
	items []T
}

func NewDeque[T fmt.Stringer](items []T) *Deque[T] {
	return new(Deque[T]).Init(items)
}

func (self *Deque[T]) Init(items []T) *Deque[T] {
	self.items = items
	return self
}

func (self *Deque[T]) Clear() {
	self.items = nil
}

func (self Deque[T]) Dump(out io.Writer) error {
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

func (self Deque[T]) Len() int {
	return len(self.items)
}

func (self Deque[T]) Top(i int) *T {
	n := len(self.items) - 1
	
	if n < i {
		return nil
	}

	return &self.items[n-i]
}

func (self *Deque[T]) Pop() *T {
	i := len(self.items)-1
	
	if i < 0 {
		return nil
	}

	v := self.items[i]
	self.items = self.items[:i]
	return &v
}

func (self *Deque[T]) Drop(n int) []T {
	l := len(self.items)
	
	if l < n {
		return nil
	}

	out := make([]T, n)
	copy(out, self.items[l-n:])
	self.items = self.items[:l-n]
	return out
}

func (self *Deque[T]) Push(val T) {
	self.items = append(self.items, val)
}

func (self *Deque[T]) String() string {
	var out strings.Builder
	self.Dump(&out)
	return out.String()
}

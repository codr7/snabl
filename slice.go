package snabl

import (
	"fmt"
	"io"
	"strings"
	"unsafe"
)

type Slice struct {
	items []V
}

func (self *Slice) Clear() {
	self.items = nil
}

func (self Slice) Len() int {
	return len(self.items)
}

func (self Slice) Dump(out io.Writer) error {
	if _, err := io.WriteString(out, "["); err != nil {
		return err
	}

	for i, v := range self.items {
		if i > 0 {
			if _, err := io.WriteString(out, " "); err != nil {
				return err
			}
		}

		if err := v.Dump(out); err != nil {
			return err
		}
	}
	
	if _, err := io.WriteString(out, "]"); err != nil {
		return err
	}
	
	return nil
}

func (self Slice) String() string {
	var out strings.Builder
	self.Dump(&out)
	return out.String()
}


type SliceType struct {
	BasicType
}

func (self *SliceType) Dump(val V, out io.Writer) error {
	return val.d.(*Slice).Dump(out)
}

func (self *SliceType) Write(val V, out io.Writer) error {
	for _, v := range val.d.(*Slice).items {
		if err := v.Write(out); err != nil {
			return err
		}
	}
	
	return nil
}

func (self *SliceType) Len(val V) int {
	return val.d.(*Slice).Len()
}

func (self *SliceType) Gt(left, right V, vm *Vm, pos *Pos) (bool, error) {
	ls := left.d.(*Slice)
	rs := right.d.(*Slice)

	for i, lv := range ls.items {
		if rs.Len() <= i {
			return true, nil
		}

		rv := rs.items[i]

		if lv.t != rv.t {
			return false, vm.E(pos, "Type mismatch: %v/%v", lv.t.String(), rv.t.String())
		}

		t, ok := lv.t.(CmpType)

		if !ok {
			return false, vm.E(pos, "> not supported: %v", lv.String())
		}

		if gt, err := t.Gt(lv, rv, vm, pos); err != nil {
			return false, err
		} else if gt {
			return true, nil
		}	
	}

	return rs.Len() == ls.Len(), nil
}

func (self *SliceType) Lt(left, right V, vm *Vm, pos *Pos) (bool, error) {
	ls := left.d.(*Slice)
	rs := right.d.(*Slice)

	for i, lv := range ls.items {
		if rs.Len() <= i {
			return false, nil
		}

		rv := rs.items[i]

		if lv.t != rv.t {
			return false, vm.E(pos, "Type mismatch: %v/%v", lv.t.String(), rv.t.String())
		}

		t, ok := lv.t.(CmpType)

		if !ok {
			return false, vm.E(pos, "< not supported: %v", lv.String())
		}

		if lt, err := t.Lt(lv, rv, vm, pos); err != nil {
			return false, err
		} else if lt {
			return true, nil
		}
	}

	return rs.Len() > ls.Len(), nil
}

func (self *SliceType) Iter(val V) Iter {
	return NewSliceIter(val.d.(*Slice))
}

type SliceIter struct {
	slice *Slice
	i int
}

func NewSliceIter(slice *Slice) *SliceIter {
	return &SliceIter{slice: slice}
}

func (self *SliceIter) Next() (*V, error) {
	if self.i < self.slice.Len() {
		return &self.slice.items[self.i], nil
		self.i++
	}

	return nil, nil
}

func (self *SliceIter) Dump(out io.Writer) error {
	_, err := fmt.Fprintf(out, "SliceIter(%v)", unsafe.Pointer(self))
	return err
}


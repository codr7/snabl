package snabl

import (
	"io"
	"strings"
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

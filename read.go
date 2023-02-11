package snabl

import (
	"bufio"
	//"fmt"
	"io"
	"strings"
	"unicode"
)

func (self *Vm) ReadForms(pos *Pos, in *bufio.Reader, out *Forms) error {
	for {
		if err := self.ReadForm(pos, in, out); err != nil {
			if err == io.EOF {
				return nil
			}

			return err
		}
	}
}

func (self *Vm) ReadForm(pos *Pos, in *bufio.Reader, out *Forms) error {
NEXT:
	c, _, err := in.ReadRune()
	
	if err != nil {
		return err
	}

	switch c {
	case '\n':
		pos.line++
		pos.column = 0
		goto NEXT
	case ' ':
		pos.column++
		goto NEXT
	case '(':
		pos.column++
		return self.ReadGroup(pos, in, out)
	case '"':
		pos.column++
		return self.ReadString(pos, in, out)
	default:
		if unicode.IsDigit(c) {
			in.UnreadRune()
			return self.ReadInt(pos, in, out)
		} else if !unicode.IsSpace(c) && !unicode.IsControl(c) {
			in.UnreadRune()
			return self.ReadId(pos, in, out)			
		}
	}

	return self.E(pos, "%v?", c)
}

func (self *Vm) ReadGroup(pos *Pos, in *bufio.Reader, out *Forms) error {
	fpos := *pos;
	var forms Forms

	for {
		c, _, err := in.ReadRune()
		
		if err != nil {
			if err == io.EOF {
				return self.E(pos, "Open group")
			}
			
			return err
		}

		if c == ')' {
			pos.column++
			break
		} else {
			in.UnreadRune()
		}

		if err := self.ReadForm(pos, in, &forms); err != nil {
			if err == io.EOF {
				return self.E(pos, "Open group")
			}

			return err
		}
	}

	out.Push(NewGroupForm(fpos, forms.items...))
	return nil
}

func (self *Vm) ReadId(pos *Pos, in *bufio.Reader, out *Forms) error {
	var buffer strings.Builder
	fpos := *pos
	
	for {
		c, _, err := in.ReadRune()
		
		if err != nil {
			if err == io.EOF {
				break
			}
			
			return err
		}

		if c == '(' || c == ')' || unicode.IsSpace(c) || unicode.IsControl(c) {
			in.UnreadRune()
			break
		}
		
		buffer.WriteRune(c)
		pos.column++
	}

	out.Push(NewIdForm(fpos, buffer.String()))
	return nil
}

func (self *Vm) ReadInt(pos *Pos, in *bufio.Reader, out *Forms) error {
	var v int
	base := 10
	fpos := *pos
	
	for {
		
		c, _, err := in.ReadRune()
		
		if err != nil {
			if err == io.EOF {
				break
			}
		
			return err
		}

		if !unicode.IsDigit(c) && (base != 16 || c < 'a' || c > 'f') {
			if err = in.UnreadRune(); err != nil {
				return err
			}
			
			break
		}
		
		var dv int
		
		if base == 16 && c >= 'a' && c <= 'f' {
			dv = 10 + int(c) - int('a')
		} else {
			dv = int(c) - int('0')
		}
		
		v = v * base + dv
		pos.column++
	}
	
	out.Push(NewLitForm(fpos, &self.AbcLib.IntType, v))
	return nil
}

func (self *Vm) ReadString(pos *Pos, in *bufio.Reader, out *Forms) error {
	fpos := *pos
	var buf strings.Builder
	
	for {
		
		c, _, err := in.ReadRune()
		
		if err != nil {
			if err == io.EOF {
				return self.E(pos, "Open string")
			}
		
			return err
		}

		if c == '"' {
			break
		}
		
		buf.WriteRune(c)
		pos.column++
	}
	
	out.Push(NewLitForm(fpos, &self.AbcLib.StringType, buf.String()))
	return nil
}

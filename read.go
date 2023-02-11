package snabl

import (
	"bufio"
	//"fmt"
	"io"
	"strings"
	"unicode"
)

type Reader = func(vm *Vm, pos *Pos, in *bufio.Reader, out *Forms) error

func ReadForms(vm *Vm, pos *Pos, in *bufio.Reader, out *Forms) error {
	for {
		if err := ReadForm(vm, pos, in, out); err != nil {
			if err == io.EOF {
				return nil
			}

			return err
		}
	}
}

func ReadForm(vm *Vm, pos *Pos, in *bufio.Reader, out *Forms) error {
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
		return ReadGroup(vm, pos, in, out)
	default:
		if unicode.IsDigit(c) {
			in.UnreadRune()
			return ReadInt(vm, pos, in, out)
		} else if !unicode.IsSpace(c) && !unicode.IsControl(c) {
			in.UnreadRune()
			return ReadId(vm, pos, in, out)			
		}
	}

	return vm.E(pos, "%v?", c)
}

func ReadGroup(vm *Vm, pos *Pos, in *bufio.Reader, out *Forms) error {
	fpos := *pos;
	var forms Forms

	for {
		c, _, err := in.ReadRune()
		
		if err != nil {
			if err == io.EOF {
				return vm.E(pos, "Open group")
			}
			
			return err
		}

		if c == ')' {
			pos.column++
			break
		} else {
			in.UnreadRune()
		}

		if err := ReadForm(vm, pos, in, &forms); err != nil {
			if err == io.EOF {
				return vm.E(pos, "Open group")
			}

			return err
		}
	}

	out.Push(NewGroupForm(fpos, forms.items...))
	return nil
}

func ReadId(vm *Vm, pos *Pos, in *bufio.Reader, out *Forms) error {
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

func ReadInt(vm *Vm, pos *Pos, in *bufio.Reader, out *Forms) error {
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
	
	out.Push(NewLitForm(fpos, &vm.AbcLib.IntType, v))
	return nil
}

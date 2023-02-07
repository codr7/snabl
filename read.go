package snabl

import (
	"bufio"
	"io"
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
	default:
		if unicode.IsDigit(c) {
			in.UnreadRune()
			return ReadInt(vm, pos, in, out)
		}
	}

	return vm.E(pos, "%v?", c)
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

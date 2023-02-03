package snabl

import (
	"errors"
	"fmt"
)

func NewE(pos *Pos, spec string, args...interface{}) error {
	msg := fmt.Sprintf("Error in %v@%v:%v %v", 
		pos.source, pos.line, pos.column, fmt.Sprintf(spec, args...))

	return errors.New(msg)
}

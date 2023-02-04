package snabl

import (
	"errors"
	"fmt"
	"strings"
)

func NewE(pos *Pos, spec string, args...interface{}) error {
	var msg strings.Builder

	if pos != nil {
		fmt.Fprintf(&msg, "%v ", pos)
	}
	
	fmt.Fprintf(&msg, spec, args...)
	return errors.New(msg.String())
}

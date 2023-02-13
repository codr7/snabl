package snabl

import (
	"io"
)

type SeqType interface {
	Iter(val V) Iter
}

type LenType interface {
	SeqType
	Len(val V) int
}

type Iter interface {
	Next() (*V, error)
	Dump(out io.Writer) error
}

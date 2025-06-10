package lexer

import "errors"

var (
	ErrRealBeforeComplex  = errors.New("real part of a complex number has to be defined before the imaginary part e.g. `3+2i` not `2i+3`")
	ErrBadSyntaxForNumber = errors.New("bad syntax for number")
	ErrImagMissing        = errors.New("missing imaginary symbol `i`")
)

package parser

type MarkerKind int

const (
	MarkerKindList = iota + 1
	MarkerKindFloat
	MarkerKindInt
	MarkerKindComplex
	MarkerKindBool
	MarkerKindString
)

//go:generate stringer -type=MarkerKind

package parser

type MarkerKind int

const (
	MarkerKindString MarkerKind = iota + 1
	MarkerKindFloat
	MarkerKindInt
	MarkerKindComplex
	MarkerKindBool
	MarkerKindList
)

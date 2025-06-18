//go:generate stringer -type=Kind

package marker

type Kind int

const (
	STRING Kind = iota + 1
	FLOAT
	INT
	COMPLEX
	BOOL
	LIST
)

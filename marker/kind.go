//go:generate stringer -type=Kind

package marker

type Kind int

const (
	INVALID Kind = iota
	STRING
	FLOAT
	INT
	COMPLEX
	BOOL
	LIST
)

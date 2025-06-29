//go:generate stringer -type=Kind

package token

type Kind int

const (
	EOF Kind = iota + 1
	ERROR
	STRING  // "codemark" (without the quotation marks)
	BOOL    // true or false
	IDENT   // codemark:valid:ident
	ASSIGN  // =
	PLUS    // +
	LBRACK  // [
	RBRACK  // ]
	INT     // 1237123, 0x283f etc.
	FLOAT   // 1.2
	COMPLEX // 3 + 2i
	COMMA   // ,
)

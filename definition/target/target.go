//go:generate stringer -type=Target

package target

// Target defines to which type of
// expression a marker is appliable
type Target int

const (
	FIELD Target = iota + 1
	NAMED
	PKG // Package
	FUNC
	CONST
	VAR
	METHOD
	IFACE
	IMPORT
	ALIAS
	IFACESIG // Interface Signature
	STRUCT
	ANY
)

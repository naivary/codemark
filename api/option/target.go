//go:generate stringer -type=Target

package option

// Target defines to which type of
// expression an option can be applied
type Target int

const (
	TargetField Target = iota + 1
	TargetNamed
	TargetPkg // Package
	TargetFunc
	TargetConst
	TargetVar
	TargetMethod
	TargetIface
	TargetImport // Import statmenet
	TargetAlias
	TargetIfaceSig // Interface Signature
	TargetStruct
	TargetAny
)

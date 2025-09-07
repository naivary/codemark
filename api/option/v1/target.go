package v1

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

func (t Target) String() string {
	switch t {
	case TargetField:
		return "Field"
	case TargetNamed:
		return "Named"
	case TargetPkg:
		return "Pkg"
	case TargetFunc:
		return "Func"
	case TargetConst:
		return "Const"
	case TargetVar:
		return "Var"
	case TargetMethod:
		return "Method"
	case TargetIface:
		return "Iface"
	case TargetImport:
		return "Import"
	case TargetAlias:
		return "Alias"
	case TargetIfaceSig:
		return "IfaceSig"
	case TargetStruct:
		return "Struct"
	case TargetAny:
		return "Any"
	default:
		return "Unknown"
	}
}

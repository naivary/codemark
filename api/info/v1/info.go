package v1

import (
	"go/types"

	"golang.org/x/tools/go/packages"
)

type Info interface {
	Options() Options
}

type Project = map[*packages.Package]*Information

type Filename = string

type Information struct {
	Structs map[types.Object]*StructInfo
	Ifaces  map[types.Object]*IfaceInfo
	Aliases map[types.Object]*AliasInfo
	Named   map[types.Object]*NamedInfo
	Consts  map[types.Object]*ConstInfo
	Vars    map[types.Object]*VarInfo
	Imports map[types.Object]*ImportInfo
	Funcs   map[types.Object]*FuncInfo
	Files   map[Filename]*FileInfo
}

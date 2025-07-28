package loader

import (
	"errors"
	"go/types"

	"golang.org/x/tools/go/packages"

	infov1 "github.com/naivary/codemark/api/info/v1"
)

var ErrPkgsEmpty = errors.New(
	"loaded packages are empty. check that the defined patterns are matching any files",
)

type Loader interface {
	Load(patterns ...string) (map[*packages.Package]*infov1.Information, error)
}

func newInformation() *infov1.Information {
	return &infov1.Information{
		Structs: make(map[types.Object]*infov1.StructInfo),
		Ifaces:  make(map[types.Object]infov1.IfaceInfo),
		Aliases: make(map[types.Object]infov1.AliasInfo),
		Named:   make(map[types.Object]*infov1.NamedInfo),
		Consts:  make(map[types.Object]infov1.ConstInfo),
		Vars:    make(map[types.Object]infov1.VarInfo),
		Imports: make(map[types.Object]infov1.ImportInfo),
		Funcs:   make(map[types.Object]infov1.FuncInfo),
		Files:   make(map[infov1.Filename]infov1.FileInfo),
	}
}

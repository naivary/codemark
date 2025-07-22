package loader

import (
	"errors"
	"go/types"

	"golang.org/x/tools/go/packages"

	loaderapi "github.com/naivary/codemark/api/loader"
)

var ErrPkgsEmpty = errors.New(
	"loaded packages are empty. check that the defined patterns are matching any files",
)

type Loader interface {
	Load(patterns ...string) (map[*packages.Package]*loaderapi.Information, error)
}

func newInformation() *loaderapi.Information {
	return &loaderapi.Information{
		Structs: make(map[types.Object]*loaderapi.StructInfo),
		Ifaces:  make(map[types.Object]loaderapi.IfaceInfo),
		Aliases: make(map[types.Object]loaderapi.AliasInfo),
		Named:   make(map[types.Object]*loaderapi.NamedInfo),
		Consts:  make(map[types.Object]loaderapi.ConstInfo),
		Vars:    make(map[types.Object]loaderapi.VarInfo),
		Imports: make(map[types.Object]loaderapi.ImportInfo),
		Funcs:   make(map[types.Object]loaderapi.FuncInfo),
		Files:   make(map[loaderapi.Filename]loaderapi.FileInfo),
	}
}

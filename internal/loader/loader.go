package loader

import (
	"errors"
	"go/types"

	"golang.org/x/tools/go/packages"

	loaderv1 "github.com/naivary/codemark/api/loader/v1"
)

var ErrPkgsEmpty = errors.New(
	"loaded packages are empty. check that the defined patterns are matching any files",
)

type Loader interface {
	Load(patterns ...string) (map[*packages.Package]*loaderv1.Information, error)
}

func newInformation() *loaderv1.Information {
	return &loaderv1.Information{
		Structs: make(map[types.Object]*loaderv1.StructInfo),
		Ifaces:  make(map[types.Object]loaderv1.IfaceInfo),
		Aliases: make(map[types.Object]loaderv1.AliasInfo),
		Named:   make(map[types.Object]*loaderv1.NamedInfo),
		Consts:  make(map[types.Object]loaderv1.ConstInfo),
		Vars:    make(map[types.Object]loaderv1.VarInfo),
		Imports: make(map[types.Object]loaderv1.ImportInfo),
		Funcs:   make(map[types.Object]loaderv1.FuncInfo),
		Files:   make(map[loaderv1.Filename]loaderv1.FileInfo),
	}
}

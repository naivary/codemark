package codemark

import (
	"go/ast"
	"go/parser"
	"go/token"

	"golang.org/x/tools/go/packages"
)

type Workspace struct{}

type Loader interface {
	Load(patterns ...string) (*Workspace, error)
}

func NewLoader(mngr *ConverterManager, cfg *packages.Config) Loader {
	l := &loader{
		mngr: mngr,
	}
	if cfg == nil {
		l.cfg = l.defaultConfig()
	}
	l.cfg.ParseFile = func(fset *token.FileSet, filename string, src []byte) (*ast.File, error) {
		return parser.ParseFile(fset, filename, src, parser.ParseComments)
	}
	return l
}

var _ Loader = (*loader)(nil)

type loader struct {
	mngr *ConverterManager
	cfg  *packages.Config
}

func (l *loader) Load(patterns ...string) (*Workspace, error) {
	return nil, nil
}

func (l *loader) defaultConfig() *packages.Config {
	return &packages.Config{
		Mode: packages.NeedTypesInfo | packages.NeedSyntax | packages.NeedTypes,
	}
}

package codemark

import (
	"go/ast"
	"go/parser"
	"go/token"

	"github.com/naivary/codemark/sdk"
	"golang.org/x/tools/go/packages"
)

func NewLocalLoader(mngr *ConverterManager, cfg *packages.Config) sdk.Loader {
	l := &localLoader{
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

var _ sdk.Loader = (*localLoader)(nil)

type localLoader struct {
	mngr *ConverterManager
	cfg  *packages.Config
}

func (l *localLoader) Load(patterns ...string) (*sdk.Project, error) {
	return nil, nil
}

func (l *localLoader) defaultConfig() *packages.Config {
	return &packages.Config{
		Mode: packages.NeedTypesInfo | packages.NeedSyntax | packages.NeedTypes,
	}
}

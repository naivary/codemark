package codemark

import (
	"github.com/naivary/codemark/converter"
	"github.com/naivary/codemark/loader"
	"github.com/naivary/codemark/sdk"
	"golang.org/x/tools/go/packages"
)

var NewLoader = loader.New

// LoadWithManager will load a package with a custom manager which usually mean
// you have custom converter to add.
func LoadWithManager(mngr sdk.ConverterManager, patterns ...string) (map[*packages.Package]*sdk.Project, error) {
	l := loader.New(mngr, nil)
	return l.Load(patterns...)
}

// Load is extracting all the type informations including, while parsing the
// markers as definitions defnied in the provided registry.
func Load(reg sdk.Registry, patterns ...string) (map[*packages.Package]*sdk.Project, error) {
	mngr, err := converter.NewManager(reg)
	if err != nil {
		return nil, err
	}
	return LoadWithManager(mngr, patterns...)
}

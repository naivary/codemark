package codemark

import (
	"fmt"

	loaderapi "github.com/naivary/codemark/api/loader"
	"github.com/naivary/codemark/converter"
	"github.com/naivary/codemark/loader"
	"github.com/naivary/codemark/registry"
	"golang.org/x/tools/go/packages"
)

// NewLoader will returns a new loader. This function should be only used if you
// need fine grained control over the configuration and options of the loader.
// Otherwise use `LoadWithManager` or `Load`.
var NewLoader = loader.New

// LoadWithManager will load a package with a custom manager which usually mean
// you have custom converter to add.
func LoadWithManager(mngr *converter.Manager, patterns ...string) (map[*packages.Package]*loaderapi.Project, error) {
	l := loader.New(mngr, nil)
	if len(patterns) == 0 {
		return nil, fmt.Errorf("patterns cannot be empty because no projects can be loaded")
	}
	return l.Load(patterns...)
}

// Load is extracting all the type informations including, while parsing the
// markers as definitions defnied in the provided registry.
func Load(reg registry.Registry, patterns ...string) (map[*packages.Package]*loaderapi.Project, error) {
	mngr, err := converter.NewManager(reg)
	if err != nil {
		return nil, err
	}
	return LoadWithManager(mngr, patterns...)
}

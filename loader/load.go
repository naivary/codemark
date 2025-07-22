package loader

import (
	"fmt"

	"golang.org/x/tools/go/packages"

	loaderapi "github.com/naivary/codemark/api/loader"
	"github.com/naivary/codemark/internal/converter"
	"github.com/naivary/codemark/registry"
)

// LoadWithManager will load a package with a custom manager which usually mean
// you have custom converter to add.
func LoadWithManager(
	mngr *converter.Manager,
	patterns ...string,
) (map[*packages.Package]*loaderapi.Information, error) {
	l := New(mngr, nil)
	if len(patterns) == 0 {
		return nil, fmt.Errorf("patterns cannot be empty because no projects can be loaded")
	}
	return l.Load(patterns...)
}

// Load is extracting all the type informations including, while parsing the
// found markers.
func Load(
	reg registry.Registry,
	patterns ...string,
) (map[*packages.Package]*loaderapi.Information, error) {
	mngr, err := converter.NewManager(reg)
	if err != nil {
		return nil, err
	}
	return LoadWithManager(mngr, patterns...)
}

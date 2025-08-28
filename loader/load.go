package loader

import (
	"fmt"

	"golang.org/x/tools/go/packages"

	convv1 "github.com/naivary/codemark/api/converter/v1"
	infov1 "github.com/naivary/codemark/api/info/v1"
	regv1 "github.com/naivary/codemark/api/registry/v1"
	"github.com/naivary/codemark/converter"
	"github.com/naivary/codemark/internal/loader"
)

// Load is extracting all the type informations including, while parsing the
// found markers.
func Load(reg regv1.Registry, convs []convv1.Converter, patterns ...string) (map[*packages.Package]*infov1.Information, error) {
	mngr, err := converter.NewManager(reg, convs...)
	if err != nil {
		return nil, err
	}
	l := loader.New(mngr, nil)
	if len(patterns) == 0 {
		return nil, fmt.Errorf("patterns cannot be empty because no projects can be loaded")
	}
	return l.Load(patterns...)
}

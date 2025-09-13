package codemark

import (
	"go/types"

	"golang.org/x/tools/go/packages"

	genv1 "github.com/naivary/codemark/api/generator/v1"
	infov1 "github.com/naivary/codemark/api/info/v1"
	optv1 "github.com/naivary/codemark/api/option/v1"
)

type Resourcer interface {
	// Resouce represented by this resource
	Resource() string

	// Options of the resource
	Options() []*optv1.Option

	// Create generated the actual artifact
	Create(pkg *packages.Package, info map[types.Object][]infov1.Info) (*genv1.Artifact, error)
}

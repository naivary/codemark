package codemark

import (
	"go/types"

	"golang.org/x/tools/go/packages"

	genv1 "github.com/naivary/codemark/api/generator/v1"
	infov1 "github.com/naivary/codemark/api/info/v1"
	optv1 "github.com/naivary/codemark/api/option/v1"
)

const _optDocResource = "option"

var _ Resourcer = (*optDocResourcer)(nil)

type optDocResourcer struct{}

// Resouce represented by this resource
func (o optDocResourcer) Resource() string {
	return _optDocResource
}

// Options of the resource
func (o optDocResourcer) Options() []*optv1.Option {
	return makeOpts(_optDocResource,
		mustMakeOpt("description", OptDocDesc(""), _unique, optv1.TargetNamed),
		mustMakeOpt("summary", OptDocDesc(""), _unique, optv1.TargetNamed),
	)
}

func (o optDocResourcer) CanCreate(info infov1.Info) bool {
	_, isNamed := info.(*infov1.NamedInfo)
	opts := info.Options()
	return opts.IsDefined("codemark:option:summary") || opts.IsDefined("codemark:option:description") && isNamed
}

// Create generated the actual artifact
func (o optDocResourcer) Create(pkg *packages.Package, obj types.Object, info infov1.Info) (*genv1.Artifact, error) {
	// we need all of the nameds to write it down to a file
	named := info.(*infov1.NamedInfo)
	return nil, nil
}

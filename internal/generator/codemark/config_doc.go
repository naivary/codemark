package codemark

import (
	"go/types"

	"golang.org/x/tools/go/packages"

	docv1 "github.com/naivary/codemark/api/doc/v1"
	genv1 "github.com/naivary/codemark/api/generator/v1"
	infov1 "github.com/naivary/codemark/api/info/v1"
	optv1 "github.com/naivary/codemark/api/option/v1"
)

const _configDocResource = "config"

type configDocResourcer struct{}

func NewConfigDocResourcer() (*configDocResourcer, error) {
	return &configDocResourcer{}, nil
}

func (c configDocResourcer) Options() []*optv1.Option {
	return makeOpts(_configDocResource,
		mustMakeOpt("description", ConfigDocDescription(""), _unique, optv1.TargetField),
		mustMakeOpt("default", ConfigDocDefault(""), _unique, optv1.TargetField),
	)
}

func (c configDocResourcer) CanCreate(info infov1.Info) bool {
	s, isStruct := info.(*infov1.StructInfo)
	if !isStruct {
		return false
	}
	for _, field := range s.Fields {
		opts := field.Options().Filter(_domain, _configDocResource)
		if opts.IsDefined("codemark:config:description") {
			return true
		}
	}
	return false
}

func (c configDocResourcer) Create(pkg *packages.Package, obj types.Object, info infov1.Info) (*genv1.Artifact, error) {
	if !c.CanCreate(info) {
		// just cannot create it. It's not really an error
		return nil, nil
	}
	// TODO: recursively add structs, maps
	s := info.(*infov1.StructInfo)
	configDocs := make(map[string]docv1.Config)
	for _, field := range s.Fields {
		configDoc := docv1.Config{}
		fieldName := field.Ident.Name
		for _, opts := range field.Options().Filter(_domain, _configDocResource) {
			for _, opt := range opts {
				switch v := opt.(type) {
				case ConfigDocDefault:
					configDoc.Default = string(v)
				case ConfigDocDescription:
					configDoc.Description = string(v)
				}
			}
		}
		configDocs[fieldName] = configDoc
	}
	return nil, nil
}

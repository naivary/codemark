package k8s

import (
	"github.com/naivary/codemark/api/doc"
	optionapi "github.com/naivary/codemark/api/option"
)

const _serviceAccountResource = "serviceaccount"

func serviceAccountOpts() []*optionapi.Option {
	return makeOpts(_serviceAccountResource,
		mustMakeOpt("name", ServiceAccountName(""), true, optionapi.TargetFunc),
	)
}

type ServiceAccountName string

func (s ServiceAccountName) Doc() doc.Option {
	return doc.Option{
		Desc:    "Name of the generated ServiceAccount name",
		Default: "REQUIRED",
	}
}

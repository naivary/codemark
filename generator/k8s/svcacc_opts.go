package k8s

import (
	docv1 "github.com/naivary/codemark/api/doc/v1"
	optionv1 "github.com/naivary/codemark/api/option/v1"
)

const _serviceAccountResource = "serviceaccount"

func serviceAccountOpts() []*optionv1.Option {
	return makeOpts(_serviceAccountResource,
		mustMakeOpt("name", ServiceAccountName(""), _required, true, optionv1.TargetFunc),
	)
}

type ServiceAccountName string

func (s ServiceAccountName) Doc() docv1.Option {
	return docv1.Option{
		Desc:    "Name of the generated ServiceAccount name",
		Default: "REQUIRED",
	}
}

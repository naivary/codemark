package k8s

import (
	"errors"

	docv1 "github.com/naivary/codemark/api/doc/v1"
	optionv1 "github.com/naivary/codemark/api/option/v1"

	corev1 "k8s.io/api/core/v1"
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

func (s ServiceAccountName) apply(svcacc *corev1.ServiceAccount) error {
	v := string(s)
	if len(v) == 0 {
		return errors.New("service account name cannot be empty")
	}
	svcacc.Name = v
	return nil
}

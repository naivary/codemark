package k8s

import optionapi "github.com/naivary/codemark/api/option"

const _serviceAccountResource = "serviceaccount"

func serviceAccountOpts() []*optionapi.Option {
	return makeOpts(_serviceAccountResource,
		newOption(Immutable(false), true, optionapi.TargetStruct),
	)
}

type ServiceAccountName string

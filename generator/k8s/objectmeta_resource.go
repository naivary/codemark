package k8s

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	loaderapi "github.com/naivary/codemark/api/loader"
	optionapi "github.com/naivary/codemark/api/option"
)

const _metaResource = "meta"

func objectMetaOpts() []*optionapi.Option {
	return makeOpts(_metaResource,
		mustMakeOpt(_typeName, Name(""), true, optionapi.TargetAny),
		mustMakeOpt(_typeName, Namespace(""), true, optionapi.TargetAny),
	)
}

func createObjectMeta(optioner loaderapi.Optioner) (metav1.ObjectMeta, error) {
	obj := metav1.ObjectMeta{}
	for _, opts := range optioner.Options() {
		for _, opt := range opts {
			var err error
			switch o := opt.(type) {
			case Name:
				err = o.apply(&obj)
			case Namespace:
				err = o.apply(&obj)
			}
			if err != nil {
				return obj, err
			}
		}
	}
	return obj, nil
}

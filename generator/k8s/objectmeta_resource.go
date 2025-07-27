package k8s

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	loaderv1 "github.com/naivary/codemark/api/loader/v1"
	optionv1 "github.com/naivary/codemark/api/option/v1"
)

const _metaResource = "meta"

func objectMetaOpts() []*optionv1.Option {
	return makeOpts(_metaResource,
		mustMakeOpt(_typeName, Name(""), _optional, _unique, optionv1.TargetAny),
		mustMakeOpt(_typeName, Namespace(""), _optional, _unique, optionv1.TargetAny),
	)
}

func createObjectMeta(info loaderv1.Info) (metav1.ObjectMeta, error) {
	obj := metav1.ObjectMeta{}
	for _, opts := range info.Options() {
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

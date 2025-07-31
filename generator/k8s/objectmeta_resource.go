package k8s

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	infov1 "github.com/naivary/codemark/api/info/v1"
	optionv1 "github.com/naivary/codemark/api/option/v1"
)

func setObjectMetaDefaults(info infov1.Info, cfg *config) {
	opts := objectMetaOpts()
	setDefaults(opts, info, cfg, optionv1.TargetAny)
}

func createObjectMeta(info infov1.Info, cfg *config) (metav1.ObjectMeta, error) {
	obj := metav1.ObjectMeta{}
	setObjectMetaDefaults(info, cfg)
	format := info.Options()["k8s:metadata:format.name"][0].(Format)
	for _, opts := range info.Options() {
		for _, opt := range opts {
			var err error
			switch o := opt.(type) {
			case Name:
				err = o.apply(&obj, format)
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

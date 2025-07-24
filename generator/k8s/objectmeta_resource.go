package k8s

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	loaderapi "github.com/naivary/codemark/api/loader"
	optionapi "github.com/naivary/codemark/api/option"
)

func objectMetaOpts() []*optionapi.Option {
	const resource = "meta"
	return makeDefs(resource, map[any][]optionapi.Target{
		Name(""):      {optionapi.TargetAny},
		Namespace(""): {optionapi.TargetAny},
	})
}

func createObjectMeta(optioner loaderapi.Optioner) (*metav1.ObjectMeta, error) {
	obj := &metav1.ObjectMeta{}
	for _, opts := range optioner.Options() {
		for _, opt := range opts {
			var err error
			switch o := opt.(type) {
			case Name:
				err = o.apply(obj)
			case Namespace:
				err = o.apply(obj)
			}
			if err != nil {
				return nil, err
			}
		}
	}
	return obj, nil
}

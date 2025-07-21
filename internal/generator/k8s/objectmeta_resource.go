package k8s

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/naivary/codemark/api/core"
	loaderapi "github.com/naivary/codemark/api/loader"
)

func objectMetaOpts() []*core.Option {
	const resource = "meta"
	return makeDefs(resource, map[any][]core.Target{
		Name(""):      {core.TargetAny},
		Namespace(""): {core.TargetAny},
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

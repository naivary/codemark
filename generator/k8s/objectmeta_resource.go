package k8s

import (
	"github.com/naivary/codemark/api/core"
	loaderapi "github.com/naivary/codemark/api/loader"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func objectMetaOpts() []*core.Option {
	const resource = "meta"
	return makeDefs(resource, map[any][]core.Target{
		Name(""):      {core.TargetAny},
		Namespace(""): {core.TargetAny},
	})
}

func createObjectMeta(optioner loaderapi.Optioner) metav1.ObjectMeta {
	obj := metav1.ObjectMeta{}
	for _, opts := range optioner.Options() {
		for _, opt := range opts {
			switch o := opt.(type) {
			case Name:
				o.apply(&obj)
			case Namespace:
				o.apply(&obj)
			}
		}
	}
	return obj
}

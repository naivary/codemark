package k8s

import (
	"github.com/naivary/codemark/api"
	loaderapi "github.com/naivary/codemark/api/loader"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func objectMetaDefs() []*api.Definition {
	const resource = "meta"
	return makeDefs(resource,
		Name(""),
		Namespace(""),
	)
}

func createObjectMeta(d loaderapi.Defs) metav1.ObjectMeta {
	obj := metav1.ObjectMeta{}
	for _, defs := range d.Definitions() {
		for _, def := range defs {
			switch d := def.(type) {
			case Name:
				d.apply(&obj)
			case Namespace:
				d.apply(&obj)
			}
		}
	}
	return obj
}

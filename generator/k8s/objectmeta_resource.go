package k8s

import (
	"reflect"

	loaderapi "github.com/naivary/codemark/api/loader"
	"github.com/naivary/codemark/definition"
	"github.com/naivary/codemark/definition/target"
	"github.com/naivary/codemark/maker"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func objectMetaDefs() []*definition.Definition {
	return []*definition.Definition{
		maker.MustMakeDefWithDoc(
			objectMetaIdent("name"),
			reflect.TypeFor[Name](),
			"name of the object",
			target.STRUCT,
		),
		maker.MustMakeDefWithDoc(
			objectMetaIdent("namespace"),
			reflect.TypeFor[Namespace](),
			"namespace for the object. It does not check if the object is namespace scoped",
			target.STRUCT,
		),
	}
}

func createObjectMeta(strc *loaderapi.StructInfo) *metav1.ObjectMeta {
	obj := &metav1.ObjectMeta{}
	for _, defs := range strc.Defs {
		for _, def := range defs {
			switch d := def.(type) {
			case Name:
				d.apply(obj)
			case Namespace:
				d.apply(obj)
			}
		}
	}
	return obj
}

package k8s

import (
	"reflect"

	loaderapi "github.com/naivary/codemark/api/loader"
	"github.com/naivary/codemark/definition"
	"github.com/naivary/codemark/definition/target"
	"github.com/naivary/codemark/maker"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func newConfigMap() *corev1.ConfigMap {
	return &corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ConfigMap",
			APIVersion: "v1",
		},
		Data: make(map[string]string),
	}
}

func configMapDefs() []*definition.Definition {
	return []*definition.Definition{
		// structs
		maker.MustMakeDefWithDoc(cmIdent("immutable"), reflect.TypeFor[Immutable](), "decides wether the configmap should be immutable", target.STRUCT),

		// fields
		maker.MustMakeDefWithDoc(cmIdent("default"), reflect.TypeFor[Default](), "default value for the config map value", target.FIELD),
	}
}

func createConfigMap(strc *loaderapi.StructInfo) (*corev1.ConfigMap, error) {
	cm := newConfigMap()
	cm.ObjectMeta = *createObjectMeta(strc)
	for _, defs := range strc.Defs {
		for _, def := range defs {
			switch d := def.(type) {
			case Immutable:
				d.apply(cm)
			}
		}
	}

	for _, field := range strc.Fields {
		for _, defs := range field.Defs {
			for _, def := range defs {
				switch d := def.(type) {
				case Default:
					d.apply(field.Idn.Name, cm)
				}
			}
		}
	}
	return cm, nil
}

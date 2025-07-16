package k8s

import (
	"slices"
	"strings"

	"github.com/naivary/codemark/api"
	loaderapi "github.com/naivary/codemark/api/loader"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func newConfigMap() *corev1.ConfigMap {
	return &corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "ConfigMap",
		},
		Data: make(map[string]string),
	}
}

func setDataInConfigMap(key, value string, cm *corev1.ConfigMap) {
	lower := strings.ToLower(key)
	cm.Data[lower] = value
}

func configMapDefs() []*api.Definition {
	const resource = "configmap"
	return makeDefs(resource,
		Immutable(false),
		Default(""),
	)
}

func createConfigMap(strc *loaderapi.StructInfo) (*corev1.ConfigMap, error) {
	cm := newConfigMap()
	cm.ObjectMeta = createObjectMeta(strc)
	for _, defs := range strc.Defs {
		for _, def := range defs {
			switch d := def.(type) {
			case Immutable:
				d.apply(cm)
			}
		}
	}
	for _, field := range strc.Fields {
		idents := keys(field.Defs)
		if !slices.Contains(idents, "k8s:configmap:default") {
			setDataInConfigMap(field.Idn.Name, "", cm)
			continue
		}
		for _, defs := range field.Defs {
			for _, def := range defs {
				switch d := def.(type) {
				case Default:
					d.apply(field.Idn, cm)
				}
			}
		}
	}
	return cm, nil
}

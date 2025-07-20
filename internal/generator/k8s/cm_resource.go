package k8s

import (
	"slices"
	"strings"

	coreapi "github.com/naivary/codemark/api/core"
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

func configMapOpts() []*coreapi.Option {
	const resource = "configmap"
	return makeDefs(resource, map[any][]coreapi.Target{
		Immutable(false): {coreapi.TargetStruct},
		Default(""):      {coreapi.TargetField},
	})
}

func createConfigMap(strc *loaderapi.StructInfo) (*corev1.ConfigMap, error) {
	cm := newConfigMap()
	cm.ObjectMeta = createObjectMeta(strc)
	for _, opts := range strc.Opts {
		for _, opt := range opts {
			switch o := opt.(type) {
			case Immutable:
				o.apply(cm)
			}
		}
	}
	for _, field := range strc.Fields {
		idents := keys(field.Opts)
		if !slices.Contains(idents, "k8s:configmap:default") {
			setDataInConfigMap(field.Idn.Name, "", cm)
			continue
		}
		for _, opts := range field.Opts {
			for _, opt := range opts {
				switch o := opt.(type) {
				case Default:
					o.apply(field.Idn, cm)
				}
			}
		}
	}
	return cm, nil
}

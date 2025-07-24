package k8s

import (
	"slices"
	"strings"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	loaderapi "github.com/naivary/codemark/api/loader"
	optionapi "github.com/naivary/codemark/api/option"
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

func configMapOpts() []*optionapi.Option {
	const resource = "configmap"
	return makeDefs(resource,
		newOption(Immutable(false), true, optionapi.TargetStruct),
		newOption(Default(""), true, optionapi.TargetField),
	)
}

func createConfigMap(strc *loaderapi.StructInfo) (*corev1.ConfigMap, error) {
	cm := newConfigMap()
	objectMeta, err := createObjectMeta(strc)
	if err != nil {
		return nil, err
	}
	cm.ObjectMeta = *objectMeta
	for _, opts := range strc.Opts {
		for _, opt := range opts {
			var err error
			switch o := opt.(type) {
			case Immutable:
				err = o.apply(cm)
			}
			if err != nil {
				return nil, err
			}
		}
	}
	for _, field := range strc.Fields {
		idents := keys(field.Opts)
		if !slices.Contains(idents, "k8s:configmap:default") {
			setDataInConfigMap(field.Ident.Name, "", cm)
			continue
		}
		for _, opts := range field.Opts {
			for _, opt := range opts {
				var err error
				switch o := opt.(type) {
				case Default:
					err = o.apply(field.Ident, cm)
				}
				if err != nil {
					return nil, err
				}
			}
		}
	}
	return cm, nil
}

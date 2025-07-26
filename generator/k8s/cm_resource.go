package k8s

import (
	"bytes"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/goccy/go-yaml"

	genv1 "github.com/naivary/codemark/api/generator/v1"
	loaderv1 "github.com/naivary/codemark/api/loader/v1"
)

func newConfigMap(strc *loaderv1.StructInfo) (corev1.ConfigMap, error) {
	cm := corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "ConfigMap",
		},
		Data: make(map[string]string),
	}
	objectMeta, err := createObjectMeta(strc)
	if err != nil {
		return cm, err
	}
	cm.ObjectMeta = objectMeta
	return cm, nil
}

func createConfigMap(strc *loaderv1.StructInfo) (*genv1.Artifact, error) {
	cm, err := newConfigMap(strc)
	if err != nil {
		return nil, err
	}
	format, err := applyStructOptsToConfigMap(strc, &cm)
	if err != nil {
		return nil, err
	}
	for _, finfo := range strc.Fields {
		if err := applyFieldOptToConfigMap(finfo, format, &cm); err != nil {
			return nil, err
		}
	}
	var file bytes.Buffer
	if err := yaml.NewEncoder(&file).Encode(&cm); err != nil {
		return nil, err
	}
	return &genv1.Artifact{
		Name:        "codemark_k8s_configmap",
		ContentType: "text/yaml",
		Data:        &file,
	}, nil
}

func applyStructOptsToConfigMap(strc *loaderv1.StructInfo, cm *corev1.ConfigMap) (KeyFormat, error) {
	format := CamelCase
	for _, opts := range strc.Opts {
		for _, opt := range opts {
			var err error
			switch o := opt.(type) {
			case Immutable:
				err = o.apply(cm)
			case KeyFormat:
				format = o
			}
			if err != nil {
				return format, err
			}
		}
	}
	return format, nil
}

func applyFieldOptToConfigMap(field loaderv1.FieldInfo, format KeyFormat, cm *corev1.ConfigMap) error {
	const defaultValue = ""
	if !field.Ident.IsExported() {
		// unexported fields will be ignored
		return nil
	}
	// set the default for an exported field to an empty string
	key := format.Format(field.Ident.Name)
	cm.Data[key] = ""

	for ident, opts := range field.Opts {
		if !isResource(ident, _configMapResource) {
			continue
		}
		for _, opt := range opts {
			var err error
			switch o := opt.(type) {
			case Default:
				err = o.apply(field, cm, format)
			}
			if err != nil {
				return err
			}
		}
	}
	if cm.Data[key] == "" && cm.Immutable != nil {
		return errImmutableConfigMapWithoutDefault
	}
	return nil
}

package k8s

import (
	"fmt"
	"strings"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	genv1 "github.com/naivary/codemark/api/generator/v1"
	infov1 "github.com/naivary/codemark/api/info/v1"
	optionv1 "github.com/naivary/codemark/api/option/v1"
)

func newConfigMap(strc *infov1.StructInfo) (corev1.ConfigMap, error) {
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

func configMapMetadataDefaults(strc *infov1.StructInfo) {
	name := strc.Opts["k8s:metadata:name"]
	if len(name) == 0 {
		cmName := strings.ToLower(strc.Spec.Name.Name)
		strc.Opts["k8s:metadata:name"] = []any{Name(cmName)}
	}
}

func configMapDefaults(strc *infov1.StructInfo) {
	opts := configMapOpts()
	setOptsDefaults(opts, strc.Opts, optionv1.TargetStruct)
	for _, finfo := range strc.Fields {
		setOptsDefaults(opts, finfo.Opts, optionv1.TargetField)
	}
	configMapMetadataDefaults(strc)
}

func createConfigMap(strc *infov1.StructInfo) (*genv1.Artifact, error) {
	configMapDefaults(strc)
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
	filename := fmt.Sprintf("%s.configmap.yaml", cm.Name)
	return newArtifact(filename, cm)
}

func applyStructOptsToConfigMap(strc *infov1.StructInfo, cm *corev1.ConfigMap) (Format, error) {
	var keyFormat Format
	for ident, opts := range strc.Opts {
		if !isResource(ident, _configMapResource) {
			continue
		}
		for _, opt := range opts {
			var err error
			switch o := opt.(type) {
			case Immutable:
				err = o.apply(cm)
			case Format:
				keyFormat = o
			}
			if err != nil {
				return keyFormat, err
			}
		}
	}
	return keyFormat, nil
}

func applyFieldOptToConfigMap(field infov1.FieldInfo, format Format, cm *corev1.ConfigMap) error {
	if !field.Ident.IsExported() {
		// unexported fields will be ignored
		return nil
	}
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
	return nil
}

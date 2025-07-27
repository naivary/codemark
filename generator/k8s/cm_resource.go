package k8s

import (
	"bytes"
	"slices"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/goccy/go-yaml"

	genv1 "github.com/naivary/codemark/api/generator/v1"
	loaderv1 "github.com/naivary/codemark/api/loader/v1"
	optionv1 "github.com/naivary/codemark/api/option/v1"
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

// TODO: write test for this function
func setOptsDefaults(opts []*optionv1.Option, infoOpts map[string][]any, targets ...optionv1.Target) {
	for _, opt := range opts {
		if !slices.Equal(opt.Targets, targets) && !slices.Contains(targets, optionv1.TargetAny) && !opt.IsRequired() {
			continue
		}
		v := infoOpts[opt.Ident]
		if len(v) > 0 {
			continue
		}
		infoOpts[opt.Ident] = append(v, opt.Default)
	}
}

func configMapDefaults(strc *loaderv1.StructInfo) {
	opts := configMapOpts()
	setOptsDefaults(opts, strc.Opts, optionv1.TargetStruct)
	for _, finfo := range strc.Fields {
		setOptsDefaults(opts, finfo.Opts, optionv1.TargetField)
	}
}

func createConfigMap(strc *loaderv1.StructInfo) (*genv1.Artifact, error) {
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

package k8s

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	genv1 "github.com/naivary/codemark/api/generator/v1"
	infov1 "github.com/naivary/codemark/api/info/v1"
	optionv1 "github.com/naivary/codemark/api/option/v1"
)

var _ Resourcer = (*configMapResourcer)(nil)

func NewConfigMapResourcer() Resourcer {
	return &configMapResourcer{resource: "configmap"}
}

type configMapResourcer struct {
	resource string
}

func (c configMapResourcer) Resource() string {
	return c.resource
}

func (c configMapResourcer) Options() []*optionv1.Option {
	return makeOpts(c.resource,
		mustMakeOpt(_typeName, Default(""), _required, _unique, optionv1.TargetField),
		mustMakeOpt(_typeName, Immutable(false), _optional, _unique, optionv1.TargetStruct),
		mustMakeOpt("format.key", Format(CamelCase), _optional, _unique, optionv1.TargetStruct),
	)
}

func (c configMapResourcer) CanCreate(info infov1.Info) bool {
	structInfo, isStructInfo := info.(*infov1.StructInfo)
	if !isStructInfo {
		return false
	}
	for _, field := range structInfo.Fields {
		for ident := range field.Opts {
			if ident == "k8s:configmap:default" {
				return true
			}
		}
	}
	return false
}

func (c configMapResourcer) Create(info infov1.Info, metadata metav1.ObjectMeta) (*genv1.Artifact, error) {
	structInfo := info.(*infov1.StructInfo)
	c.setDefaultOpts(structInfo)
	cm := c.newConfigMap(structInfo, metadata)
	format, err := c.applyStructOpts(structInfo, &cm)
	if err != nil {
		return nil, err
	}
	err = c.applyFieldOpts(structInfo, format, &cm)
	if err != nil {
		return nil, err
	}
	filename := fmt.Sprintf("%s.configmap.yaml", cm.Name)
	return newArtifact(filename, cm)
}

func (c configMapResourcer) newConfigMap(info *infov1.StructInfo, metadata metav1.ObjectMeta) corev1.ConfigMap {
	cm := corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "ConfigMap",
		},
		ObjectMeta: metadata,
		Data:       make(map[string]string, len(info.Fields)),
	}
	if cm.ObjectMeta.Name == "" {
		cm.ObjectMeta.Name = KebabCase.Format(info.Spec.Name.Name)
	}
	return cm
}

func (c configMapResourcer) setDefaultOpts(info *infov1.StructInfo) {
	opts := c.Options()
	setOptsDefaults(opts, info.Opts, optionv1.TargetStruct)
	for _, finfo := range info.Fields {
		setOptsDefaults(opts, finfo.Opts, optionv1.TargetField)
	}
}

func (c configMapResourcer) applyStructOpts(info *infov1.StructInfo, cm *corev1.ConfigMap) (Format, error) {
	var keyFormat Format
	for ident, opts := range info.Opts {
		if !isResource(ident, c.resource) {
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

func (c configMapResourcer) applyFieldOpts(info *infov1.StructInfo, format Format, cm *corev1.ConfigMap) error {
	for _, finfo := range info.Fields {
		if !finfo.Ident.IsExported() {
			continue
		}
		for ident, opts := range finfo.Opts {
			if !isResource(ident, c.resource) {
				continue
			}
			for _, opt := range opts {
				var err error
				switch o := opt.(type) {
				case Default:
					err = o.apply(finfo, cm, format)
				}
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

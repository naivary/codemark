package k8s

import (
	"go/ast"
	"strings"

	corev1 "k8s.io/api/core/v1"

	"github.com/naivary/codemark/api/doc"
	optionapi "github.com/naivary/codemark/api/option"
)

func configMapOpts() []*optionapi.Option {
	const resource = "configmap"
	return makeDefs(resource,
		newOption(Immutable(false), true, optionapi.TargetStruct),
		newOption(Default(""), true, optionapi.TargetField),
	)
}

type Default string

func (d Default) apply(field *ast.Ident, cm *corev1.ConfigMap) error {
	lower := strings.ToLower(field.Name)
	cm.Data[lower] = string(d)
	return nil
}

func (d Default) Doc() doc.Option {
	return doc.Option{
		Desc:    `Default defines the default value of the config parameter`,
		Default: "REQUIRED",
	}
}

type Immutable bool

func (i Immutable) apply(cm *corev1.ConfigMap) error {
	b := bool(i)
	cm.Immutable = &b
	return nil
}

func (i Immutable) Doc() doc.Option {
	return doc.Option{
		Desc:    `Defines if the ConfigMap is immutable or not`,
		Default: "false",
	}
}

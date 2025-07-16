package k8s

import (
	"go/ast"
	"strings"

	"github.com/naivary/codemark/api"
	"github.com/naivary/codemark/definition/target"
	corev1 "k8s.io/api/core/v1"
)

type Default string

func (d Default) apply(field *ast.Ident, cm *corev1.ConfigMap) error {
	lower := strings.ToLower(field.Name)
	cm.Data[lower] = string(d)
	return nil
}

func (d Default) Doc() api.OptionDoc {
	return api.OptionDoc{
		Targets: []target.Target{target.FIELD},
		Doc:     "Default value for the field. If no default value can be provided dont set the marker",
		Default: "<empty string>",
	}
}

type Immutable bool

func (i Immutable) apply(cm *corev1.ConfigMap) error {
	b := bool(i)
	cm.Immutable = &b
	return nil
}

func (i Immutable) Doc() api.OptionDoc {
	return api.OptionDoc{
		Targets: []target.Target{target.STRUCT},
		Doc:     "whether the ConfigMap is immutable or not",
		Default: "false",
	}
}

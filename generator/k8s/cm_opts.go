package k8s

import (
	"go/ast"
	"strings"

	corev1 "k8s.io/api/core/v1"

	optionapi "github.com/naivary/codemark/api/option"
)

type Default string

func (d Default) apply(field *ast.Ident, cm *corev1.ConfigMap) error {
	lower := strings.ToLower(field.Name)
	cm.Data[lower] = string(d)
	return nil
}

func (d Default) Doc() optionapi.OptionDoc {
	return optionapi.OptionDoc{
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

func (i Immutable) Doc() optionapi.OptionDoc {
	return optionapi.OptionDoc{
		Desc:    `Defines if the ConfigMap is immutable or not`,
		Default: "false",
	}
}

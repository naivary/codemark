package k8s

import (
	"go/ast"
	"strings"

	corev1 "k8s.io/api/core/v1"
)

type Default string

func (d Default) apply(field *ast.Ident, cm *corev1.ConfigMap) error {
	lower := strings.ToLower(field.Name)
	cm.Data[lower] = string(d)
	return nil
}

type Immutable bool

func (i Immutable) apply(cm *corev1.ConfigMap) error {
	b := bool(i)
	cm.Immutable = &b
	return nil
}

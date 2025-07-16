package k8s

import (
	"fmt"
	"strings"

	corev1 "k8s.io/api/core/v1"
)

func cmIdent(option string) string {
	return fmt.Sprintf("k8s:configmap:%s", option)
}

type Default string

func (d Default) apply(ident string, cm *corev1.ConfigMap) error {
	lower := strings.ToLower(ident)
	cm.Data[lower] = string(d)
	return nil
}

type Immutable bool

func (i Immutable) apply(cm *corev1.ConfigMap) error {
	b := bool(i)
	cm.Immutable = &b
	return nil
}

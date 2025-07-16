package k8s

import (
	"fmt"
	"strings"

	"github.com/naivary/codemark/api"
	"github.com/naivary/codemark/definition/target"
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

func (d Default) Doc() api.OptionDoc {
	return api.OptionDoc{
		Targets: []target.Target{target.FIELD},
		Doc:     "Default value for the field. If no default value can be provided dont set the marker",
	}
}

type Immutable bool

func (i Immutable) apply(cm *corev1.ConfigMap) error {
	b := bool(i)
	cm.Immutable = &b
	return nil
}

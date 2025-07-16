package k8s

import (
	"fmt"
	"strings"

	"github.com/naivary/codemark/definition/target"
	"github.com/naivary/codemark/sdk"
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

func (d Default) Doc() sdk.OptionDoc {
	return sdk.OptionDoc{
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

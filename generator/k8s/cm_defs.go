package k8s

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/naivary/codemark/definition"
	"github.com/naivary/codemark/definition/target"
	"github.com/naivary/codemark/maker"
	corev1 "k8s.io/api/core/v1"
)

func cmIdent(option string) string {
	return fmt.Sprintf("k8s:configmap:%s", option)
}

type Name string

func (n Name) apply(ident string, cm *corev1.ConfigMap) error {
	cm.Name = string(n)
	return nil
}

type Default string

func (d Default) apply(ident string, cm *corev1.ConfigMap) error {
	lower := strings.ToLower(ident)
	cm.Data[lower] = string(d)
	return nil
}

type Immutable bool

func (i Immutable) apply(ident string, cm *corev1.ConfigMap) error {
	b := bool(i)
	cm.Immutable = &b
	return nil
}

func configMapDefs() []*definition.Definition {
	return []*definition.Definition{
		// structs
		maker.MustMakeDefWithDoc(cmIdent("name"), reflect.TypeFor[Name](), "name of the configmap", target.STRUCT),
		maker.MustMakeDefWithDoc(cmIdent("immutable"), reflect.TypeFor[Immutable](), "decides wether the configmap should be immutable", target.STRUCT),

		// fields
		maker.MustMakeDefWithDoc(cmIdent("default"), reflect.TypeFor[Default](), "default value for the config map value", target.FIELD),
	}
}

// ...
// data:
//   cpu: "3"
//   memory: "1024"

type Config struct {
	// +k8s:configmap:default="3"
	CPU int

	// +k8s:configmap:default="1024"
	Memory int
}

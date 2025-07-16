package k8s

import (
	"github.com/naivary/codemark/api"
	"github.com/naivary/codemark/definition/target"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Name string

func (n Name) apply(m *metav1.ObjectMeta) error {
	m.Name = string(n)
	return nil
}

func (n Name) Doc() api.OptionDoc {
	return api.OptionDoc{
		Targets: []target.Target{target.ANY},
		Doc:     "metadata.name",
		Default: "go identifier of the struct",
	}
}

type Namespace string

// TODO: check if object is namespaced??
func (n Namespace) apply(m *metav1.ObjectMeta) error {
	m.Namespace = string(n)
	return nil
}

func (n Namespace) Doc() api.OptionDoc {
	return api.OptionDoc{
		Targets: []target.Target{target.ANY},
		Doc:     "metadata.namespace",
		Default: "default",
	}
}

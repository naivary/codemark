package k8s

import (
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func objectMetaIdent(option string) string {
	return fmt.Sprintf("k8s:meta:%s", option)
}

type Name string

func (n Name) apply(m *metav1.ObjectMeta) error {
	m.Name = string(n)
	return nil
}

type Namespace string

// TODO: check if object is namespaced??
func (n Namespace) apply(m *metav1.ObjectMeta) error {
	m.Namespace = string(n)
	return nil
}

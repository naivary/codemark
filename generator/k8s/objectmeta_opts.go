package k8s

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Name string

func (n Name) apply(m *metav1.ObjectMeta) error {
	m.Name = string(n)
	return nil
}

type Namespace string

func (n Namespace) apply(m *metav1.ObjectMeta) error {
	m.Namespace = string(n)
	return nil
}

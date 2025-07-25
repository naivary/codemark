package k8s

import (
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	docv1 "github.com/naivary/codemark/api/doc/v1"
)

type Name string

func (n Name) apply(m *metav1.ObjectMeta) error {
	s := string(n)
	if s == "" {
		return fmt.Errorf("name cannot be empty")
	}
	m.Name = string(n)
	return nil
}

func (n Name) Doc() docv1.Option {
	return docv1.Option{
		Desc:    `Name defines the name in the metadata field of the object`,
		Default: "Identifier name of the expression the Option is used on",
	}
}

type Namespace string

func (n Namespace) apply(m *metav1.ObjectMeta) error {
	s := string(n)
	if s == "" {
		s = "default"
	}
	m.Namespace = s
	return nil
}

func (n Namespace) Doc() docv1.Option {
	return docv1.Option{
		Desc:    `Namespace of the object`,
		Default: `default`,
	}
}

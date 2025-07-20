package k8s

import (
	"fmt"

	"github.com/naivary/codemark/api/core"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

func (n Name) Doc() core.OptionDoc {
	return core.OptionDoc{
		Doc:     `Name defines the name in the metadata field of the object`,
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

func (n Namespace) Doc() core.OptionDoc {
	return core.OptionDoc{
		Doc:     `Namespace of the object`,
		Default: `default`,
	}
}

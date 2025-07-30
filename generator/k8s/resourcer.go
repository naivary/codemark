package k8s

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	genv1 "github.com/naivary/codemark/api/generator/v1"
	infov1 "github.com/naivary/codemark/api/info/v1"
	optionv1 "github.com/naivary/codemark/api/option/v1"
)

type Resourcer[I infov1.Info] interface {
	// Resouce represented by this resouce
	Resource() string
	// Options of the resource
	Options() []*optionv1.Option
	// Create generated the actual artifact
	Create(info I, metadata metav1.ObjectMeta) (*genv1.Artifact, error)
}

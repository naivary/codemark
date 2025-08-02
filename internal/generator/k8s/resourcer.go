package k8s

import (
	"go/types"

	"golang.org/x/tools/go/packages"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	genv1 "github.com/naivary/codemark/api/generator/v1"
	infov1 "github.com/naivary/codemark/api/info/v1"
	optionv1 "github.com/naivary/codemark/api/option/v1"
)

type Resourcer interface {
	// Resouce represented by this resouce
	Resource() string

	// Options of the resource
	Options() []*optionv1.Option

	// Create generated the actual artifact
	Create(pkg *packages.Package, obj types.Object, info infov1.Info, metadata metav1.ObjectMeta, config *config) (*genv1.Artifact, error)

	CanCreate(info infov1.Info) bool
}

package k8s

import (
	loaderapi "github.com/naivary/codemark/api/loader"
	"github.com/naivary/codemark/sdk"
)

var _ sdk.Generator = (*k8sGenerator)(nil)

type k8sGenerator struct{}

func (g k8sGenerator) Generate(proj *loaderapi.Project) ([]byte, error) {
	return nil, nil
}

package k8s

import (
	"testing"

	corev1 "k8s.io/api/core/v1"

	genv1 "github.com/naivary/codemark/api/generator/v1"
	loaderv1 "github.com/naivary/codemark/api/loader/v1"
	"github.com/naivary/codemark/loader"
)

func load(path string) (genv1.Generator, loaderv1.Project) {
	gen, err := New()
	if err != nil {
		panic(err)
	}
	proj, err := loader.Load(gen.Registry(), path)
	if err != nil {
		panic(err)
	}
	return gen, proj
}

func TestResource_ConfigMap(t *testing.T) {
	tests := []struct {
		name string
		path string
		want corev1.ConfigMap
	}{
		{
			name: "valid configmap",
			path: "tests/configmap/valid.go",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gen, proj := load(tc.path)
			artifacts, err := gen.Generate(proj)
			if err != nil {
				t.Errorf("err occured: %s\n", err)
			}
			t.Log(artifacts[0].Data)
		})
	}
}

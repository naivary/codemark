package k8s

import (
	"testing"
)

func TestResource_RBAC(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		isValid bool
	}{
		{
			name:    "valid",
			path:    "./tests/rbac/valid.go",
			isValid: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gen, proj := load(tc.path)
			_, err := gen.Generate(proj)
			if err != nil {
				t.Errorf("err occured: %s\n", err)
			}
		})
	}
}

package codemark

import (
	"path/filepath"
	"testing"
)

func TestConfigDocResourcer(t *testing.T) {
	tests := []struct {
		path    string
		cfgFile string
		isValid bool
	}{
		{
			path:    "testdata/config_simple.go",
			isValid: true,
		},
		{
			path:    "testdata/config_with_map.go",
			isValid: true,
		},
	}
	for _, tc := range tests {
		name := filepath.Base(tc.path)
		t.Run(name, func(t *testing.T) {
			arts, err := gen(tc.path, tc.cfgFile)
			if err != nil && tc.isValid {
				t.Errorf("expected no error: %s\n", err)
			}
			t.Log(arts)
		})
	}
}

package openapi

import "testing"

func TestResourcer_Schema(t *testing.T) {
	tests := []struct {
		path    string
		name    string
		isValid bool
	}{
		{
			path:    "testdata/schema/valid.go",
			name:    "valid.go",
			isValid: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := gen(tc.path)
			if err != nil && tc.isValid {
				t.Errorf("unexpected err occured: %s", err)
			}
		})
	}
}

package main

import (
	"testing"
)

func TestLoader(t *testing.T) {
	tests := []struct {
		name  string
		paths []string
	}{
		{
			name:  "simple file",
			paths: []string{"testdata/auth_req.go"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			l := NewLoader(nil)
			_, err := l.Load(tc.paths...)
			if err != nil {
				t.Fatal(err.Error())
			}
		})
	}
}

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
			paths: []string{"testdata/auth_req.go", "testdata/consts.go"},
		},
		{
			name:  "sub directory",
			paths: []string{"testdata/pkg1/consts.go"},
		},
	}

	for _, tc := range tests {
		reg := NewRegistry()
		mapper, err := NewMapper(reg)
        if err != nil {
            t.Error(err)
        }
		t.Run(tc.name, func(t *testing.T) {
			l := NewLoader(mapper, nil)
			info, err := l.Load(tc.paths...)
			if err != nil {
				t.Fatal(err.Error())
			}
			t.Log(info)
		})
	}
}

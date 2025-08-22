package openapi

import (
	"encoding/json"
	"path/filepath"
	"slices"
	"testing"
)

func TestResourcer_Schema(t *testing.T) {
	tests := []struct {
		path    string
		isValid bool
		want    Schema
	}{
		{
			path:    "testdata/schema/required.go",
			isValid: true,
			want: Schema{
				ID:    "auth_request.json",
				Draft: "https://json-schema.org/draft/2020-12/schema",
				Title: "authentication request",
				Desc:  "authentication request data type",
				Type:  objectType,
				Properties: map[string]*Schema{
					"email": {
						Type: stringType,
					},
					"password": {
						Type: stringType,
					},
				},
				Required: []string{"email", "password"},
			},
		},
		{
			path:    "testdata/schema/named_ref.go",
			isValid: true,
			want: Schema{
				ID:    "struct.json",
				Draft: "https://json-schema.org/draft/2020-12/schema",
				Desc:  "example",
				Type:  objectType,
				Properties: map[string]*Schema{
					"f1": {
						Type: stringType,
					},
					"f2": {
						Type: integerType,
					},
				},
			},
		},
		{
			path:    "testdata/schema/enum.go",
			isValid: true,
			want: Schema{
				ID:    "enum.json",
				Draft: "https://json-schema.org/draft/2020-12/schema",
				Title: "enum-test",
				Type:  objectType,
				Properties: map[string]*Schema{
					"f1": {
						Type: stringType,
						Enum: []any{"e1", "e2"},
					},
					"f2": {
						Type: integerType,
						Enum: []any{1, 2},
					},
					"f3": {
						Type: numberType,
						Enum: []any{1.1, 2.2},
					},
					"f4": {
						Type: arrayType,
						Items: &Schema{
							Type: stringType,
							Enum: []any{"e1", "e2"},
						},
					},
					"f5": {
						Type: arrayType,
						Items: &Schema{
							Type: integerType,
							Enum: []any{1, 2},
						},
					},
					"f6": {
						Type: arrayType,
						Items: &Schema{
							Type: numberType,
							Enum: []any{1.1, 2.2},
						},
					},
				},
			},
		},
		{
			path:    "testdata/schema/examples_invalid.go",
			isValid: false,
			want:    _schemaz,
		},
	}
	for _, tc := range tests {
		name := filepath.Base(tc.path)
		t.Run(name, func(t *testing.T) {
			artifacts, err := gen(tc.path)
			if err != nil && tc.isValid {
				t.Errorf("unexpected err occured: %s", err)
			}
			if err != nil && !tc.isValid {
				t.Logf("skipping rest of the test because expected error occured: %s", err)
				t.SkipNow()
			}
			artifact := artifacts[0]
			got := Schema{}
			err = json.NewDecoder(artifact.Data).Decode(&got)
			if err != nil {
				t.Errorf("unexpected err occured: %s", err)
			}
			wantJSON := mustMarshal(tc.want)
			gotJSON := mustMarshal(got)
			if !slices.Equal(wantJSON, gotJSON) {
				t.Errorf("schemas are not equal.\n want: %s\n got: %s", wantJSON, gotJSON)
				t.FailNow()
			}
			t.Logf("Success!\n want: %s\n got: %s", wantJSON, gotJSON)
		})
	}
}

package parser

import (
	"reflect"
	"testing"

	"github.com/naivary/codemark/parser/marker"
	sdktesting "github.com/naivary/codemark/sdk/testing"
)

func valueFor[T any](v T) reflect.Value {
	return reflect.ValueOf(v)
}

type want struct {
	kind  marker.Kind
	value reflect.Value
}

// TODO: better testing of parser..
func TestParse(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		isValid bool
		want    map[string]want
	}{
		{
			name: "multi marker all types",
			input: `+codemark:parser:string="codemark"
			+codemark:parser:int=3`,
			isValid: true,
			want: map[string]want{
				"codemark:parser:string": {kind: marker.STRING, value: valueFor("codemark")},
				"codemark:parser:int":    {kind: marker.INT, value: valueFor[int64](3)},
			},
		},
		// 		{
		// 			name:    "catch error",
		// 			input:   "+jsonschema:validation:max=",
		// 			isValid: false,
		// 		},
		// 		{
		// 			name: "multiline string",
		// 			input: `+codemark:parser:string.multiline=` + "`" + `multi line
		// string` + "`",
		// 			isValid: true,
		// 		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			markers, err := Parse(tc.input)
			if err != nil && tc.isValid {
				t.Errorf("expected to be valid but got an error marker: %v", err)
			}
			if err == nil && !tc.isValid {
				t.Errorf("expected error but err was nil. got: %v\n", markers)
			}
			for _, m := range markers {
				want := tc.want[m.Ident]
				if m.Kind != want.kind {
					t.Errorf("marker kind not equal. got: %s; want: %s\n", m.Kind, want.kind)
				}
				vvfn := sdktesting.GetVVFn(m.Value.Type())
				if !vvfn(m.Value, want.value) {
					t.Errorf("marker value not equal. got: %v; want: %s\n", m.Value, want.value)
				}
				t.Logf("SUCCESS. kind: `%v`; value: `%v`\n", m.Kind, m.Value)
			}
		})
	}
}

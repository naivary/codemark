package parser

import (
	"reflect"
	"testing"

	"github.com/naivary/codemark/marker"
)

const multiLineString = `this is a multi line 
string idk a
something
`

func valueFor[T any](v T) want {
	value := reflect.ValueOf(v)
	mkind := marker.KindFromRType(value.Type())
	return want{kind: mkind, value: value}
}

type want struct {
	kind  marker.Kind
	value reflect.Value
}

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
					+codemark:parser:int=3
					+codemark:parser:float=3.3
					+codemark:parser:complex=3+3i
					+codemark:parser:bool=false
					`,
			isValid: true,
			want: map[string]want{
				"codemark:parser:string":  valueFor("codemark"),
				"codemark:parser:int":     valueFor[int64](3),
				"codemark:parser:float":   valueFor(3.3),
				"codemark:parser:complex": valueFor(3 + 3i),
				"codemark:parser:bool":    valueFor(false),
			},
		},
		{
			name:    "catch error",
			input:   "+jsonschema:validation:max=",
			isValid: false,
		},
		{
			name:    "multiline string",
			input:   `+codemark:parser:string.multiline=` + "`" + multiLineString + "`",
			isValid: true,
			want: map[string]want{
				"codemark:parser:string.multiline": valueFor(multiLineString),
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			markers, err := Parse(tc.input)
			t.Log(markers)
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
				if !m.IsEqual(want.value) {
					t.Errorf("marker value not equal. got: %v; want: %v\n", m.Value, want.value)
				}
				t.Logf("SUCCESS. kind: `%v`; value: `%v`\n", m.Kind, m.Value)
			}
		})
	}
}

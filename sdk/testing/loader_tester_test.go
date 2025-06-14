package testing

import "testing"

func TestRandField(t *testing.T) {
	s := randStruct()
	t.Log(s.Name)
	for _, m := range s.Markers {
		t.Log(m)
	}
	for _, f := range s.Fields {
		t.Log(f.F.Name)
	}
}

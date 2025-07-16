package generator

import "testing"

func TestManager_Load(t *testing.T) {
	mngr, err := NewManager()
	if err != nil {
		t.Errorf("err occured: %s\n", err)
	}
	err = mngr.Load("k8s.so")
	if err != nil {
		t.Errorf("err occured: %s\n", err)
	}
}

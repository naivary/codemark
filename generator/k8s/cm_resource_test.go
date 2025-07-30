package k8s

import (
	"testing"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/goccy/go-yaml"
)

func TestResource_ConfigMap(t *testing.T) {
	immutable := true
	mutable := false
	_ = immutable
	tests := []struct {
		name    string
		path    string
		isValid bool
		want    corev1.ConfigMap
	}{
		{
			name:    "configmap with defaults",
			path:    "tests/configmap/defaults.go",
			isValid: true,
			want: corev1.ConfigMap{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "v1",
					Kind:       "ConfigMap",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "valid",
					Namespace: "default",
				},
				Immutable: &mutable,
				Data: map[string]string{
					"int":    "4",
					"string": "1024",
				},
			},
		},
		{
			name:    "immutable configmap without default",
			path:    "tests/configmap/immutable.go",
			isValid: false,
		},
		{
			name:    "configmap with format",
			path:    "tests/configmap/format.go",
			isValid: true,
			want: corev1.ConfigMap{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "v1",
					Kind:       "ConfigMap",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "format",
					Namespace: "codemark",
				},
				Immutable: &mutable,
				Data: map[string]string{
					"int":    "4",
					"string": "1024",
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			artifacts, err := gen(tc.path)
			if err != nil && tc.isValid {
				t.Errorf("err occured: %s\n", err)
			}
			if err != nil && !tc.isValid {
				t.Skipf("expected error occucred: %s", err)
			}
			if len(artifacts) == 0 {
				t.Errorf("no artifacts generated")
			}
			got := corev1.ConfigMap{}
			err = yaml.NewDecoder(artifacts[0].Data).Decode(&got)
			if err != nil {
				t.Errorf("err occured: %s\n", err)
			}
			if got.Name != tc.want.Name {
				t.Errorf("name not equal. got: %s; want: %s\n", got.Name, tc.want.Name)
			}
			if tc.want.Immutable == nil && got.Immutable != nil || tc.want.Immutable != nil && got.Immutable == nil {
				t.Errorf("immutable check failed. got: %v; want: %v", got.Immutable, tc.want.Immutable)
			}
			for key, wantValue := range tc.want.Data {
				gotValue, found := got.Data[key]
				if !found {
					t.Errorf("missing key: %s\n", key)
				}
				if wantValue != gotValue {
					t.Errorf("values not equal. got: %s; want: %s\n", gotValue, wantValue)
				}
				if wantValue != "" {
					t.Logf("%s=%s\n", key, wantValue)
					continue
				}
				t.Logf("%s=%s\n", key, "<empty-string>")
			}
		})
	}
}

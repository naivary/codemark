package k8s

import (
	"testing"

	"github.com/goccy/go-yaml"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestResource_RBAC(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		isValid bool
		role    rbacv1.Role
		binding rbacv1.RoleBinding
		sva     corev1.ServiceAccount
	}{
		{
			name:    "valid",
			path:    "./tests/rbac/valid.go",
			isValid: true,
			role: rbacv1.Role{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "my-application-name",
					Namespace: "default",
				},
				TypeMeta: metav1.TypeMeta{
					APIVersion: "rbac.authorization.k8s.io/v1",
					Kind:       "Role",
				},
				Rules: []rbacv1.PolicyRule{
					{
						APIGroups: []string{""},
						Verbs:     []string{"get", "list"},
						Resources: []string{"pods"},
					},
				},
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			artifacts, err := gen(tc.path)
			if err != nil {
				t.Errorf("err occured: %s\n", err)
			}
			role := rbacv1.Role{}
			binding := rbacv1.RoleBinding{}
			sva := corev1.ServiceAccount{}
			dec := yaml.NewDecoder(artifacts[0].Data)

			dec.Decode(&role)
			dec.Decode(&binding)
			dec.Decode(&sva)
			t.Log(role.Rules, binding, sva)
		})
	}
}

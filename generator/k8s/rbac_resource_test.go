package k8s

import (
	"fmt"
	"slices"
	"testing"

	"github.com/goccy/go-yaml"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	genv1 "github.com/naivary/codemark/api/generator/v1"
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
			path:    "testdata/rbac/valid.go",
			isValid: true,
			role: rbacv1.Role{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "codemark-rbac",
					Namespace: "codemark-rbac",
				},
				Rules: []rbacv1.PolicyRule{
					{
						APIGroups: []string{""},
						Verbs:     []string{"get", "list"},
						Resources: []string{"pods"},
					},
				},
			},
			binding: rbacv1.RoleBinding{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "codemark-rbac",
					Namespace: "codemark-rbac",
				},
				RoleRef: rbacv1.RoleRef{
					APIGroup: "rbac.authorization.k8s.io",
					Kind:     "Role",
					Name:     "codemark-rbac",
				},
				Subjects: []rbacv1.Subject{
					{
						Kind:      "ServiceAccount",
						Name:      "codemark-rbac",
						Namespace: "codemark-rbac",
					},
				},
			},
			sva: corev1.ServiceAccount{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "codemark-rbac",
					Namespace: "codemark-rbac",
				},
			},
		},
		{
			name:    "partially defined",
			path:    "testdata/rbac/partial.go",
			isValid: false,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			artifacts, err := gen(tc.path)
			if err != nil && tc.isValid {
				t.Errorf("err occured: %s\n", err)
				t.FailNow()
			}
			if err != nil && !tc.isValid {
				t.Skipf("expected error occucred: %s", err)
			}
			if len(artifacts) == 0 {
				t.Errorf("no artifacts generated")
				t.FailNow()
			}
			role, binding, sva := readRBACFromAritfact(artifacts[0])
			err = compareRole(role, tc.role)
			if err != nil {
				t.Errorf("role error: %s", err)
			}
			err = compareRoleBinding(binding, tc.binding)
			if err != nil {
				t.Errorf("role binding error: %s", err)
			}
			err = compareServiceAccount(sva, tc.sva)
			if err != nil {
				t.Errorf("service account error: %s", err)
			}
		})
	}
}

func readRBACFromAritfact(artifact *genv1.Artifact) (rbacv1.Role, rbacv1.RoleBinding, corev1.ServiceAccount) {
	dec := yaml.NewDecoder(artifact.Data)
	role := rbacv1.Role{}
	binding := rbacv1.RoleBinding{}
	sva := corev1.ServiceAccount{}
	err := dec.Decode(&role)
	if err != nil {
		panic(err)
	}
	err = dec.Decode(&binding)
	if err != nil {
		panic(err)
	}
	err = dec.Decode(&sva)
	if err != nil {
		panic(err)
	}
	return role, binding, sva
}

func compareServiceAccount(got, want corev1.ServiceAccount) error {
	if got.Name != want.Name {
		return fmt.Errorf("service account name not equal: got: %s; want: %s", got.Name, want.Name)
	}
	if got.Namespace != want.Namespace {
		return fmt.Errorf("service account namespace not equal: got: %s; want: %s", got.Namespace, want.Namespace)
	}
	return nil
}

func compareRoleBinding(got, want rbacv1.RoleBinding) error {
	if got.Name != want.Name {
		return fmt.Errorf("role binding name not equal: got: %s; want: %s", got.Name, want.Name)
	}
	if got.Namespace != want.Namespace {
		return fmt.Errorf("role binding namespace not equal: got: %s; want: %s", got.Namespace, want.Namespace)
	}
	if got.RoleRef != want.RoleRef {
		return fmt.Errorf("roleRef not equal. got: %v; want: %v", got.RoleRef, want.RoleRef)
	}
	for i, subject := range got.Subjects {
		if subject != want.Subjects[i] {
			return fmt.Errorf("subject not equal. got: %v; want: %v", subject, want.Subjects[i])
		}
	}
	return nil
}

func compareRole(got, want rbacv1.Role) error {
	if got.Name != want.Name {
		return fmt.Errorf("role name not equal: got: %s; want: %s", got.Name, want.Name)
	}
	if got.Namespace != want.Namespace {
		return fmt.Errorf("role namespace not equal: got: %s; want: %s", got.Namespace, want.Namespace)
	}
	for i, rule := range got.Rules {
		if !slices.Equal(rule.APIGroups, want.Rules[i].APIGroups) {
			return fmt.Errorf("APIGroups not equal. got: %v; want: %v", rule.APIGroups, want.Rules[i].APIGroups)
		}
		if !slices.Equal(rule.Resources, want.Rules[i].Resources) {
			return fmt.Errorf("Resources not equal. got: %v; want: %v", rule.Resources, want.Rules[i].Resources)
		}
		if !slices.Equal(rule.Verbs, want.Rules[i].Verbs) {
			return fmt.Errorf("Verbs not equal. got: %v; want: %v", rule.Verbs, want.Rules[i].Verbs)
		}
	}
	return nil
}

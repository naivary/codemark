package k8s

import (
	"strings"

	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	loaderapi "github.com/naivary/codemark/api/loader"
)

func newRBACRole() *rbacv1.Role {
	role := &rbacv1.Role{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "rbac.authorization.k8s.io/v1",
			Kind:       "Role",
		},
		Rules: make([]rbacv1.PolicyRule, 1),
	}
	return role
}

func isRBACResource(ident string) bool {
	return strings.Split(ident, ":")[1] == "rbac"
}

func createRBACRole(fn loaderapi.FuncInfo) (*rbacv1.Role, error) {
	role := newRBACRole()
	objectMeta, err := createObjectMeta(fn)
	if err != nil {
		return nil, err
	}
	role.ObjectMeta = *objectMeta
	for ident, opts := range fn.Opts {
		if !isRBACResource(ident) {
			continue
		}
		for i, opt := range opts {
			if len(role.Rules) <= i {
				role.Rules = append(role.Rules, rbacv1.PolicyRule{})
			}
			rule := role.Rules[i]
			var err error
			switch o := opt.(type) {
			case APIGroups:
				err = o.apply(&rule)
			case Resources:
				err = o.apply(&rule)
			case Verbs:
				err = o.apply(&rule)
			}
			if err != nil {
				return nil, err
			}
			role.Rules[i] = rule
		}
	}
	return role, nil
}

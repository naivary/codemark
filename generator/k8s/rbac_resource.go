package k8s

import (
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	loaderapi "github.com/naivary/codemark/api/loader"
)

func newRBACRole(fn loaderapi.FuncInfo) (rbacv1.Role, error) {
	role := rbacv1.Role{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "rbac.authorization.k8s.io/v1",
			Kind:       "Role",
		},
		Rules: make([]rbacv1.PolicyRule, 1),
	}
	objectMeta, err := createObjectMeta(fn)
	if err != nil {
		return role, err
	}
	role.ObjectMeta = objectMeta
	return role, nil
}

func createRBACRole(fn loaderapi.FuncInfo) (rbacv1.Role, error) {
	role, err := newRBACRole(fn)
	if err != nil {
		return role, err
	}
	for ident, opts := range fn.Opts {
		if !isResource(ident, _rbacResource) {
			continue
		}
		for i, opt := range opts {
			if len(role.Rules) <= i {
				role.Rules = append(role.Rules, rbacv1.PolicyRule{})
			}
			rule := role.Rules[i]
			err = applyOptToRBACRule(opt, &rule)
			if err != nil {
				return role, err
			}
			role.Rules[i] = rule
		}
	}
	return role, nil
}

func applyOptToRBACRule(opt any, rule *rbacv1.PolicyRule) error {
	var err error
	switch o := opt.(type) {
	case APIGroups:
		err = o.apply(rule)
	case Resources:
		err = o.apply(rule)
	case Verbs:
		err = o.apply(rule)
	}
	return err
}

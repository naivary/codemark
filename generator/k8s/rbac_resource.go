package k8s

import (
	"errors"
	"fmt"

	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	genv1 "github.com/naivary/codemark/api/generator/v1"
	infov1 "github.com/naivary/codemark/api/info/v1"
)

var errInvalidRole = errors.New(`
one of the rules in the role is incomplete. This is usually a sign that you want to 
define multiple rules for the rbac role but have missed one of the markers:
// BAD 
// +k8s:rbac:apigroups=[""]
// +k8s:rbac:verbs=["get"]
// +k8s:rbac:resources=["pod"]
// +k8s:rbac:resources=["pod", "service"] <-- second rule in RBAC role partially defined
----
// GOOD
// +k8s:rbac:apigroups=[""]
// +k8s:rbac:verbs=["get"]
// +k8s:rbac:resources=["pod"]
// +k8s:rbac:apigroups=[""] <-- added apigroups
// +k8s:rbac:verbs=["get", "list"] <-- added verbs
// +k8s:rbac:resources=["pod", "service"]`)

func newRBACRole(fn infov1.FuncInfo) (rbacv1.Role, error) {
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

func newRBACRoleBinding(fn infov1.FuncInfo) (rbacv1.RoleBinding, error) {
	roleBinding := rbacv1.RoleBinding{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "rbac.authorization.k8s.io/v1",
			Kind:       "RoleBinding",
		},
	}
	objectMeta, err := createObjectMeta(fn)
	if err != nil {
		return roleBinding, err
	}
	roleBinding.ObjectMeta = objectMeta
	return roleBinding, nil
}

func validateRole(role rbacv1.Role) error {
	for _, rule := range role.Rules {
		if rule.APIGroups == nil || rule.Resources == nil || rule.Verbs == nil {
			return errInvalidRole
		}
	}
	return nil
}

// createRBAC creates the set of RBAC manifests e.g. Role, RoleBinding and
// ServiceAccount.
func createRBAC(fn infov1.FuncInfo) (*genv1.Artifact, error) {
	role, err := createRBACRole(fn)
	if err != nil {
		return nil, err
	}
	binding, err := createRBACRoleBinding(fn)
	if err != nil {
		return nil, err
	}
	filename := fmt.Sprintf("%s.rbac.yaml", role.Name)
	rbacArtifact, err := newArtifact(filename, role, binding)
	if err != nil {
		return nil, err
	}
	const svaName = "k8s:serviceaccount:name"
	name := fn.Opts[svaName]
	if len(name) == 0 {
		fn.Opts[svaName] = []any{ServiceAccountName(role.Name)}
	}
	sva, err := createServiceAccount(fn)
	if err != nil {
		return nil, err
	}
	return mergeArtifacts(rbacArtifact, sva)
}

func createRBACRoleBinding(fn infov1.FuncInfo) (rbacv1.RoleBinding, error) {
	binding, err := newRBACRoleBinding(fn)
	if err != nil {
		return binding, err
	}
	return binding, nil
}

func createRBACRole(fn infov1.FuncInfo) (rbacv1.Role, error) {
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
	return role, validateRole(role)
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

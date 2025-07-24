package k8s

import (
	rbacv1 "k8s.io/api/rbac/v1"

	loaderapi "github.com/naivary/codemark/api/loader"
)

func newRBACRole() *rbacv1.Role {
	role := &rbacv1.Role{}
	return role
}

func createRBACRole(fn loaderapi.FuncInfo) (*rbacv1.Role, error) {
	role := newRBACRole()
	objectMeta, err := createObjectMeta(fn)
	if err != nil {
		return nil, err
	}
	role.ObjectMeta = *objectMeta
	for _, opts := range fn.Opts {
		for _, opt := range opts {
			var err error
			switch o := opt.(type) {
			case APIGroups:
				err = o.apply(&role.Rules[0])
			}
			if err != nil {
				return nil, err
			}
		}
	}
	return role, nil
}

package k8s

import (
	rbacv1 "k8s.io/api/rbac/v1"

	docv1 "github.com/naivary/codemark/api/doc/v1"
)

type APIGroups []string

func (a APIGroups) apply(r *rbacv1.PolicyRule) error {
	r.APIGroups = a
	return nil
}

func (a APIGroups) Doc() docv1.Option {
	return docv1.Option{
		Desc:    `API Groups of the Role`,
		Default: "REQUIRED",
	}
}

type Resources []string

func (r Resources) apply(ru *rbacv1.PolicyRule) error {
	ru.Resources = r
	return nil
}

func (r Resources) Doc() docv1.Option {
	return docv1.Option{
		Desc:    `Resources`,
		Default: "REQUIRED",
	}
}

type Verbs []string

func (v Verbs) apply(ru *rbacv1.PolicyRule) error {
	ru.Verbs = v
	return nil
}

func (v Verbs) Doc() docv1.Option {
	return docv1.Option{
		Desc:    `Verbs`,
		Default: "REQUIRED",
	}
}

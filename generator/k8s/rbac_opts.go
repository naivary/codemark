package k8s

import (
	rbacv1 "k8s.io/api/rbac/v1"

	"github.com/naivary/codemark/api/doc"
)

type APIGroups []string

func (a APIGroups) apply(r *rbacv1.PolicyRule) error {
	r.APIGroups = a
	return nil
}

func (a APIGroups) Doc() doc.Option {
	return doc.Option{
		Desc:    `API Groups of the Role`,
		Default: "REQUIRED",
	}
}

type Resources []string

type Verbs []string

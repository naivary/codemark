package k8s

import (
	rbacv1 "k8s.io/api/rbac/v1"

	"github.com/naivary/codemark/api/doc"
	optionapi "github.com/naivary/codemark/api/option"
)

const _rbacResource = "rbac"

func rbacOpts() []*optionapi.Option {
	return makeOpts(_rbacResource,
		mustMakeOpt(_typeName, APIGroups(nil), false, optionapi.TargetFunc),
		mustMakeOpt(_typeName, Resources(nil), false, optionapi.TargetFunc),
		mustMakeOpt(_typeName, Verbs(nil), false, optionapi.TargetFunc),
	)
}

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

func (r Resources) apply(ru *rbacv1.PolicyRule) error {
	ru.Resources = r
	return nil
}

func (r Resources) Doc() doc.Option {
	return doc.Option{
		Desc:    `Resources`,
		Default: "REQUIRED",
	}
}

type Verbs []string

func (v Verbs) apply(ru *rbacv1.PolicyRule) error {
	ru.Verbs = v
	return nil
}

func (v Verbs) Doc() doc.Option {
	return doc.Option{
		Desc:    `Verbs`,
		Default: "REQUIRED",
	}
}

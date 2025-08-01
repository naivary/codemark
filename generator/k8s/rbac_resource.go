package k8s

import (
	"errors"
	"fmt"
	"go/types"

	"golang.org/x/tools/go/packages"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	genv1 "github.com/naivary/codemark/api/generator/v1"
	infov1 "github.com/naivary/codemark/api/info/v1"
	optionv1 "github.com/naivary/codemark/api/option/v1"
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

const _rbacResource = "rbac"

var _ Resourcer = (*rbacResourcer)(nil)

type rbacResourcer struct {
	resource string
}

func NewRBACResourcer() Resourcer {
	return &rbacResourcer{resource: _rbacResource}
}

func (r *rbacResourcer) Resource() string {
	return r.resource
}

func (r *rbacResourcer) Options() []*optionv1.Option {
	return makeOpts(_rbacResource,
		mustMakeOpt(_typeName, APIGroups(nil), _repetable, optionv1.TargetFunc),
		mustMakeOpt(_typeName, Resources(nil), _repetable, optionv1.TargetFunc),
		mustMakeOpt(_typeName, Verbs(nil), _repetable, optionv1.TargetFunc),
	)
}

func (r *rbacResourcer) CanCreate(info infov1.Info) bool {
	fn, isFunc := info.(*infov1.FuncInfo)
	if !isFunc {
		return false
	}
	if fn.Decl.Name.Name != "main" {
		return false
	}
	return true
}

func (r *rbacResourcer) Create(
	pkg *packages.Package,
	obj types.Object,
	info infov1.Info,
	metadata metav1.ObjectMeta,
	cfg *config,
) (*genv1.Artifact, error) {
	fn := info.(*infov1.FuncInfo)
	role, err := r.newRole(metadata, fn)
	if err != nil {
		return nil, err
	}
	sva, err := r.newServiceAccount(metadata)
	if err != nil {
		return nil, err
	}
	binding, err := r.newRoleBinding(metadata, role, sva)
	if err != nil {
		return nil, err
	}
	filename := fmt.Sprintf("%s.rbac.yaml", role.Name)
	return newArtifact(filename, role, binding, sva)
}

func (r *rbacResourcer) newRole(metadata metav1.ObjectMeta, fn *infov1.FuncInfo) (rbacv1.Role, error) {
	role := rbacv1.Role{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "rbac.authorization.k8s.io/v1",
			Kind:       "Role",
		},
		ObjectMeta: metadata,
		Rules:      make([]rbacv1.PolicyRule, 1),
	}
	if role.Name == "" {
		return role, errors.New("you have to define a +k8s:metadata:name")
	}
	err := r.applyOptsToRole(fn, &role)
	if err != nil {
		return role, err
	}
	return role, r.validateRole(role)
}

func (r *rbacResourcer) applyOptsToRole(fn *infov1.FuncInfo, role *rbacv1.Role) error {
	for ident, opts := range fn.Opts {
		if !isResource(ident, _rbacResource) {
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
				return err
			}
			role.Rules[i] = rule
		}
	}
	return nil
}

func (r *rbacResourcer) newRoleBinding(metadata metav1.ObjectMeta, role rbacv1.Role, sva corev1.ServiceAccount) (rbacv1.RoleBinding, error) {
	roleBinding := rbacv1.RoleBinding{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "rbac.authorization.k8s.io/v1",
			Kind:       "RoleBinding",
		},
		ObjectMeta: metadata,
		RoleRef: rbacv1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "Role",
			Name:     role.Name,
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      sva.Kind,
				Name:      sva.Name,
				Namespace: sva.Namespace,
			},
		},
	}
	return roleBinding, nil
}

func (r *rbacResourcer) newServiceAccount(metadata metav1.ObjectMeta) (corev1.ServiceAccount, error) {
	sva := corev1.ServiceAccount{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "ServiceAccount",
		},
		ObjectMeta: metadata,
	}
	return sva, nil
}

func (r *rbacResourcer) validateRole(role rbacv1.Role) error {
	for _, rule := range role.Rules {
		if rule.APIGroups == nil || rule.Resources == nil || rule.Verbs == nil {
			return errInvalidRole
		}
	}
	return nil
}

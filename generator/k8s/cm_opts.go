package k8s

import (
	"errors"

	corev1 "k8s.io/api/core/v1"

	"github.com/iancoleman/strcase"

	docv1 "github.com/naivary/codemark/api/doc/v1"
	loaderv1 "github.com/naivary/codemark/api/loader/v1"
	optionv1 "github.com/naivary/codemark/api/option/v1"
)

var errImmutableConfigMapWithoutDefault = errors.New(`
when you decide to have an immutable ConfigMap the default value of the marker cannot be empty. 
If you would like to have the empty default value remove the marker completly and use the default 
value of the struct field in go.
`)

const _configMapResource = "configmap"

func configMapOpts() []*optionv1.Option {
	return makeOpts(_configMapResource,
		mustMakeOpt(_typeName, Immutable(false), _optional, _unique, optionv1.TargetStruct),
		mustMakeOpt(_typeName, Default(""), _optional, _unique, optionv1.TargetField),
		mustMakeOpt(_typeName, KeyFormat(CamelCase), _optional, _unique, optionv1.TargetStruct),
	)
}

type Default string

func (d Default) apply(info loaderv1.FieldInfo, cm *corev1.ConfigMap, format KeyFormat) error {
	isImmutable := *cm.Immutable
	if !isImmutable {
		// TODO: check if the type of the field matches the default value
		cm.Data[format.Format(info.Ident.Name)] = string(d)
		return nil
	}
	if len(string(d)) == 0 {
		return errImmutableConfigMapWithoutDefault
	}
	cm.Data[format.Format(info.Ident.Name)] = string(d)
	return nil
}

func (d Default) Doc() docv1.Option {
	return docv1.Option{
		Desc:    `Default defines the default value of the config parameter`,
		Default: "REQUIRED",
	}
}

type Immutable bool

func (i Immutable) apply(cm *corev1.ConfigMap) error {
	b := bool(i)
	cm.Immutable = &b
	return nil
}

func (i Immutable) Doc() docv1.Option {
	return docv1.Option{
		Desc:    `Defines if the ConfigMap is immutable or not`,
		Default: "false",
	}
}

type KeyFormat string

const (
	SnakeCase  KeyFormat = "snake_case"
	CamelCase  KeyFormat = "camelCase"
	PascalCase KeyFormat = "pascalCase"
	KebabCase  KeyFormat = "kebab-case"
	Env        KeyFormat = "env"
)

func (k KeyFormat) Format(key string) string {
	switch k {
	case SnakeCase:
		return strcase.ToSnake(key)
	case CamelCase:
		return strcase.ToLowerCamel(key)
	case PascalCase:
		return strcase.ToCamel(key)
	case Env:
		return strcase.ToScreamingDelimited(key, '_', "", true)
	case KebabCase:
		return strcase.ToKebab(key)
	default:
		return key
	}
}

func (k KeyFormat) Doc() docv1.Option {
	return docv1.Option{
		Desc:    `Format of the key. It be used to manipulate based on the context of the configuration. For example if the configuration is ssettable via environment variable it is useful to choose the env formation.`,
		Default: "lowercase",
	}
}

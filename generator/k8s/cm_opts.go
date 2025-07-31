package k8s

import (
	"errors"

	"github.com/iancoleman/strcase"

	corev1 "k8s.io/api/core/v1"

	docv1 "github.com/naivary/codemark/api/doc/v1"
	infov1 "github.com/naivary/codemark/api/info/v1"
)

var errEmptyDefault = errors.New(`
the default value of a config option cannot be empty. If you want it to be empty just remove the marker 
and use the empty value of the go type itself.
`)

type Default string

func (d Default) apply(info *infov1.FieldInfo, cm *corev1.ConfigMap, format Format) error {
	if len(string(d)) == 0 {
		return errEmptyDefault
	}
	// TODO: check if the type of the field matches the default value
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

type Format string

const (
	SnakeCase  Format = "snake_case"
	CamelCase  Format = "camelCase"
	PascalCase Format = "pascalCase"
	KebabCase  Format = "kebab-case"
	Env        Format = "env"
)

func (k Format) Format(key string) string {
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

func (k Format) Doc() docv1.Option {
	return docv1.Option{
		Desc:    `Format of the key. It be used to manipulate based on the context of the configuration. For example if the configuration is ssettable via environment variable it is useful to choose the env formation.`,
		Default: "camelCase",
	}
}

package k8s

import (
	"go/ast"

	corev1 "k8s.io/api/core/v1"

	"github.com/iancoleman/strcase"

	"github.com/naivary/codemark/api/doc"
	optionapi "github.com/naivary/codemark/api/option"
)

func configMapOpts() []*optionapi.Option {
	const resource = "configmap"
	return makeDefs(resource,
		newOption(Immutable(false), true, optionapi.TargetStruct),
		newOption(Default(""), true, optionapi.TargetField),
		newOption(KeyFormat(""), true, optionapi.TargetStruct),
	)
}

type Default string

func (d Default) apply(field *ast.Ident, cm *corev1.ConfigMap, format KeyFormat) error {
	cm.Data[format.Format(field.Name)] = string(d)
	return nil
}

func (d Default) Doc() doc.Option {
	return doc.Option{
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

func (i Immutable) Doc() doc.Option {
	return doc.Option{
		Desc:    `Defines if the ConfigMap is immutable or not`,
		Default: "false",
	}
}

type KeyFormat string

const (
	SnakeCase  KeyFormat = "snake_case"
	CamelCase  KeyFormat = "camelCalse"
	PascalCase KeyFormat = "pascalCase"
	KebabCase  KeyFormat = "kebabCase"
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

func (k KeyFormat) Doc() doc.Option {
	return doc.Option{
		Desc:    `Format of the key. It be used to manipulate based on the context of the configuration. For example if the configuration is ssettable via environment variable it is useful to choose the env formation.`,
		Default: "lowercase",
	}
}

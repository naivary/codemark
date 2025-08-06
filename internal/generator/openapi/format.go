package openapi

import (
	"github.com/iancoleman/strcase"
)

type NamingConvention string

const (
	SnakeCase  NamingConvention = "snake_case"
	CamelCase  NamingConvention = "camelCase"
	PascalCase NamingConvention = "PascalCase"
	KebabCase  NamingConvention = "kebab-case"
	Env        NamingConvention = "ENV"
)

func (nc NamingConvention) Format(key string) string {
	switch nc {
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
		return ""
	}
}

package openapi

import docv1 "github.com/naivary/codemark/api/doc/v1"

type Pattern string

func (p Pattern) Doc() docv1.Option {
	return docv1.Option{
		Desc:    "Pattern which the given string must fullfill",
		Default: "",
	}
}

type MinLength int

func (m MinLength) Doc() docv1.Option {
	return docv1.Option{
		Desc:    "Minimum length of the string",
		Default: "",
	}
}

type MaxLength int

func (m MaxLength) Doc() docv1.Option {
	return docv1.Option{
		Desc:    "Maximum length of the string",
		Default: "",
	}
}

type ContentEncoding string

func (c ContentEncoding) Doc() docv1.Option {
	return docv1.Option{
		Desc:    "Content encoding",
		Default: "",
	}
}

type ContentMediaType string

func (c ContentMediaType) Doc() docv1.Option {
	return docv1.Option{
		Desc:    "Content media type",
		Default: "",
	}
}

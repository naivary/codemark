package codemark

import docv1 "github.com/naivary/codemark/api/doc/v1"

type OptDocDesc string

func (o OptDocDesc) apply(doc *docv1.Option) {
	doc.Desc = string(o)
}

type OptDocSummary string

func (o OptDocSummary) apply(doc *docv1.Option) {
	doc.Summary = string(o)
}

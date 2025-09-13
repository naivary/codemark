package codemark

import (
	"bytes"
	"fmt"
	"go/types"

	"golang.org/x/tools/go/packages"

	docv1 "github.com/naivary/codemark/api/doc/v1"
	genv1 "github.com/naivary/codemark/api/generator/v1"
	infov1 "github.com/naivary/codemark/api/info/v1"
	optv1 "github.com/naivary/codemark/api/option/v1"
)

const _optDocResource = "option"

type optDocResourcer struct {
	docPkgPath string
}

func NewOptDocResourcer() *optDocResourcer {
	return &optDocResourcer{
		docPkgPath: "github.com/naivary/codemark/api/doc/v1",
	}
}

// Resouce represented by this resource
func (o optDocResourcer) Resource() string {
	return _optDocResource
}

// Options of the resource
func (o optDocResourcer) Options() []*optv1.Option {
	return makeOpts(_optDocResource,
		mustMakeOpt("description", OptDocDesc(""), _unique, optv1.TargetNamed),
		mustMakeOpt("summary", OptDocDesc(""), _unique, optv1.TargetNamed),
	)
}

func (o optDocResourcer) CanCreate(info infov1.Info) bool {
	_, isNamed := info.(*infov1.NamedInfo)
	opts := info.Options()
	return opts.IsDefined("codemark:option:summary") || opts.IsDefined("codemark:option:description") && isNamed
}

// Create generated the actual artifact
func (o optDocResourcer) Create(pkg *packages.Package, info map[types.Object]*infov1.NamedInfo) (*genv1.Artifact, error) {
	docs := make(map[types.Object]*docv1.Option, len(info))
	for obj, named := range info {
		if !o.CanCreate(named) {
			continue
		}
		doc := o.newDoc(named)
		docs[obj] = doc
	}
	return o.createFile(pkg, docs)
}

func (o optDocResourcer) newDoc(named *infov1.NamedInfo) *docv1.Option {
	doc := &docv1.Option{}
	for _, opts := range named.Options().Filter(_domain, _optDocResource) {
		for _, opt := range opts {
			switch t := opt.(type) {
			case OptDocDesc:
				doc.Desc = string(t)
			case OptDocSummary:
				doc.Summary = string(t)
			}
		}
	}
	return doc
}

func (o optDocResourcer) createFile(pkg *packages.Package, docs map[types.Object]*docv1.Option) (*genv1.Artifact, error) {
	var b bytes.Buffer
	fmt.Fprintf(&b, "// CODE GENERATE BY CODEMARK\n")
	fmt.Fprintf(&b, "// DO NOT EDIT!\n")
	fmt.Fprintf(&b, "package %s\n", pkg.Name)
	// import docv1.Option package
	fmt.Fprintf(&b, "import docv1 %s\n", o.docPkgPath)

	for obj, doc := range docs {
		fmt.Fprintf(&b, "func (%s) Doc() *docv1.Option {\n", obj.Name())
		fmt.Fprintf(&b, "\t")
		fmt.Fprintf(&b, "return &docv1.Option{\n")
		fmt.Fprintf(&b, "\t\t")
		fmt.Fprintf(&b, `Desc: "%s"`, doc.Desc)
		fmt.Fprintf(&b, ",\n")
		fmt.Fprintf(&b, "\t\t")
		fmt.Fprintf(&b, `Summary: "%s"`, doc.Desc)
		fmt.Fprintf(&b, ",\n")
		fmt.Fprintf(&b, "\t")
		fmt.Fprintf(&b, "}\n")
		fmt.Fprintf(&b, "}")
		fmt.Fprintf(&b, "\n")
	}
	fmt.Println(b.String())
	return nil, nil
}

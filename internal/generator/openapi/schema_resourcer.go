package openapi

import (
	"fmt"
	"go/types"

	"golang.org/x/tools/go/packages"

	"github.com/iancoleman/strcase"

	genv1 "github.com/naivary/codemark/api/generator/v1"
	infov1 "github.com/naivary/codemark/api/info/v1"
	optionv1 "github.com/naivary/codemark/api/option/v1"
)

const _schemaResource = "schema"

var _ Resourcer = (*schemaResourcer)(nil)

func NewSchemaResourcer() Resourcer {
	return &schemaResourcer{_schemaResource}
}

type schemaResourcer struct {
	resource string
}

func (s schemaResourcer) Resource() string {
	return s.resource
}

func (s schemaResourcer) Options() []*optionv1.Option {
	return makeOpts(s.resource,
		// type agnostic
		mustMakeOpt(_typeName, ID(""), _unique, optionv1.TargetStruct),
		mustMakeOpt(_typeName, Draft(""), _unique, optionv1.TargetStruct),
		mustMakeOpt(_typeName, Description(""), _unique, optionv1.TargetStruct),
		// string
		mustMakeOpt(_typeName, Pattern(""), _unique, optionv1.TargetField),
		mustMakeOpt(_typeName, ContentMediaType(""), _unique, optionv1.TargetField),
		mustMakeOpt(_typeName, ContentEncoding(""), _unique, optionv1.TargetField),
		mustMakeOpt(_typeName, MinLength(0), _unique, optionv1.TargetField),
		mustMakeOpt(_typeName, MaxLength(0), _unique, optionv1.TargetField),
		// numeric
		mustMakeOpt(_typeName, Minimum(0), _unique, optionv1.TargetField),
		mustMakeOpt(_typeName, Maximum(0), _unique, optionv1.TargetField),
		mustMakeOpt(_typeName, ExclusiveMaximum(0), _unique, optionv1.TargetField),
		mustMakeOpt(_typeName, ExclusiveMinimum(0), _unique, optionv1.TargetField),
	)
}

func (s schemaResourcer) CanCreate(info infov1.Info) bool {
	_, isStruct := info.(*infov1.StructInfo)
	return isStruct
}

func (s schemaResourcer) Create(pkg *packages.Package, obj types.Object, info infov1.Info, cfg *config) (*genv1.Artifact, error) {
	structInfo := info.(*infov1.StructInfo)
	s.setDefaultOpts(structInfo, cfg)
	schema := s.newSchema(structInfo)
	err := s.applyStructOpts(&schema, structInfo, cfg)
	if err != nil {
		return nil, err
	}
	for _, finfo := range structInfo.Fields {
		name := cfg.Schema.Formats.Property.Format(finfo.Ident.Name)
		err = s.applyFieldOpts(schema.Properties[name], finfo)
		if err != nil {
			return nil, err
		}
	}
	fmt.Println(schema.ID)
	filename := fmt.Sprintf("%s.json", cfg.Schema.Formats.Filename.Format(structInfo.Spec.Name.Name))
	return newArtifact(filename, schema)
}

func (s schemaResourcer) setDefaultOpts(info *infov1.StructInfo, cfg *config) {
	opts := s.Options()
	setDefaults(opts, info, cfg, optionv1.TargetStruct, map[string]any{
		"openapi:schema:id": ID(info.Spec.Name.Name),
	})
	for _, finfo := range info.Fields {
		setDefaults(opts, finfo, cfg, optionv1.TargetField, nil)
	}
}

func (s schemaResourcer) newSchema(info *infov1.StructInfo) Schema {
	schema := Schema{
		Type:       objectType,
		Properties: make(map[string]*Schema, len(info.Fields)),
	}
	for obj, finfo := range info.Fields {
		fieldName := strcase.ToLowerCamel(finfo.Ident.Name)
		schema.Properties[fieldName] = &Schema{
			Type: jsonTypeOf(obj.Type()),
		}
	}
	return schema
}

func (s schemaResourcer) applyFieldOpts(schema *Schema, info *infov1.FieldInfo) error {
	for ident, opts := range info.Options() {
		if !isResource(ident, _schemaResource) {
			continue
		}
		for _, opt := range opts {
			var err error
			switch o := opt.(type) {
			case Maximum:
				err = o.apply(schema)
			case Minimum:
				err = o.apply(schema)
			}
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (s schemaResourcer) applyStructOpts(schema *Schema, info *infov1.StructInfo, cfg *config) error {
	for ident, opts := range info.Options() {
		if !isResource(ident, _schemaResource) {
			continue
		}
		for _, opt := range opts {
			var err error
			switch o := opt.(type) {
			case Description:
				err = o.apply(schema)
			case Draft:
				err = o.apply(schema)
			case ID:
				err = o.apply(schema, cfg.Schema.IDBaseURL, cfg.Schema.Formats.Filename)
			}
			if err != nil {
				return err
			}
		}
	}
	return nil
}

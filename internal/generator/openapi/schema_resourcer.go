package openapi

import (
	"encoding/json"
	"fmt"
	"go/types"
	"os"
	"path/filepath"

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
		// annotations
		mustMakeOpt(_typeName, Description(""), _unique, optionv1.TargetStruct, optionv1.TargetField),
		mustMakeOpt(_typeName, Title(""), _unique, optionv1.TargetStruct, optionv1.TargetField),
		mustMakeOpt(_typeName, Examples(nil), _unique, optionv1.TargetField),
		mustMakeOpt(_typeName, Default(""), _unique, optionv1.TargetField),
		mustMakeOpt(_typeName, Deprecated(false), _unique, optionv1.TargetStruct, optionv1.TargetField),
		mustMakeOpt(_typeName, WriteOnly(false), _unique, optionv1.TargetField),
		mustMakeOpt(_typeName, ReadOnly(false), _unique, optionv1.TargetField),
		// object
		mustMakeOpt(_typeName, Required(false), _unique, optionv1.TargetField),
		// string
		mustMakeOpt(_typeName, Pattern(""), _unique, optionv1.TargetField),
		mustMakeOpt(_typeName, Format(""), _unique, optionv1.TargetField),
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
	root, err := s.newSchema(structInfo, cfg)
	if err != nil {
		return nil, err
	}
	err = s.applyStructOpts(&root, structInfo)
	if err != nil {
		return nil, err
	}
	for _, finfo := range structInfo.Fields {
		name := cfg.Schema.Formats.Property.Format(finfo.Ident.Name)
		fieldSchema := root.Properties[name]
		err = s.applyFieldOpts(&root, fieldSchema, finfo, cfg)
		if err != nil {
			return nil, err
		}
	}
	fmt.Println(json.NewEncoder(os.Stdout).Encode(&root))
	filename := filepath.Base(root.ID)
	return newArtifact(filename, root)
}

func (s schemaResourcer) setDefaultOpts(info *infov1.StructInfo, cfg *config) {
	opts := s.Options()
	setDefaults(opts, info, cfg, optionv1.TargetStruct, nil)
	for _, finfo := range info.Fields {
		setDefaults(opts, finfo, cfg, optionv1.TargetField, nil)
	}
}

func (s schemaResourcer) newSchema(info *infov1.StructInfo, cfg *config) (Schema, error) {
	id, err := id(info.Spec.Name.Name, cfg.Schema.IDBaseURL, cfg.Schema.Formats.Filename)
	if err != nil {
		return Schema{}, err
	}
	schema := Schema{
		ID:         id,
		Draft:      cfg.Schema.Draft,
		Type:       objectType,
		Properties: make(map[string]*Schema, len(info.Fields)),
	}
	for obj, finfo := range info.Fields {
		fieldName := strcase.ToLowerCamel(finfo.Ident.Name)
		schema.Properties[fieldName] = &Schema{
			Type: jsonTypeOf(obj.Type()),
		}
	}
	return schema, nil
}

func (s schemaResourcer) applyFieldOpts(root, fieldSchema *Schema, info *infov1.FieldInfo, cfg *config) error {
	for ident, opts := range info.Options() {
		if !isResource(ident, _schemaResource) {
			continue
		}
		for _, opt := range opts {
			var err error
			switch o := opt.(type) {
			// annotations
			case Title:
				err = o.apply(fieldSchema)
			case Description:
				err = o.apply(fieldSchema)
			case Examples:
				err = o.apply(fieldSchema)
			case ReadOnly:
				err = o.apply(fieldSchema)
			case WriteOnly:
				err = o.apply(fieldSchema)
			case Default:
				err = o.apply(fieldSchema)
			case Deprecated:
				err = o.apply(fieldSchema)
			// numeric
			case Maximum:
				err = o.apply(fieldSchema)
			case Minimum:
				err = o.apply(fieldSchema)
			case ExclusiveMaximum:
				err = o.apply(fieldSchema)
			case ExclusiveMinimum:
				err = o.apply(fieldSchema)
			case MultipleOf:
				err = o.apply(fieldSchema)
			// string
			case Format:
				err = o.apply(fieldSchema)
			case Pattern:
				err = o.apply(fieldSchema)
			case MaxLength:
				err = o.apply(fieldSchema)
			case MinLength:
				err = o.apply(fieldSchema)
			case ContentEncoding:
				err = o.apply(fieldSchema)
			case ContentMediaType:
				err = o.apply(fieldSchema)
			// object
			case Required:
				err = o.apply(root, info, cfg)
			}
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (s schemaResourcer) applyStructOpts(schema *Schema, info *infov1.StructInfo) error {
	for ident, opts := range info.Options() {
		if !isResource(ident, _schemaResource) {
			continue
		}
		for _, opt := range opts {
			var err error
			switch o := opt.(type) {
			case Description:
				err = o.apply(schema)
			case Deprecated:
				err = o.apply(schema)
			}
			if err != nil {
				return err
			}
		}
	}
	return nil
}

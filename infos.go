package codemark

import (
	"go/ast"
	"go/types"
)

const UnknownMethodName = "UNKNOWN_METHOD_NAME"

type info interface {
	// The raw documentation of the declaration
	Doc() string
	// Converted markers of the declaration
	Defs() Definitions

	Name() string
}

type Info struct {
	Doc   string
	Defs  Definitions
	Type  types.Type
	Obj   types.Object
	Expr  ast.Expr
	Spec  ast.Spec
	Decl  ast.Decl
	Ident *ast.Ident
}

// Converted markers to the given definitions. If multiple equal marker are
// found they will be appended to the slice.
type Definitions map[string][]any

// add adds `def` to the definitions indexed by `idn`
func (d Definitions) add(idn string, def any) {
	defs, ok := d[idn]
	if !ok {
		d[idn] = []any{def}
		return
	}
	d[idn] = append(defs, def)
}

// Get returns the definitions of the given identifier. The second return value
// is indicating if any definitions are found for the given `idn`.
func (d Definitions) Get(idn string) ([]any, bool) {
	defs, ok := d[idn]
	return defs, ok
}

// IsUnique is returning the definition of the given `idn` iff the identifier is
// associated with definition.
func (d Definitions) IsUnique(idn string) (any, bool) {
	defs, ok := d[idn]
	if !ok {
		return nil, false
	}
	if len(defs) == 1 {
		return defs[0], true
	}
	return nil, false
}

func (d Definitions) IsEmpty() bool {
	return len(d) == 0
}

var _ info = (*PackageInfo)(nil)

type PackageInfo struct {
	Info *Info
	File *ast.File
}

func (p *PackageInfo) Doc() string {
	return p.Info.Doc
}

func (p *PackageInfo) Defs() Definitions {
	return p.Info.Defs
}

func (p *PackageInfo) Name() string {
	return p.File.Name.Name
}

var _ info = (*MethodInfo)(nil)

type MethodInfo struct {
	Info *Info
}

func (m MethodInfo) funcDecl() *ast.FuncDecl {
	return m.Info.Decl.(*ast.FuncDecl)
}

func (m MethodInfo) Doc() string {
	return m.Info.Doc
}

func (m MethodInfo) Defs() Definitions {
	return m.Info.Defs
}

func (m MethodInfo) Name() string {
	return m.funcDecl().Name.Name
}

func (m MethodInfo) ReceiverName() string {
	field := m.funcDecl().Recv.List[0]
	if len(field.Names) == 0 {
		return UnknownMethodName
	}
	return field.Names[0].Name
}

func (m MethodInfo) ReceiverType() ast.Expr {
	return m.funcDecl().Recv.List[0].Type
}

func (m MethodInfo) Params() *ast.FieldList {
	return m.funcDecl().Type.Params
}

func (m MethodInfo) Returns() *ast.FieldList {
	return m.funcDecl().Type.Results
}

func (m MethodInfo) Decl() *ast.FuncDecl {
	return m.Info.Decl.(*ast.FuncDecl)
}

var _ info = (*FuncInfo)(nil)

type FuncInfo struct {
	Info *Info
}

func (f FuncInfo) funcDecl() *ast.FuncDecl {
	return f.Info.Decl.(*ast.FuncDecl)
}

func (f FuncInfo) Doc() string {
	return f.Info.Doc
}

func (f FuncInfo) Defs() Definitions {
	return f.Info.Defs
}

func (f FuncInfo) Name() string {
	return f.funcDecl().Name.Name
}

func (f FuncInfo) Params() *ast.FieldList {
	return f.funcDecl().Type.Params
}

func (f FuncInfo) Returns() *ast.FieldList {
	return f.funcDecl().Type.Results
}

var _ info = (*ConstInfo)(nil)

type ConstInfo struct {
	Info *Info
}

func (c ConstInfo) Doc() string {
	return c.Info.Doc
}

func (c ConstInfo) Defs() Definitions {
	return c.Info.Defs
}

func (c ConstInfo) Name() string {
	return c.Info.Ident.Name
}

var _ info = (*VarInfo)(nil)

type VarInfo struct {
	Info *Info
}

func (v VarInfo) Doc() string {
	return v.Info.Doc
}

func (v VarInfo) Defs() Definitions {
	return v.Info.Defs
}

func (v VarInfo) Name() string {
	return v.Info.Ident.Name
}

var _ info = (*StructInfo)(nil)

type StructInfo struct {
	Info   *Info
	Fields []*FieldInfo
}

func (s StructInfo) Doc() string {
	return s.Info.Doc
}

func (s StructInfo) Defs() Definitions {
	return s.Info.Defs
}

func (s StructInfo) Name() string {
	return s.Info.Ident.Name
}

var _ info = (*FieldInfo)(nil)

type FieldInfo struct {
	Info *Info
}

func (f FieldInfo) Doc() string {
	return f.Info.Doc
}

func (f FieldInfo) Defs() Definitions {
	return f.Info.Defs
}

func (f FieldInfo) IsEmbedded() bool {
	return f.Info.Ident == nil
}

func (f FieldInfo) Name() string {
	return f.Info.Ident.Name
}

var _ info = (*InterfaceInfo)(nil)

type InterfaceInfo struct {
	Info       *Info
	Signatures []*SignatureInfo
}

func (i InterfaceInfo) Doc() string {
	return i.Info.Doc
}

func (i InterfaceInfo) Defs() Definitions {
	return i.Info.Defs
}

func (i InterfaceInfo) Name() string {
	return i.Info.Ident.Name
}

type SignatureInfo struct {
	Info       *Info
	IsEmbedded bool
}

func (s SignatureInfo) Doc() string {
	return s.Info.Doc
}

func (s SignatureInfo) Defs() Definitions {
	return s.Info.Defs
}

var _ info = (*TypeInfo)(nil)

type TypeInfo struct {
	Info *Info
}

func (t TypeInfo) Doc() string {
	return t.Info.Doc
}

func (t TypeInfo) Defs() Definitions {
	return t.Info.Defs
}

func (t TypeInfo) Name() string {
	return t.Info.Ident.Name
}

func (t TypeInfo) IsPointer() *types.Pointer {
	pointer, isPointer := t.Info.Type.(*types.Pointer)
	if !isPointer {
		return nil
	}
	return pointer
}

func (t TypeInfo) IsBasic() *types.Basic {
	basic, isBasic := t.Info.Type.(*types.Basic)
	if !isBasic {
		return nil
	}
	return basic
}

var _ info = (*AliasInfo)(nil)

type AliasInfo struct {
	Info *Info
}

func (a AliasInfo) Doc() string {
	return a.Info.Doc
}

func (a AliasInfo) Defs() Definitions {
	return a.Info.Defs
}

func (a AliasInfo) Name() string {
	return a.Info.Ident.Name
}

func (a AliasInfo) Rhs() types.Type {
	alias := a.Info.Type.(*types.Alias)
	return alias.Rhs()
}

var _ info = (*ImportInfo)(nil)

type ImportInfo struct {
	Info *Info
	Pkgs []*ImportedPackageInfo
}

func (i ImportInfo) Doc() string {
	return i.Info.Doc
}

func (i ImportInfo) Defs() Definitions {
	return i.Info.Defs
}

func (i ImportInfo) Name() string {
	return "IMPORT_STMT"
}

var _ info = (*ImportedPackageInfo)(nil)

type ImportedPackageInfo struct {
	Info *Info
}

func (i ImportedPackageInfo) Doc() string {
	return i.Info.Doc
}

func (i ImportedPackageInfo) Defs() Definitions {
	return i.Info.Defs
}

func (i ImportedPackageInfo) Name() string {
	spec := i.Info.Spec.(*ast.ImportSpec)
	return spec.Name.Name
}

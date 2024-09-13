package main

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

var _ info = (*PackageInfo)(nil)

type PackageInfo struct {
	doc  string
	defs Definitions
	file *ast.File
}

func (p *PackageInfo) Doc() string {
	return p.doc
}

func (p *PackageInfo) Defs() Definitions {
	return p.defs
}

func (p *PackageInfo) Name() string {
	return p.file.Name.Name
}

var _ info = (*MethodInfo)(nil)

type MethodInfo struct {
	doc  string
	defs Definitions
	typ  types.Type
	obj  types.Object

	Decl *ast.FuncDecl
}

func (m MethodInfo) Type() types.Type {
	return m.typ
}

func (m MethodInfo) FuncType() *types.Func {
	return m.obj.(*types.Func)
}

func (m MethodInfo) Doc() string {
	return m.doc
}

func (m MethodInfo) Defs() Definitions {
	return m.defs
}

func (m MethodInfo) Name() string {
	return m.Decl.Name.Name
}

func (m MethodInfo) ReceiverName() string {
	field := m.Decl.Recv.List[0]
	if len(field.Names) == 0 {
		return UnknownMethodName
	}
	return field.Names[0].Name
}

func (m MethodInfo) ReceiverExpr() ast.Expr {
	return m.Decl.Recv.List[0].Type
}

func (m MethodInfo) Params() *ast.FieldList {
	return m.Decl.Type.Params
}

func (m MethodInfo) Returns() *ast.FieldList {
	return m.Decl.Type.Results
}

var _ info = (*FuncInfo)(nil)

type FuncInfo struct {
	doc  string
	defs Definitions
	obj  types.Object
	typ  types.Type

	Decl *ast.FuncDecl
}

func (f FuncInfo) Doc() string {
	return f.doc
}

func (f FuncInfo) Defs() Definitions {
	return f.defs
}

func (f FuncInfo) Name() string {
	return f.Decl.Name.Name
}

func (f FuncInfo) Params() *ast.FieldList {
	return f.Decl.Type.Params
}

func (f FuncInfo) Returns() *ast.FieldList {
	return f.Decl.Type.Results
}

func (f FuncInfo) FuncType() *types.Func {
	return f.obj.(*types.Func)
}

func (f FuncInfo) Type() types.Type {
	return f.typ
}

var _ info = (*ConstInfo)(nil)

type ConstInfo struct {
	doc   string
	defs  Definitions
	idn   *ast.Ident
	value ast.Expr
	typ   types.Type
	obj   types.Object

	decl *ast.GenDecl
}

func (c ConstInfo) Doc() string {
	return c.doc
}

func (c ConstInfo) Defs() Definitions {
	return c.defs
}

func (c ConstInfo) Name() string {
	return c.idn.Name
}

func (c ConstInfo) Value() ast.Expr {
	return c.value
}

func (c ConstInfo) Ident() *ast.Ident {
	return c.idn
}

var _ info = (*VarInfo)(nil)

type VarInfo struct {
	doc   string
	defs  Definitions
	idn   *ast.Ident
	value ast.Expr
	typ   types.Type
	obj   types.Object

	decl *ast.GenDecl
}

func (v VarInfo) Doc() string {
	return v.doc
}

func (v VarInfo) Defs() Definitions {
	return v.defs
}

func (v VarInfo) Name() string {
	return v.idn.Name
}

func (v VarInfo) Value() ast.Expr {
	return v.value
}

func (v VarInfo) Ident() *ast.Ident {
	return v.idn
}

var _ info = (*StructInfo)(nil)

type StructInfo struct {
	doc    string
	defs   Definitions
	fields []*FieldInfo
	idn    *ast.Ident
	typ    types.Type
	obj    types.Object
	spec   *ast.TypeSpec
	decl   *ast.GenDecl
}

func (s StructInfo) Doc() string {
	return s.doc
}

func (s StructInfo) Defs() Definitions {
	return s.defs
}

func (s StructInfo) Name() string {
	return s.idn.Name
}

var _ info = (*FieldInfo)(nil)

type FieldInfo struct {
	idn  *ast.Ident
	doc  string
	expr ast.Expr
	defs Definitions

	typ types.Type
	obj types.Object
}

func (f FieldInfo) Doc() string {
	return f.doc
}

func (f FieldInfo) Defs() Definitions {
	return f.defs
}

func (f FieldInfo) IsEmbedded() bool {
	return f.idn == nil
}

func (f FieldInfo) Name() string {
	return f.idn.Name
}

func (f FieldInfo) Type() types.Type {
	return f.typ
}

var _ info = (*InterfaceInfo)(nil)

type InterfaceInfo struct {
	doc  string
	defs Definitions
	idn  *ast.Ident
	typ  types.Type
	obj  types.Object

	signatures []*SignatureInfo
	decl       *ast.GenDecl
}

func (i InterfaceInfo) Doc() string {
	return i.doc
}

func (i InterfaceInfo) Defs() Definitions {
	return i.defs
}

func (i InterfaceInfo) Name() string {
	return i.idn.Name
}

type SignatureInfo struct {
	doc  string
	defs Definitions
	idn  *ast.Ident
	obj  types.Object
	typ  types.Type

	isEmbedded bool
}

func (s SignatureInfo) Doc() string {
	return s.doc
}

func (s SignatureInfo) Defs() Definitions {
	return s.defs
}

func (s SignatureInfo) SignatureType() *types.Signature {
	return s.typ.(*types.Signature)
}

func (s SignatureInfo) IsEmbedded() bool {
	return s.isEmbedded
}

var _ info = (*TypeInfo)(nil)

type TypeInfo struct {
	doc  string
	defs Definitions
	decl *ast.GenDecl
	idn  *ast.Ident

	typ types.Type
	obj types.Object
}

func (t TypeInfo) Doc() string {
	return t.doc
}

func (t TypeInfo) Defs() Definitions {
	return t.defs
}

func (t TypeInfo) Name() string {
	return t.idn.Name
}

func (t TypeInfo) IsPointer() *types.Pointer {
	pointer, isPointer := t.typ.(*types.Pointer)
	if !isPointer {
		return nil
	}
	return pointer
}

func (t TypeInfo) IsBasic() *types.Basic {
	basic, isBasic := t.typ.(*types.Basic)
	if !isBasic {
		return nil
	}
	return basic
}

var _ info = (*AliasInfo)(nil)

type AliasInfo struct {
	doc  string
	idn  *ast.Ident
	defs Definitions
	typ  types.Type
	obj  types.Object

	alias *types.Alias
}

func (a AliasInfo) Doc() string {
	return a.doc
}

func (a AliasInfo) Defs() Definitions {
	return a.defs
}

func (a AliasInfo) Name() string {
	return a.idn.Name
}

func (a AliasInfo) Rhs() types.Type {
	return a.alias.Rhs()
}

var _ info = (*ImportInfo)(nil)

type ImportInfo struct {
	doc  string
	defs Definitions

	pkgs []*ImportedPackageInfo
}

func (i ImportInfo) Doc() string {
	return i.doc
}

func (i ImportInfo) Defs() Definitions {
	return i.defs
}

func (i ImportInfo) Pkgs() []*ImportedPackageInfo {
	return i.pkgs
}

var _ info = (*ImportedPackageInfo)(nil)

type ImportedPackageInfo struct {
	doc  string
	defs Definitions
	spec *ast.ImportSpec
}

func (i ImportedPackageInfo) Doc() string {
	return i.doc
}

func (i ImportedPackageInfo) Defs() Definitions {
	return i.defs
}

func (i ImportedPackageInfo) Path() *ast.BasicLit {
	return i.spec.Path
}

func (i ImportedPackageInfo) Name() string {
	return i.spec.Name.Name
}

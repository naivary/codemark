package main

import (
	"go/ast"
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

	Import     *ImportInfo
	Consts     []*ConstInfo
	Vars       []*VarInfo
	Funcs      []*FuncInfo
	Methods    []*MethodInfo
	Structs    []*StructInfo
	Types      []*TypeInfo
	Aliases    []*AliasInfo
	Interfaces []*InterfaceInfo
}

func (p PackageInfo) Doc() string {
	return p.doc
}

func (p PackageInfo) Defs() Definitions {
	return p.defs
}

var _ info = (*MethodInfo)(nil)

type MethodInfo struct {
	doc  string
	defs Definitions
	Decl *ast.FuncDecl
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

var _ info = (*ConstInfo)(nil)

type ConstInfo struct {
	doc  string
	defs Definitions
	Decl *ast.GenDecl
	Spec *ast.ValueSpec
}

func (c ConstInfo) Doc() string {
	return c.doc
}

func (c ConstInfo) Defs() Definitions {
	return c.defs
}

func (c ConstInfo) Name() string {
	return c.Spec.Names[0].Name
}

func (c ConstInfo) Value() ast.Expr {
	return c.Spec.Values[0]
}

func (c ConstInfo) Type() ast.Expr {
	return c.Spec.Type
}

var _ info = (*VarInfo)(nil)

type VarInfo struct {
	doc  string
	defs Definitions
	Decl *ast.GenDecl
	Spec *ast.ValueSpec
}

func (v VarInfo) Doc() string {
	return v.doc
}

func (v VarInfo) Defs() Definitions {
	return v.defs
}

var _ info = (*StructInfo)(nil)

type StructInfo struct {
	doc  string
	defs Definitions
	decl *ast.GenDecl
}

func (s StructInfo) Doc() string {
	return s.doc
}

func (s StructInfo) Defs() Definitions {
	return s.defs
}

var _ info = (*InterfaceInfo)(nil)

type InterfaceInfo struct {
	doc  string
	defs Definitions
	decl *ast.GenDecl
}

func (i InterfaceInfo) Doc() string {
	return i.doc
}

func (i InterfaceInfo) Defs() Definitions {
	return i.defs
}

var _ info = (*TypeInfo)(nil)

type TypeInfo struct {
	doc  string
	defs Definitions
	decl *ast.GenDecl
}

func (t TypeInfo) Doc() string {
	return t.doc
}

func (t TypeInfo) Defs() Definitions {
	return t.defs
}

func (t TypeInfo) IsPointer() bool {
	return false
}

var _ info = (*AliasInfo)(nil)

type AliasInfo struct {
	doc  string
	defs Definitions
	decl *ast.GenDecl
}

func (a AliasInfo) Doc() string {
	return a.doc
}

func (a AliasInfo) Defs() Definitions {
	return a.defs
}

var _ info = (*ImportInfo)(nil)

type ImportInfo struct {
	doc  string
	defs Definitions
	decl *ast.GenDecl
}

func (i ImportInfo) Doc() string {
	return i.doc
}

func (i ImportInfo) Defs() Definitions {
	return i.defs
}

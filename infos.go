package main

type Definitions map[string][]any

func (d Definitions) Add(idn string, def any) {
	defs, ok := d[idn]
	if !ok {
		d[idn] = []any{def}
		return
	}
	d[idn] = append(defs, def)
}

func (d Definitions) Get(idn string) ([]any, bool) {
	defs, ok := d[idn]
	return defs, ok
}

type info interface {
	Doc() string
	Defs() Definitions
}

type Infos struct {
	ConstInfos []*ConstInfo
	VarInfos   []*VarInfo
}

func NewInfos() *Infos {
	return &Infos{}
}

var _ info = (*MethodInfo)(nil)

type MethodInfo struct {
	doc  string
	defs Definitions
}

func (m MethodInfo) Doc() string {
	return m.doc
}

func (m MethodInfo) Defs() Definitions {
	return m.defs
}

var _ info = (*FuncInfo)(nil)

type FuncInfo struct {
	doc  string
	defs Definitions
}

func (f FuncInfo) Doc() string {
	return f.doc
}

func (f FuncInfo) Defs() Definitions {
	return f.defs
}

var _ info = (*ConstInfo)(nil)

type ConstInfo struct {
	doc  string
	defs Definitions
}

func (c ConstInfo) Doc() string {
	return c.doc
}

func (c ConstInfo) Defs() Definitions {
	return c.defs
}

var _ info = (*VarInfo)(nil)

type VarInfo struct {
	doc  string
	defs Definitions
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
}

func (t TypeInfo) Doc() string {
	return t.doc
}

func (t TypeInfo) Defs() Definitions {
	return t.defs
}

var _ info = (*AliasInfo)(nil)

type AliasInfo struct {
	doc  string
	defs Definitions
}

func (a AliasInfo) Doc() string {
	return a.doc
}

func (a AliasInfo) Defs() Definitions {
	return a.defs
}

var _ info = (*PackageInfo)(nil)

type PackageInfo struct {
	doc  string
	defs Definitions
}

func (p PackageInfo) Doc() string {
	return p.doc
}

func (p PackageInfo) Defs() Definitions {
	return p.defs
}

var _ info = (*ImportInfo)(nil)

type ImportInfo struct {
	doc  string
	defs Definitions
}

func (i ImportInfo) Doc() string {
	return i.doc
}

func (i ImportInfo) Defs() Definitions {
	return i.defs
}

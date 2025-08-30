package loader

import (
	"github.com/naivary/codemark/marker"
)

type markers interface {
	markers() []marker.Marker
}

type File struct {
	Markers []marker.Marker
}

func (f File) markers() []marker.Marker {
	return f.Markers
}

type VarDecl struct {
	Ident   string
	Value   string
	Markers []marker.Marker
}

func (v VarDecl) markers() []marker.Marker {
	return v.Markers
}

type ConstDecl struct {
	Ident   string
	Value   string
	Markers []marker.Marker
}

func (c ConstDecl) markers() []marker.Marker {
	return c.Markers
}

type FuncDecl struct {
	Ident       string
	ReturnTypes []string
	Params      []string
	Markers     []marker.Marker
}

func (fn FuncDecl) markers() []marker.Marker {
	return fn.Markers
}

type NamedDecl struct {
	Ident   string
	Type    string
	Methods map[string]FuncDecl
	Markers []marker.Marker
}

func (n NamedDecl) markers() []marker.Marker {
	return n.Markers
}

type ImportDecl struct {
	Alias       string
	PackagePath string
	Markers     []marker.Marker
}

func (i ImportDecl) markers() []marker.Marker {
	return i.Markers
}

type StructDecl struct {
	Ident   string
	Fields  map[string]FieldDecl
	Methods map[string]FuncDecl
	Markers []marker.Marker
}

func (s StructDecl) markers() []marker.Marker {
	return s.Markers
}

type FieldDecl struct {
	Ident   string
	Type    string
	Markers []marker.Marker
}

func (f FieldDecl) markers() []marker.Marker {
	return f.Markers
}

type AliasDecl struct {
	Ident   string
	Rhs     string
	Markers []marker.Marker
}

func (a AliasDecl) markers() []marker.Marker {
	return a.Markers
}

type IfaceDecl struct {
	Ident      string
	Signatures map[string]FuncDecl
	Markers    []marker.Marker
}

func (i IfaceDecl) markers() []marker.Marker {
	return i.Markers
}

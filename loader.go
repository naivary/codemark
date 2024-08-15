package main

import (
	"errors"
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"

	"golang.org/x/tools/go/packages"
)

type Loader interface {
	Load(files ...string) (map[string]*Info, error)
}

func NewLoader(cfg *packages.Config) Loader {
	l := &loader{
		methods: make(map[string][]*MethodInfo),
	}
	if cfg == nil {
		l.cfg = l.defaultConfig()
	}
	l.cfg.ParseFile = func(fset *token.FileSet, filename string, src []byte) (*ast.File, error) {
		return parser.ParseFile(fset, filename, src, parser.ParseComments)
	}
	return l
}

var _ Loader = (*loader)(nil)

type loader struct {
	cfg *packages.Config

	consts     []*ConstInfo
	vars       []*VarInfo
	aliases    []*AliasInfo
	structs    []*StructInfo
	basicTypes []*BasicTypeInfo
	ifaces     []*InterfaceInfo
	funcs      []*FuncInfo
	// indexed by receiver type
	methods map[string][]*MethodInfo
}

// types und packages nutzen um die verschiedenen Expression reinzuladen
// docs für die expressions laden
// Marker für die Expression parsen und überprüfen ob die Marker auf diese
// Expression sein dürfen
// Info struct erstellen für die Expression e.g. FuncInfo oder ConstInfo etc.
// Eine Struct mit allen Info als result wiedergeben
func (l *loader) Load(paths ...string) (map[string]*Info, error) {
	infos := make(map[string]*Info, len(paths))
	pkgs, err := packages.Load(l.cfg, paths...)
	if err != nil {
		return nil, err
	}
	if len(pkgs) == 0 {
		return nil, errors.New("empty packages")
	}
	for _, pkg := range pkgs {
		for _, file := range pkg.Syntax {
			l.fileToInfo(pkg, file)
		}
	}
	return infos, nil
}

func (l *loader) fileToInfo(pkg *packages.Package, file *ast.File) (*Info, error) {
	info := &Info{}
	if file.Decls == nil {
		return nil, errors.New("no top-level declarations found")
	}
	types := make([]*ast.GenDecl, 0, 0)
	for _, decl := range file.Decls {
		// FuncDecl -> Method or Func
		funcDecl, isFuncDecl := decl.(*ast.FuncDecl)
		if isFuncDecl {
			l.funcDecl(funcDecl)
			continue
		}
		genDecl := decl.(*ast.GenDecl)
		switch genDecl.Tok {
		case token.CONST:
			l.constDecl(pkg, genDecl)
		case token.VAR:
			l.varDecl(pkg, genDecl)
		case token.IMPORT:
			l.importDecl(genDecl)
		case token.TYPE:
			types = append(types, genDecl)
		}
	}
	for _, typ := range types {
		l.typeDecl(pkg, typ)
	}
	return info, nil
}

func (l *loader) constDecl(pkg *packages.Package, decl *ast.GenDecl) {
	specs := convertSpecs[*ast.ValueSpec](decl.Specs)
	for _, spec := range specs {
		infos := NewConstInfo(spec, pkg)
		l.consts = append(l.consts, infos...)
	}
}

func (l *loader) varDecl(pkg *packages.Package, decl *ast.GenDecl) {
	specs := convertSpecs[*ast.ValueSpec](decl.Specs)
	varInfos := make([]*VarInfo, 0, len(specs))
	for _, spec := range specs {
		infos := NewVarInfo(spec, pkg)
		varInfos = append(varInfos, infos...)
	}
	l.vars = varInfos
}

func (l *loader) importDecl(decl *ast.GenDecl) {}

func (l *loader) funcDecl(decl *ast.FuncDecl) {
	if decl.Recv != nil {
		l.methodDecl(decl)
		return
	}
}

func (l *loader) methodDecl(decl *ast.FuncDecl) {
	typ := decl.Recv.List[0].Type
	ident, isIdent := typ.(*ast.Ident)
	if isIdent {
		info := NewMethodInfo(decl, ident, nil)
		l.methods[ident.Name] = append(l.methods[ident.Name], info)
		return
	}
	pointer := typ.(*ast.StarExpr)
	ptrIdent := pointer.X.(*ast.Ident)
	info := NewMethodInfo(decl, ptrIdent, pointer)
	l.methods[ptrIdent.Name] = append(l.methods[ptrIdent.Name], info)
}

func (l *loader) typeDecl(pkg *packages.Package, decl *ast.GenDecl) {
	typeSpecs := convertSpecs[*ast.TypeSpec](decl.Specs)
	for _, typeSpec := range typeSpecs {
		typ := pkg.TypesInfo.TypeOf(typeSpec.Name)
		alias, isAlias := typ.(*types.Alias)
		if isAlias {
			info := NewAliasInfo(typeSpec, alias, decl)
			l.aliases = append(l.aliases, info)
			continue
		}

		named := typ.(*types.Named).Underlying()
		/*
			strc, isStruct := named.(*types.Struct)
			if isStruct {
				structType := typeSpec.Type.(*ast.StructType)
				continue
			}
		*/
		iface, isIface := named.(*types.Interface)
		if isIface {
			info := NewInterfaceInfo(typeSpec, iface, decl)
			l.ifaces = append(l.ifaces, info)
			continue
		}

		basic, isBasic := named.(*types.Basic)
		if isBasic {
			info := NewBasicTypeInfo(typeSpec, basic, decl, nil)
			info.Methods = l.methods[typeSpec.Name.Name]
			l.basicTypes = append(l.basicTypes, info)
			continue
		}

		ptr, isPtr := named.(*types.Pointer)
		if isPtr {
			basic := ptr.Elem().(*types.Basic)
			info := NewBasicTypeInfo(typeSpec, basic, decl, ptr)
			info.Methods = l.methods[typeSpec.Name.Name]
			l.basicTypes = append(l.basicTypes, info)
		}
	}
}

func (l *loader) defaultConfig() *packages.Config {
	return &packages.Config{
		Mode: packages.NeedTypesInfo | packages.NeedSyntax | packages.NeedTypes,
	}
}

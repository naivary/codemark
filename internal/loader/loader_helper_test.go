package loader

import (
	randv2 "math/rand/v2"
	"os"
	"path"
	"reflect"
	"strconv"
	"strings"
	"text/template"

	"github.com/naivary/codemark/internal/rand"
	"github.com/naivary/codemark/marker"
	"github.com/naivary/codemark/marker/markertest"
)

const _defaultNum = 0

type project struct {
	Structs map[string]StructDecl
	Funcs   map[string]FuncDecl
	Vars    map[string]VarDecl
	Consts  map[string]ConstDecl
	Ifaces  map[string]IfaceDecl
	Aliases map[string]AliasDecl
	Named   map[string]NamedDecl
	Imports map[string]ImportDecl
	Files   map[string]File
}

func newProject() *project {
	p := &project{
		Structs: make(map[string]StructDecl),
		Funcs:   make(map[string]FuncDecl),
		Consts:  make(map[string]ConstDecl),
		Vars:    make(map[string]VarDecl),
		Ifaces:  make(map[string]IfaceDecl),
		Aliases: make(map[string]AliasDecl),
		Named:   make(map[string]NamedDecl),
		Imports: make(map[string]ImportDecl),
		Files:   make(map[string]File),
	}
	p.randomize()
	return p
}

func (p *project) randomize() {
	const n = 10
	var (
		numOfStructs = randv2.IntN(n)
		numOfIfaces  = randv2.IntN(n)
		numOfFuncs   = randv2.IntN(n)
		numOfVars    = randv2.IntN(n)
		numOfConsts  = randv2.IntN(n)
		numOfImports = randv2.IntN(n)
		numOfAliases = randv2.IntN(n)
		numOfNamed   = randv2.IntN(n)
	)
	for range numOfStructs {
		s := randStructDecl()
		p.Structs[s.Ident] = s
	}
	for range numOfIfaces {
		iface := randIfaceDecl()
		p.Ifaces[iface.Ident] = iface
	}
	for range numOfFuncs {
		fn := randFuncDecl()
		p.Funcs[fn.Ident] = fn
	}
	for range numOfVars {
		const randomValue = ""
		v := randVarDecl(randomValue)
		p.Vars[v.Ident] = v
	}
	for range numOfConsts {
		const randomValue = ""
		c := randConstDecl(randomValue)
		p.Consts[c.Ident] = c
	}
	for range numOfImports {
		imp := randImportDecl()
		for {
			_, found := p.Imports[imp.PackagePath]
			if !found {
				break
			}
		}
		p.Imports[imp.PackagePath] = imp
		p.Vars[imp.Use.Ident] = imp.Use
	}
	for range numOfAliases {
		const randomRhs = ""
		alias := randAliasDecl(randomRhs)
		p.Aliases[alias.Ident] = alias
	}
	for range numOfNamed {
		named := randNamedDecl()
		p.Named[named.Ident] = named
	}
}

func (p *project) parse(glob string) (string, error) {
	funcs := template.FuncMap{
		"join": strings.Join,
	}
	tmpl, err := template.New("").Funcs(funcs).ParseGlob(glob)
	if err != nil {
		return "", err
	}
	tmpDir := os.TempDir()
	dir, err := os.MkdirTemp(tmpDir, "cm-project")
	if err != nil {
		return "", err
	}
	for _, t := range tmpl.Templates() {
		filename, _ := strings.CutSuffix(t.Name(), ".tmpl")
		if filename == "go.mod" {
			continue
		}
		p.Files[filename] = File{Markers: randMarkers(_defaultNum)}
	}
	for _, t := range tmpl.Templates() {
		filename, _ := strings.CutSuffix(t.Name(), ".tmpl")
		path := path.Join(dir, filename)
		file, err := os.Create(path)
		if err != nil {
			return "", err
		}
		if err := t.Execute(file, p); err != nil {
			return "", err
		}
	}
	return dir, nil
}

func randAliasDecl(rhs string) AliasDecl {
	if rhs == "" {
		rhs = randType().String()
	}
	return AliasDecl{
		Ident: rand.GoIdent(),
		Rhs:   rhs,
	}
}

func randImportDecl() ImportDecl {
	pkgPath, use := randStdPkgPath()
	return ImportDecl{
		Markers:     randMarkers(_defaultNum),
		PackagePath: pkgPath,
		Use: VarDecl{
			Ident:       rand.GoExpIdent(),
			Value:       use,
			IsImportUse: true,
		},
	}
}

func randFuncDecl() FuncDecl {
	return FuncDecl{
		Ident:   rand.GoIdent(),
		Markers: randMarkers(_defaultNum),
	}
}

func randVarDecl(value string) VarDecl {
	if value == "" {
		value = strconv.Itoa(int(rand.Int64()))
	}
	v := VarDecl{
		Ident:   rand.GoIdent(),
		Value:   value,
		Markers: randMarkers(_defaultNum),
	}
	return v
}

func randConstDecl(value string) ConstDecl {
	if value == "" {
		value = strconv.Itoa(int(rand.Int64()))
	}
	c := ConstDecl{
		Ident:   rand.GoIdent(),
		Value:   value,
		Markers: randMarkers(_defaultNum),
	}
	return c
}

func randStructDecl() StructDecl {
	const (
		numOfFields  = 5
		numOfMethods = 5
	)
	s := StructDecl{
		Ident:   rand.GoIdent(),
		Markers: randMarkers(_defaultNum),
		Fields:  make(map[string]FieldDecl, numOfFields),
		Methods: make(map[string]FuncDecl, numOfMethods),
	}
	for range randv2.IntN(numOfFields) {
		field := FieldDecl{
			Ident:   rand.GoIdent(),
			Type:    randType().String(),
			Markers: randMarkers(_defaultNum),
		}
		s.Fields[field.Ident] = field
	}
	for range randv2.IntN(numOfMethods) {
		method := randFuncDecl()
		s.Methods[method.Ident] = method
	}
	return s
}

func randNamedDecl() NamedDecl {
	const numOfMethods = 5
	n := NamedDecl{
		Ident:   rand.GoIdent(),
		Type:    randType().String(),
		Markers: randMarkers(_defaultNum),
		Methods: make(map[string]FuncDecl, numOfMethods),
	}
	for range randv2.IntN(numOfMethods) {
		method := randFuncDecl()
		n.Methods[method.Ident] = method
	}
	return n
}

func randIfaceDecl() IfaceDecl {
	const numOfSigs = 5
	iface := IfaceDecl{
		Ident:      rand.GoIdent(),
		Markers:    randMarkers(_defaultNum),
		Signatures: make(map[string]FuncDecl, numOfSigs),
	}
	for range randv2.IntN(numOfSigs) {
		sig := FuncDecl{
			Ident:   rand.GoIdent(),
			Markers: randMarkers(_defaultNum),
		}
		iface.Signatures[sig.Ident] = sig
	}
	return iface
}

func randMarkers(n int) []marker.Marker {
	if n == 0 {
		n = 5
	}
	markers := make([]marker.Marker, 0, n)
	for range n {
		m, err := markertest.Rand(randType())
		if err != nil {
			panic(err)
		}
		markers = append(markers, *m)
	}
	return markers
}

func randType() reflect.Type {
	i := randv2.IntN(11) + 1
	switch i {
	case 1:
		return reflect.TypeFor[string]()
	case 2:
		return reflect.TypeFor[int]()
	case 3:
		return reflect.TypeFor[float32]()
	case 4:
		return reflect.TypeFor[complex64]()
	case 5:
		return reflect.TypeFor[bool]()
	case 6:
		return reflect.TypeFor[uint]()
	case 7:
		return reflect.TypeFor[[]string]()
	case 8:
		return reflect.TypeFor[[]int]()
	case 9:
		return reflect.TypeFor[[]float32]()
	case 10:
		return reflect.TypeFor[[]complex64]()
	case 11:
		return reflect.TypeFor[[]bool]()
	}
	return nil
}

// randStdPkgPath randomly selects a Go standard library package and returns two values:
//  1. The package import path as a string.
//  2. A reference to a symbol (function, variable, constant, or type) from that package.
//
// The second return value ensures that the imported package is actually "used",
// preventing the compiler from raising an UnusedImportError.
func randStdPkgPath() (string, string) {
	switch randv2.IntN(19) + 1 {
	case 1:
		return "math", "math.Abs"
	case 2:
		return "strings", "strings.ToUpper"
	case 3:
		return "time", "time.Now"
	case 4:
		return "os", "os.Open"
	case 5:
		return "io", "io.Copy"
	case 6:
		return "bytes", "bytes.NewBuffer"
	case 7:
		return "bufio", "bufio.NewReader"
	case 8:
		return "crypto/md5", "md5.New"
	case 9:
		return "crypto/sha256", "sha256.New"
	case 10:
		return "net/http", "http.Get"
	case 11:
		return "net", "net.Dial"
	case 12:
		return "encoding/json", "json.Marshal"
	case 13:
		return "encoding/base64", "base64.StdEncoding"
	case 14:
		return "sort", "sort.Ints"
	case 15:
		return "strconv", "strconv.Itoa"
	case 16:
		return "context", "context.Background"
	case 17:
		return "errors", "errors.New"
	case 18:
		return "regexp", "regexp.MustCompile"
	case 19:
		return "path/filepath", "filepath.Join"
	default:
		return "", ""
	}
}

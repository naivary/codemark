package testing

import (
	"os"
	"path"
	"reflect"
	"strings"
	"text/template"
	"unicode"

	"github.com/naivary/codemark/parser"
)

type LoaderTestCase struct {
	Dir     string
	Structs map[string]Struct
	Funcs   map[string]Func
	Consts  map[string]Const
	Vars    map[string]Var
	Ifaces  map[string]Iface
	Aliases map[string]Alias
	Named   map[string]Named
	Imports map[string]Import
	Pkgs    []parser.Marker
}

type Import struct {
	Name    string
	Markers []parser.Marker
}

type Alias struct {
	Name    string
	Type    reflect.Type
	Markers []parser.Marker
}

type Named struct {
	Name    string
	Type    reflect.Type
	Markers []parser.Marker
}

type Func struct {
	Name    string
	Fn      reflect.Type
	Markers []parser.Marker
}

type Struct struct {
	Name    string
	Markers []parser.Marker
	Fields  map[string]Field
	Methods map[string]Func
}

type Field struct {
	F       reflect.StructField
	Markers []parser.Marker
}

type Const struct {
	Name    string
	Value   int64
	Markers []parser.Marker
}

type Var struct {
	Name    string
	Value   int64
	Markers []parser.Marker
}

type Iface struct {
	Name       string
	Signatures map[string]Func
	Markers    []parser.Marker
}

func NewGoFiles() (LoaderTestCase, error) {
	tc := LoaderTestCase{
		Structs: make(map[string]Struct),
		Funcs:   make(map[string]Func),
		Consts:  make(map[string]Const),
		Vars:    make(map[string]Var),
		Ifaces:  make(map[string]Iface),
		Aliases: make(map[string]Alias),
		Named:   make(map[string]Named),
		Imports: make(map[string]Import),
		Pkgs:    randMarkers(),
	}
	structQuantity := quantity(6)
	funcQuantity := quantity(4)
	constQuantity := quantity(10)
	ifaceQuantity := quantity(6)
	aliasQuantity := quantity(10)
	for range structQuantity {
		s := RandStruct()
		tc.Structs[s.Name] = s
	}
	for range funcQuantity {
		fn := RandFunc()
		tc.Funcs[fn.Name] = fn
	}
	for range constQuantity {
		c := RandConst()
		tc.Consts[c.Name] = c
	}
	for range constQuantity {
		v := RandVar()
		tc.Vars[v.Name] = v
	}
	for range ifaceQuantity {
		iface := RandIface()
		tc.Ifaces[iface.Name] = iface
	}
	for range aliasQuantity {
		alias := RandAlias()
		tc.Aliases[alias.Name] = alias
	}
	for range ifaceQuantity {
		named := RandNamed()
		tc.Named[named.Name] = named
	}
	for range funcQuantity {
		imported := RandImport()
		for {
			_, found := tc.Imports[imported.Name]
			if found {
				imported = RandImport()
				continue
			}
			break
		}
		tc.Imports[imported.Name] = imported
	}
	tmpl, err := template.ParseGlob("sdk/testing/tmpl/*")
	if err != nil {
		return tc, err
	}
	tmpDir := os.TempDir()
	dir, err := os.MkdirTemp(tmpDir, "cm-project")
	if err != nil {
		return tc, err
	}
	tc.Dir = dir
	for _, t := range tmpl.Templates() {
		name, _ := strings.CutSuffix(t.Name(), ".tmpl")
		path := path.Join(dir, name)
		file, err := os.Create(path)
		if err != nil {
			return tc, err
		}
		if err := t.Execute(file, &tc); err != nil {
			return tc, err
		}
	}
	return tc, err
}

func RandImport() Import {
	return Import{
		Name:    randStdPkg(),
		Markers: randMarkers(),
	}
}

func RandNamed() Named {
	return Named{
		Name:    randName(),
		Markers: randMarkers(),
		Type:    randType(),
	}
}

func RandAlias() Alias {
	return Alias{
		Name:    randName(),
		Markers: randMarkers(),
		Type:    randType(),
	}
}

func RandIface() Iface {
	sigQuantity := quantity(5)
	iface := Iface{
		Name:       randName(),
		Markers:    randMarkers(),
		Signatures: make(map[string]Func),
	}
	for range sigQuantity {
		fn := RandFunc()
		iface.Signatures[fn.Name] = fn
	}
	return iface
}

func RandStruct() Struct {
	fieldQuantity := (randInt64() % 6) + 1
	methodQuantity := (randInt64() % 2) + 1
	fields := make(map[string]Field, fieldQuantity)
	methods := make(map[string]Func, methodQuantity)
	for range fieldQuantity {
		f := randField()
		fields[f.F.Name] = f
	}
	for range methodQuantity {
		m := RandFunc()
		methods[m.Name] = m
	}
	s := Struct{
		Name:    randName(),
		Fields:  fields,
		Markers: randMarkers(),
		Methods: methods,
	}
	return s
}

func RandFunc() Func {
	fn := reflect.FuncOf([]reflect.Type{}, []reflect.Type{}, false)
	return Func{
		Name:    randName(),
		Fn:      fn,
		Markers: randMarkers(),
	}
}

func RandConst() Const {
	return Const{
		Name:    randName(),
		Markers: randMarkers(),
		Value:   randInt64(),
	}
}

func RandVar() Var {
	return Var{
		Name:    randName(),
		Markers: randMarkers(),
		Value:   randInt64(),
	}
}

func randField() Field {
	return Field{
		F: reflect.StructField{
			Name: randName(),
			Type: randType(),
		},
		Markers: randMarkers(),
	}
}

func randStdPkg() string {
	i := (randInt("int64")() % 7) + 1
	switch i {
	case 1:
		return "os"
	case 2:
		return "log/slog"
	case 3:
		return "fmt"
	case 4:
		return "io"
	case 5:
		return "bytes"
	case 6:
		return "flag"
	case 7:
		return "net/http"
	default:
		return "slices"
	}
}

func randMarkers() []parser.Marker {
	markerQuantity := (randInt64() % 5) + 1
	markers := make([]parser.Marker, 0, markerQuantity)
	for range markerQuantity {
		markers = append(markers, *RandMarker(randType()))
	}
	return markers
}

func randType() reflect.Type {
	// string, int, float32, complex64, bool, uint
	i := (randInt("int64")() % 11) + 1
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

func randName() string {
	name := randString()
	for {
		firstLetter := rune(name[0])
		if !unicode.IsDigit(firstLetter) {
			break
		}
		name = randString()
	}
	return name
}

func quantity(mx int) int {
	q := (randInt64() % int64(mx)) + 1
	return int(q)
}

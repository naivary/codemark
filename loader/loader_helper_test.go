package loader

import (
	"os"
	"path"
	"reflect"
	"strings"
	"text/template"

	"github.com/naivary/codemark/parser"
	sdktesting "github.com/naivary/codemark/sdk/testing"
)

// RandLoaderTestCase returns a random LaoderTestCase which can be used to test
// the local loader against randomized markers on any types.
func RandLoaderTestCase() (LoaderTestCase, error) {
	tc := newLoaderTestCase()
	tc.randomize()
	return tc, genRandFiles("tmpl/*", &tc)
}

func genRandFiles(glob string, tc *LoaderTestCase) error {
	tmpl, err := template.ParseGlob(glob)
	if err != nil {
		return err
	}
	tmpDir := os.TempDir()
	dir, err := os.MkdirTemp(tmpDir, "cm-project")
	if err != nil {
		return err
	}
	tc.Dir = dir
	for _, t := range tmpl.Templates() {
		name, _ := strings.CutSuffix(t.Name(), ".tmpl")
		path := path.Join(dir, name)
		file, err := os.Create(path)
		if err != nil {
			return err
		}
		if err := t.Execute(file, &tc); err != nil {
			return err
		}
	}
	return err

}

type LoaderTestCase struct {
	// Absolute Path to a directory containing all the randomly generated files
	// for testing.
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

func newLoaderTestCase() LoaderTestCase {
	return LoaderTestCase{
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
}

func (l LoaderTestCase) randomize() {
	l.randStructs(sdktesting.RandLen)
	l.randFuncs(sdktesting.RandLen)
	l.randAliases(sdktesting.RandLen)
	l.randVars(sdktesting.RandLen)
	l.randConsts(sdktesting.RandLen)
	l.randIfaces(sdktesting.RandLen)
	l.randImports(sdktesting.RandLen)
}

func (l *LoaderTestCase) randStructs(n int) {
	if n <= 0 {
		n = 10
	}
	q := quantity(n)
	for range q {
		s := RandStruct()
		l.Structs[s.Name] = s
	}
}

func (l *LoaderTestCase) randFuncs(n int) {
	if n <= 0 {
		n = 6
	}
	q := quantity(n)
	for range q {
		fn := RandFunc()
		l.Funcs[fn.Name] = fn
	}
}

func (l *LoaderTestCase) randConsts(n int) {
	if n <= 0 {
		n = 10
	}
	q := quantity(n)
	for range q {
		c := RandConst()
		l.Consts[c.Name] = c
	}
}

func (l *LoaderTestCase) randIfaces(n int) {
	if n <= 0 {
		n = 6
	}
	q := quantity(n)
	for range q {
		iface := RandIface()
		l.Ifaces[iface.Name] = iface
	}
}

func (l *LoaderTestCase) randAliases(n int) {
	if n <= 0 {
		n = 10
	}
	q := quantity(n)
	for range q {
		alias := RandAlias()
		l.Aliases[alias.Name] = alias
	}
}

func (l *LoaderTestCase) randVars(n int) {
	if n <= 0 {
		n = 10
	}
	q := quantity(n)
	for range q {
		v := RandVar()
		l.Vars[v.Name] = v
	}
}

func (l *LoaderTestCase) randImports(n int) {
	if n <= 0 {
		n = 4
	}
	q := quantity(n)
	for range q {
		imported := RandImport()
		for {
			_, found := l.Imports[imported.Name]
			if found {
				imported = RandImport()
				continue
			}
			break
		}
		l.Imports[imported.Name] = imported
	}
}

func RandImport() Import {
	return Import{
		Name:    randStdPkg(),
		Markers: randMarkers(),
	}
}

func RandNamed() Named {
	return Named{
		Name:    sdktesting.RandName(),
		Markers: randMarkers(),
		Type:    randType(),
	}
}

func RandAlias() Alias {
	return Alias{
		Name:    sdktesting.RandName(),
		Markers: randMarkers(),
		Type:    randType(),
	}
}

func RandIface() Iface {
	sigQuantity := quantity(5)
	iface := Iface{
		Name:       sdktesting.RandName(),
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
	fieldQuantity := quantity(6)
	methodQuantity := quantity(2)
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
		Name:    sdktesting.RandName(),
		Fields:  fields,
		Markers: randMarkers(),
		Methods: methods,
	}
	return s
}

func randField() Field {
	return Field{
		F: reflect.StructField{
			Name: sdktesting.RandName(),
			Type: randType(),
		},
		Markers: randMarkers(),
	}
}

func RandFunc() Func {
	fn := reflect.FuncOf([]reflect.Type{}, []reflect.Type{}, false)
	return Func{
		Name:    sdktesting.RandName(),
		Fn:      fn,
		Markers: randMarkers(),
	}
}

func RandConst() Const {
	return Const{
		Name:    sdktesting.RandName(),
		Markers: randMarkers(),
		Value:   sdktesting.RandInt64(),
	}
}

func RandVar() Var {
	return Var{
		Name:    sdktesting.RandName(),
		Markers: randMarkers(),
		Value:   sdktesting.RandInt64(),
	}
}

func randStdPkg() string {
	i := quantity(7)
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
	q := quantity(5)
	markers := make([]parser.Marker, 0, q)
	for range q {
		markers = append(markers, *sdktesting.RandMarker(randType()))
	}
	return markers
}

func randType() reflect.Type {
	i := quantity(11)
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

// quantity returns a random number from [1, n)
func quantity(mx int) int {
	q := (sdktesting.RandInt64() % int64(mx)) + 1
	return int(q)
}

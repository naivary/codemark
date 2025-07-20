package loader

import (
	"os"
	"path"
	"reflect"
	"strings"
	"text/template"

	"github.com/naivary/codemark/internal/rand"
	"github.com/naivary/codemark/marker"
	"github.com/naivary/codemark/marker/markertest"
)

// randloaderTestCase returns a random LaoderTestCase which can be used to test
// the local loader against randomized markers on any types.
func randLoaderTestCase() (loaderTestCase, error) {
	tc := newloaderTestCase()
	tc.randomize()
	return tc, randProj("tmpl/*", &tc)
}

func randProj(glob string, tc *loaderTestCase) error {
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
		filename, _ := strings.CutSuffix(t.Name(), ".tmpl")
		if filename == "go.mod" {
			continue
		}
		tc.Files[filename] = File{Markers: randMarkers()}
	}
	for _, t := range tmpl.Templates() {
		filename, _ := strings.CutSuffix(t.Name(), ".tmpl")
		path := path.Join(dir, filename)
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

type loaderTestCase struct {
	// Absolute Path to a directory containing all the randomly generated files
	// for testing.
	Dir string

	Structs map[string]Struct
	Funcs   map[string]Func
	Consts  map[string]Const
	Vars    map[string]Var
	Ifaces  map[string]Iface
	Aliases map[string]Alias
	Named   map[string]Named
	Imports map[string]Import
	Files   map[string]File
}

func newloaderTestCase() loaderTestCase {
	return loaderTestCase{
		Structs: make(map[string]Struct),
		Funcs:   make(map[string]Func),
		Consts:  make(map[string]Const),
		Vars:    make(map[string]Var),
		Ifaces:  make(map[string]Iface),
		Aliases: make(map[string]Alias),
		Named:   make(map[string]Named),
		Imports: make(map[string]Import),
		Files:   make(map[string]File),
	}
}

func (l loaderTestCase) randomize() {
	l.randStructs(rand.RandLen)
	l.randFuncs(rand.RandLen)
	l.randAliases(rand.RandLen)
	l.randVars(rand.RandLen)
	l.randConsts(rand.RandLen)
	l.randIfaces(rand.RandLen)
	l.randImports(rand.RandLen)
	l.randNameds(rand.RandLen)
}

func (l *loaderTestCase) randStructs(n int) {
	if n <= 0 {
		n = 10
	}
	q := quantity(n)
	for range q {
		s := randStruct()
		l.Structs[s.Name] = s
	}
}

func (l *loaderTestCase) randFuncs(n int) {
	if n <= 0 {
		n = 6
	}
	q := quantity(n)
	for range q {
		fn := randFunc()
		l.Funcs[fn.Name] = fn
	}
}

func (l *loaderTestCase) randConsts(n int) {
	if n <= 0 {
		n = 10
	}
	q := quantity(n)
	for range q {
		c := randConst()
		l.Consts[c.Name] = c
	}
}

func (l *loaderTestCase) randIfaces(n int) {
	if n <= 0 {
		n = 6
	}
	q := quantity(n)
	for range q {
		iface := randIface()
		l.Ifaces[iface.Name] = iface
	}
}

func (l *loaderTestCase) randAliases(n int) {
	if n <= 0 {
		n = 10
	}
	q := quantity(n)
	for range q {
		alias := randAlias()
		l.Aliases[alias.Name] = alias
	}
}

func (l *loaderTestCase) randVars(n int) {
	if n <= 0 {
		n = 10
	}
	q := quantity(n)
	for range q {
		v := randVar()
		l.Vars[v.Name] = v
	}
}

func (l *loaderTestCase) randImports(n int) {
	if n <= 0 {
		n = 4
	}
	q := quantity(n)
	for range q {
		imported := randImport()
		for {
			_, found := l.Imports[imported.Name]
			if found {
				imported = randImport()
				continue
			}
			break
		}
		l.Imports[imported.Name] = imported
	}
}

func (l *loaderTestCase) randNameds(n int) {
	if n <= 0 {
		n = 8
	}
	q := quantity(n)
	for range q {
		named := randNamed()
		l.Named[named.Name] = named
	}
}

func randImport() Import {
	return Import{
		Name:    randStdPkg(),
		Markers: randMarkers(),
	}
}

func randNamed() Named {
	n := Named{
		Name:    rand.GoIdent(),
		Markers: randMarkers(),
		Type:    randType(),
		Methods: make(map[string]Func),
	}
	methodQuantity := quantity(2)
	for range methodQuantity {
		method := randFunc()
		n.Methods[method.Name] = method
	}
	return n
}

func randAlias() Alias {
	return Alias{
		Name:    rand.GoIdent(),
		Markers: randMarkers(),
		Type:    randType(),
	}
}

func randIface() Iface {
	sigQuantity := quantity(5)
	iface := Iface{
		Name:       rand.GoIdent(),
		Markers:    randMarkers(),
		Signatures: make(map[string]Func),
	}
	for range sigQuantity {
		fn := randFunc()
		iface.Signatures[fn.Name] = fn
	}
	return iface
}

func randStruct() Struct {
	fieldQuantity := quantity(6)
	methodQuantity := quantity(2)
	fields := make(map[string]Field, fieldQuantity)
	methods := make(map[string]Func, methodQuantity)
	for range fieldQuantity {
		f := randField()
		fields[f.F.Name] = f
	}
	for range methodQuantity {
		m := randFunc()
		methods[m.Name] = m
	}
	s := Struct{
		Name:    rand.GoIdent(),
		Fields:  fields,
		Markers: randMarkers(),
		Methods: methods,
	}
	return s
}

func randField() Field {
	return Field{
		F: reflect.StructField{
			Name: rand.GoIdent(),
			Type: randType(),
		},
		Markers: randMarkers(),
	}
}

func randFunc() Func {
	fn := reflect.FuncOf([]reflect.Type{}, []reflect.Type{}, false)
	return Func{
		Name:    rand.GoIdent(),
		Fn:      fn,
		Markers: randMarkers(),
	}
}

func randConst() Const {
	return Const{
		Name:    rand.GoIdent(),
		Markers: randMarkers(),
		Value:   rand.Int64(),
	}
}

func randVar() Var {
	return Var{
		Name:    rand.GoIdent(),
		Markers: randMarkers(),
		Value:   rand.Int64(),
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

func randMarkers() []marker.Marker {
	q := quantity(5)
	markers := make([]marker.Marker, 0, q)
	for range q {
		m, err := markertest.RandMarker(randType())
		if err != nil {
			panic(err)
		}
		markers = append(markers, *m)
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
	q := (rand.Int64() % int64(mx)) + 1
	return int(q)
}

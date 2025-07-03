package converter

import (
	"fmt"
	"log/slog"
	"reflect"
	"slices"
	"strings"

	"github.com/naivary/codemark/definition/target"
	"github.com/naivary/codemark/parser"
	"github.com/naivary/codemark/parser/marker"
	"github.com/naivary/codemark/sdk"
	sdkutil "github.com/naivary/codemark/sdk/utils"
)

type definitions map[string][]any

func (d definitions) Add(idn string, value any) {
	defs, ok := d[idn]
	if !ok {
		d[idn] = []any{value}
	}
	d[idn] = append(defs, value)
}

var _ sdk.ConverterManager = (*Manager)(nil)

type Manager struct {
	reg   sdk.Registry
	convs map[reflect.Type]sdk.Converter
}

func NewManager(reg sdk.Registry, convs ...sdk.Converter) (*Manager, error) {
	if len(reg.All()) == 0 {
		return nil, sdk.ErrRegistryEmpty
	}
	mngr := &Manager{
		reg:   reg,
		convs: make(map[reflect.Type]sdk.Converter),
	}
	convs = slices.Concat(mngr.builtinConvs(), convs)
	for _, conv := range convs {
		if err := mngr.addConverter(conv); err != nil {
			return nil, err
		}
	}
	return mngr, nil
}

func (m *Manager) GetConverter(rtype reflect.Type) (sdk.Converter, error) {
	conv := m.builtin(rtype)
	if conv != nil {
		return conv, nil
	}
	conv, found := m.convs[rtype]
	if !found {
		return nil, fmt.Errorf("no converter found: %v\n", rtype)
	}
	return conv, nil
}

// addConverter is exactly the same as AddConverter but does not include any
// assertions to be able to use it for internal usage.
func (m *Manager) addConverter(conv sdk.Converter) error {
	for _, rtype := range conv.SupportedTypes() {
		_, found := m.convs[rtype]
		if found {
			return fmt.Errorf("converter already exists: %v\n", conv)
		}
		m.convs[rtype] = conv
	}
	return nil
}

func (m *Manager) AddConverter(conv sdk.Converter) error {
	if err := m.isValidName(conv.Name()); err != nil {
		return fmt.Errorf("converter is not following naming conventions: %s\n", err.Error())
	}
	return m.addConverter(conv)
}

// Convert converts the marker by finding the correlating definition in the
// registry with respect to the target.
func (m *Manager) Convert(mrk marker.Marker, t target.Target) (any, error) {
	idn := mrk.Ident
	def, err := m.reg.Get(idn)
	if err != nil {
		return nil, err
	}
	if inFavorOf, isDepcrecated := def.IsDeprecated(); isDepcrecated {
		msg := fmt.Sprintf("MARKER[%s] IS DEPRECATED IN FAVOR OF `%s`\n", idn, inFavorOf)
		slog.Warn(msg)
	}
	if !(slices.Contains(def.Targets, t) || slices.Contains(def.Targets, target.ANY)) {
		return nil, fmt.Errorf("marker `%s` is appliable to `%v`. Was applied to `%s`", idn, def.Targets, t)
	}
	conv, err := m.GetConverter(def.Output)
	if err != nil {
		return nil, err
	}
	if err := conv.CanConvert(mrk, def); err != nil {
		return nil, err
	}
	out, err := conv.Convert(mrk, def)
	if err != nil {
		return nil, err
	}
	return out.Interface(), nil
}

func (m *Manager) ParseDefs(doc string, t target.Target) (map[string][]any, error) {
	markers, err := parser.Parse(doc)
	if err != nil {
		return nil, err
	}
	defs := make(definitions, len(markers))
	for _, marker := range markers {
		value, err := m.Convert(marker, t)
		if err != nil {
			return nil, err
		}
		defs.Add(marker.Ident, value)
	}
	return defs, nil
}

func (m *Manager) builtinConvs() []sdk.Converter {
	return []sdk.Converter{
		String(),
		Integer(),
		Float(),
		Bool(),
		Complex(),
		List(m),
	}
}

// isValidName checks if the choosen name of a custom converter is following the
// convention of prefixing the name with the project name and that the project
// name is not "codemark".
func (m *Manager) isValidName(name string) error {
	if strings.HasPrefix(name, _codemark) {
		return fmt.Errorf(`the name of your custom converter cannot start with "codemark" because it is reserved for the builtin converters: %s`, name)
	}
	if len(strings.Split(name, ".")) != 2 {
		return fmt.Errorf(`the name of your custom converter has to be seperated with "%s" and must be composed of two segments e.g. "codemark.integer"`, ".")
	}
	return nil
}

// builtin returns a bultin converter if the given rtype can be converterted by
// one of the builtin converters. If no converter is found then nil will be
// returned.
func (m *Manager) builtin(rtype reflect.Type) sdk.Converter {
	if !sdkutil.IsSupported(rtype) {
		return nil
	}
	// rtype can be a native type e.g. string and for that the converter can be
	// retrieved without further validation.
	conv, found := m.convs[rtype]
	if found {
		return conv
	}
	// NOTE: The types choose in the function `reflect.TypeFor` are one of the
	// supported types of the converter. The concrete type choice has no meaning.
	if sdkutil.IsValidSlice(rtype) {
		return m.convs[reflect.TypeFor[[]string]()]
	}
	if sdkutil.IsBool(rtype) {
		return m.convs[reflect.TypeFor[bool]()]
	}
	if sdkutil.IsString(rtype) {
		return m.convs[reflect.TypeFor[string]()]
	}
	if sdkutil.IsInt(rtype) || sdkutil.IsUint(rtype) {
		return m.convs[reflect.TypeFor[int]()]
	}
	if sdkutil.IsFloat(rtype) {
		return m.convs[reflect.TypeFor[float32]()]
	}
	if sdkutil.IsComplex(rtype) {
		return m.convs[reflect.TypeFor[complex64]()]
	}
	return nil
}

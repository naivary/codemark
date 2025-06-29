package converter

import (
	"fmt"
	"log/slog"
	"reflect"
	"slices"
	"strings"

	"github.com/naivary/codemark/parser"
	"github.com/naivary/codemark/sdk"
	sdkutil "github.com/naivary/codemark/sdk/utils"
)

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
	convs = slices.Concat(mngr.defaultConvs(), convs)
	for _, conv := range convs {
		if err := mngr.addConverter(conv); err != nil {
			return nil, err
		}
	}
	return mngr, nil
}

func (c *Manager) GetConverter(rtype reflect.Type) (sdk.Converter, error) {
	conv := c.builtin(rtype)
	if conv != nil {
		return conv, nil
	}
	conv, found := c.convs[rtype]
	if !found {
		return nil, fmt.Errorf("no converter found: %v\n", rtype)
	}
	return conv, nil
}

// addConverter is exactly the same as AddConverter but does not include name
// assertions to be able to use it for internal usage.
func (c *Manager) addConverter(conv sdk.Converter) error {
	for _, rtype := range conv.SupportedTypes() {
		_, found := c.convs[rtype]
		if found {
			return fmt.Errorf("converter already exists: %v\n", conv)
		}
		c.convs[rtype] = conv
	}
	return nil
}

func (c *Manager) AddConverter(conv sdk.Converter) error {
	if err := c.isValidName(conv.Name()); err != nil {
		return fmt.Errorf("converter is not following naming conventions: %s\n", err.Error())
	}
	return c.addConverter(conv)
}

// Convert converts the marker by finding the correlating definition in the
// registry with respect to the target.
func (c *Manager) Convert(mrk parser.Marker, target sdk.Target) (any, error) {
	idn := mrk.Ident()
	def, err := c.reg.Get(idn)
	if err != nil {
		return nil, err
	}
	if inFavorOf, isDepcrecated := def.IsDeprecated(); isDepcrecated {
		msg := fmt.Sprintf("MARKER `%s` IS DEPRECATED IN FAVOR OF `%s`\n", idn, *inFavorOf)
		slog.Warn(msg)
	}
	if !(slices.Contains(def.Targets, target) || slices.Contains(def.Targets, sdk.TargetAny)) {
		return nil, fmt.Errorf("marker `%s` is appliable to `%v`. Was applied to `%s`", idn, def.Targets, target)
	}
	conv, err := c.GetConverter(def.Output)
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

func (c *Manager) ParseDefs(doc string, t sdk.Target) (map[string][]any, error) {
	markers, err := parser.Parse(doc)
	if err != nil {
		return nil, err
	}
	defs := make(map[string][]any, len(markers))
	for _, marker := range markers {
		value, err := c.Convert(marker, t)
		if err != nil {
			return nil, err
		}
		// check if marker is used multiple times and if so append it to the
		// values.
		midn := marker.Ident()
		values, ok := defs[midn]
		if !ok {
			defs[midn] = []any{value}
			continue
		}
		defs[midn] = append(values, value)
	}
	return defs, nil
}

func (c *Manager) defaultConvs() []sdk.Converter {
	return []sdk.Converter{
		&stringConverter{},
		&intConverter{},
		&floatConverter{},
		&boolConverter{},
		&complexConverter{},
		&listConverter{c},
	}
}

// isValidName checks if the choosen name of a custom converter is following the
// convention of prefixing the name with the project name and that the project
// name is not "codemark".
func (c *Manager) isValidName(name string) error {
	if strings.HasPrefix(name, _namePrefix) {
		return fmt.Errorf(`the name of your custom converter cannot start with "codemark" because it is reserved for the builtin converters: %s`, name)
	}
	if len(strings.Split(name, ".")) != 2 {
		return fmt.Errorf(`the name of your custom converter has to be seperated with "%s" and must be composed of two segments e.g. "codemark.integer"`, ".")
	}
	return nil
}

// builtin is checking if the given rtype is convertible using a builtin
// converter. If not nil will be returned.
func (c *Manager) builtin(rtype reflect.Type) sdk.Converter {
	if !sdkutil.IsSupported(rtype) {
		return nil
	}
	// rtype can be a native type e.g. string and for that the converter can be
	// retrieved without further validation.
	conv, found := c.convs[rtype]
	if found {
		return conv
	}
	// NOTE: The types choose in the function `reflect.TypeFor` are one of the
	// supported types of the converter. The concrete choice has no meaning.
	if sdkutil.IsValidSlice(rtype) {
		return c.convs[reflect.TypeFor[[]string]()]
	}
	if sdkutil.IsBool(rtype) {
		return c.convs[reflect.TypeFor[bool]()]
	}
	if sdkutil.IsString(rtype) {
		return c.convs[reflect.TypeFor[string]()]
	}
	if sdkutil.IsInt(rtype) || sdkutil.IsUint(rtype) {
		return c.convs[reflect.TypeFor[int]()]
	}
	if sdkutil.IsFloat(rtype) {
		return c.convs[reflect.TypeFor[float32]()]
	}
	if sdkutil.IsComplex(rtype) {
		return c.convs[reflect.TypeFor[complex64]()]
	}
	return nil
}

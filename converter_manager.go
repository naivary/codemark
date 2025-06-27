package codemark

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

var _ sdk.ConverterManager = (*ConverterManager)(nil)

type ConverterManager struct {
	reg   sdk.Registry
	convs map[reflect.Type]sdk.Converter
}

func NewConvMngr(reg sdk.Registry, convs ...sdk.Converter) (*ConverterManager, error) {
	if len(reg.All()) == 0 {
		return nil, sdk.ErrRegistryEmpty
	}
	mngr := &ConverterManager{
		reg:   reg,
		convs: make(map[reflect.Type]sdk.Converter),
	}
	convs = slices.Concat(mngr.defaultConvs(), convs)
	for _, conv := range convs {
		if err := mngr.AddConverter(conv); err != nil {
			return nil, err
		}
	}
	return mngr, nil
}

func (c *ConverterManager) GetConverter(rtype reflect.Type) (sdk.Converter, error) {
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

func (c *ConverterManager) AddConverter(conv sdk.Converter) error {
	if !c.isValidName(conv.Name()) {
		return fmt.Errorf("%s is reserverd for internal usage: %s\n", _namePrefix, conv.Name())
	}
	for _, rtype := range conv.SupportedTypes() {
		_, found := c.convs[rtype]
		if found {
			return fmt.Errorf("converter already exists: %v\n", conv)
		}
		c.convs[rtype] = conv
	}
	return nil
}

func (c *ConverterManager) Convert(mrk parser.Marker, target sdk.Target) (any, error) {
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

func (c *ConverterManager) ParseDefs(doc string, t sdk.Target) (map[string][]any, error) {
	markers, err := parser.Parse(doc)
	if err != nil {
		return nil, err
	}
	defs := make(map[string][]any, len(markers))
	for _, marker := range markers {
		def, err := c.Convert(marker, t)
		if err != nil {
			return nil, err
		}
		midn := marker.Ident()
		defss, ok := defs[midn]
		if !ok {
			defs[midn] = []any{def}
			continue
		}
		defs[midn] = append(defss, def)
	}
	return defs, nil
}

func (c *ConverterManager) defaultConvs() []sdk.Converter {
	return []sdk.Converter{
		&stringConverter{},
		&intConverter{},
		&floatConverter{},
		&boolConverter{},
		&complexConverter{},
		&listConverter{c},
	}
}

func (c *ConverterManager) isValidName(name string) bool {
	return strings.HasPrefix(name, _namePrefix)
}

// builtin is checking if the given rtype is convertible using a builtin
// converter. If not nil will be returned.
func (c *ConverterManager) builtin(rtype reflect.Type) sdk.Converter {
	if !c.isSupported(rtype) {
		return nil
	}
	// NOTE: The types choose in the function `reflect.TypeFor` are one of the
	// supported types of the converter.
	if c.isValidSlice(rtype) {
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

// isSupported is returning true iff the given rtype is supported by the default
// converters.
func (c *ConverterManager) isSupported(rtype reflect.Type) bool {
	return c.isPrimitive(rtype) || rtype.Kind() == reflect.Slice
}

// isPrimitive is returning true iff the given type is non-slice and a type
// which can be converted by a builtin converter.
func (c *ConverterManager) isPrimitive(rtype reflect.Type) bool {
	return sdkutil.IsInt(rtype) || sdkutil.IsUint(rtype) || sdkutil.IsFloat(rtype) || sdkutil.IsString(rtype) || sdkutil.IsBool(rtype) || sdkutil.IsComplex(rtype)
}

func (c *ConverterManager) isValidSlice(rtype reflect.Type) bool {
	return rtype.Kind() == reflect.Slice && c.isPrimitive(rtype.Elem())
}

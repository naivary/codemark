package converter

import (
	"fmt"
	"log/slog"
	"reflect"

	coreapi "github.com/naivary/codemark/api/core"
	"github.com/naivary/codemark/converter"
	"github.com/naivary/codemark/marker"
	"github.com/naivary/codemark/parser"
	"github.com/naivary/codemark/registry"
	"github.com/naivary/codemark/typeutil"
)

type options map[string][]any

func (o options) Add(idn string, value any) {
	opts, ok := o[idn]
	if !ok {
		o[idn] = []any{value}
	}
	o[idn] = append(opts, value)
}

type Manager struct {
	reg   registry.Registry
	convs map[reflect.Type]converter.Converter
}

func NewManager(reg registry.Registry, convs ...converter.Converter) (*Manager, error) {
	if len(reg.All()) == 0 {
		return nil, registry.ErrRegistryEmpty
	}
	mngr := &Manager{
		reg:   reg,
		convs: make(map[reflect.Type]converter.Converter),
	}
	for _, conv := range convs {
		if err := mngr.addConverter(conv); err != nil {
			return nil, err
		}
	}
	return mngr, nil
}

func (m *Manager) GetConverter(rtype reflect.Type) (converter.Converter, error) {
	conv := m.builtin(rtype)
	if conv != nil {
		return conv, nil
	}
	conv, found := m.convs[rtype]
	if !found {
		return nil, fmt.Errorf("no converter found: %v", rtype)
	}
	return conv, nil
}

// addConverter is exactly the same as AddConverter but does not include any
// assertions to be able to use it for internal usage.
func (m *Manager) addConverter(conv converter.Converter) error {
	for _, rtype := range conv.SupportedTypes() {
		_, found := m.convs[rtype]
		if found {
			return fmt.Errorf("converter already exists: %v", conv)
		}
		m.convs[rtype] = conv
	}
	return nil
}

func (m *Manager) AddConverter(conv converter.Converter) error {
	if err := isValidName(conv.Name()); err != nil {
		return fmt.Errorf("converter is not following naming conventions: %s", err.Error())
	}
	return m.addConverter(conv)
}

// Convert converts the marker to a defined option with respect to the target `t`
func (m *Manager) Convert(mrk marker.Marker, t coreapi.Target) (any, error) {
	idn := mrk.Ident
	opt, err := m.reg.Get(idn)
	if err != nil {
		return nil, err
	}
	if inFavorOf, isDepcrecated := opt.IsDeprecated(); isDepcrecated {
		msg := fmt.Sprintf("%s IS DEPRECATED IN FAVOR OF `%s`\n", idn, inFavorOf)
		slog.Warn(msg)
	}
	if isCorrectTarget(*opt, t) {
		return nil, fmt.Errorf(
			"marker `%s` is appliable to `%v`. Was applied to `%s`",
			idn,
			opt.Targets,
			t,
		)
	}
	conv, err := m.GetConverter(opt.Output)
	if err != nil {
		return nil, err
	}
	if err := conv.CanConvert(mrk, opt.Output); err != nil {
		return nil, err
	}
	out, err := conv.Convert(mrk, opt.Output)
	if err != nil {
		return nil, err
	}
	return out.Interface(), nil
}

func (m *Manager) ParseDefs(doc string, t coreapi.Target) (map[string][]any, error) {
	markers, err := parser.Parse(doc)
	if err != nil {
		return nil, err
	}
	opts := make(options, len(markers))
	for _, marker := range markers {
		value, err := m.Convert(marker, t)
		if err != nil {
			return nil, err
		}
		opts.Add(marker.Ident, value)
	}
	return opts, nil
}

// builtin returns a bultin converter if the given rtype can be converterted by
// one of the builtin converters. If no converter is found then nil will be
// returned.
func (m *Manager) builtin(rtype reflect.Type) converter.Converter {
	if !typeutil.IsSupported(rtype) {
		return nil
	}
	// rtype can be a native type e.g. string and for that the converter can be
	// retrieved without further validation.
	conv, found := m.convs[rtype]
	if found {
		return conv
	}
	if typeutil.IsValidSlice(rtype) {
		return List()
	}
	return Get(rtype)
}

package converter

import (
	"fmt"
	"log/slog"
	"reflect"

	convv1 "github.com/naivary/codemark/api/converter/v1"
	optionapi "github.com/naivary/codemark/api/option"
	"github.com/naivary/codemark/internal/parser"
	"github.com/naivary/codemark/marker"
	"github.com/naivary/codemark/registry"
	"github.com/naivary/codemark/typeutil"
)

type options map[string][]any

func (o options) add(idn string, value any, isUnique bool) error {
	opts, ok := o[idn]
	if !ok {
		o[idn] = []any{value}
		return nil
	}
	if isUnique {
		return fmt.Errorf("option is unique but was used more than once: %s", idn)
	}
	o[idn] = append(opts, value)
	return nil
}

type Manager struct {
	reg   registry.Registry
	convs map[reflect.Type]convv1.Converter
}

func NewManager(reg registry.Registry, convs ...convv1.Converter) (*Manager, error) {
	if len(reg.All()) == 0 {
		return nil, registry.ErrRegistryEmpty
	}
	mngr := &Manager{
		reg:   reg,
		convs: make(map[reflect.Type]convv1.Converter),
	}
	for _, conv := range convs {
		if err := mngr.Add(conv); err != nil {
			return nil, err
		}
	}
	return mngr, nil
}

func (m *Manager) Get(rtype reflect.Type) (convv1.Converter, error) {
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

// addis exactly the same as AddConverter but does not include any
// assertions to be able to use it for internal usage.
func (m *Manager) add(conv convv1.Converter) error {
	for _, rtype := range conv.SupportedTypes() {
		_, found := m.convs[rtype]
		if found {
			return fmt.Errorf("converter already exists: %v", conv)
		}
		m.convs[rtype] = conv
	}
	return nil
}

func (m *Manager) Add(conv convv1.Converter) error {
	if err := isValidName(conv.Name()); err != nil {
		return fmt.Errorf("converter is not following naming conventions: %s", err.Error())
	}
	return m.add(conv)
}

// Convert converts the marker to a defined option with respect to the target `t`
func (m *Manager) Convert(mrk marker.Marker, t optionapi.Target) (any, error) {
	idn := mrk.Ident
	opt, err := m.reg.Get(idn)
	if err != nil {
		return nil, err
	}
	if opt.IsDeprecated() {
		msg := fmt.Sprintf("%s IS DEPRECATED IN FAVOR OF `%s`\n", idn, opt.DeprecatedInFavorOf)
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
	conv, err := m.Get(opt.Output)
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

func (m *Manager) ParseDefs(doc string, t optionapi.Target) (map[string][]any, error) {
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
		opt, err := m.reg.Get(marker.Ident)
		if err != nil {
			return nil, err
		}
		err = opts.add(marker.Ident, value, opt.IsUnique)
		if err != nil {
			return nil, err
		}
	}
	return opts, nil
}

// builtin returns a bultin converter if the given rtype can be converterted by
// one of the builtin converters. If no converter is found then nil will be
// returned.
func (m *Manager) builtin(rtype reflect.Type) convv1.Converter {
	if !typeutil.IsSupported(rtype) {
		return nil
	}
	if typeutil.IsValidSlice(rtype) {
		return NewList(m)
	}
	if typeutil.IsBool(rtype) {
		return NewBool()
	}
	if typeutil.IsString(rtype) {
		return NewString()
	}
	if typeutil.IsInt(rtype) || typeutil.IsUint(rtype) {
		return NewInteger()
	}
	if typeutil.IsFloat(rtype) {
		return NewFloat()
	}
	if typeutil.IsComplex(rtype) {
		return NewComplex()
	}
	if typeutil.IsAny(rtype) {
		return NewAny()
	}
	return nil
}

package converter

import (
	"fmt"
	"log/slog"
	"reflect"

	convv1 "github.com/naivary/codemark/api/converter/v1"
	infov1 "github.com/naivary/codemark/api/info/v1"
	optionv1 "github.com/naivary/codemark/api/option/v1"
	regv1 "github.com/naivary/codemark/api/registry/v1"
	"github.com/naivary/codemark/internal/parser"
	"github.com/naivary/codemark/marker"
	"github.com/naivary/codemark/registry"
	"github.com/naivary/codemark/rtypeutil"
)

type Manager struct {
	reg   regv1.Registry
	convs map[reflect.Type]convv1.Converter
}

func NewManager(reg regv1.Registry, convs ...convv1.Converter) (*Manager, error) {
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
func (m *Manager) Convert(mrk marker.Marker, t optionv1.Target) (any, error) {
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
	conv, err := m.Get(opt.Type)
	if err != nil {
		return nil, err
	}
	if err := conv.CanConvert(mrk, opt.Type); err != nil {
		return nil, err
	}
	out, err := conv.Convert(mrk, opt.Type)
	if err != nil {
		return nil, err
	}
	return out.Interface(), nil
}

func (m *Manager) ParseMarkers(doc string, t optionv1.Target) (infov1.Options, error) {
	markers, err := parser.Parse(doc)
	if err != nil {
		return nil, err
	}
	opts := make(infov1.Options, len(markers))
	for _, marker := range markers {
		value, err := m.Convert(marker, t)
		if err != nil {
			return nil, err
		}
		opt, err := m.reg.Get(marker.Ident)
		if err != nil {
			return nil, err
		}
		err = opts.Add(marker.Ident, value, opt.IsUnique)
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
	if !rtypeutil.IsSupported(rtype) {
		return nil
	}
	if rtypeutil.IsValidSlice(rtype) {
		return NewList(m)
	}
	if rtypeutil.IsBool(rtype) {
		return NewBool()
	}
	if rtypeutil.IsString(rtype) {
		return NewString()
	}
	if rtypeutil.IsInt(rtype) || rtypeutil.IsUint(rtype) {
		return NewInteger()
	}
	if rtypeutil.IsFloat(rtype) {
		return NewFloat()
	}
	if rtypeutil.IsComplex(rtype) {
		return NewComplex()
	}
	if rtypeutil.IsAny(rtype) {
		return NewAny()
	}
	return nil
}

package codemark

import (
	"errors"
	"fmt"
	"log/slog"
	"reflect"
	"slices"

	"github.com/naivary/codemark/parser"
	"github.com/naivary/codemark/sdk"
)

type ConverterManager struct {
	reg sdk.Registry

	convs map[string]sdk.Converter
}

func NewConvMngr(reg sdk.Registry, convs ...sdk.Converter) (*ConverterManager, error) {
	if len(reg.All()) == 0 {
		return nil, errors.New("registry is empty")
	}
	mngr := &ConverterManager{
		reg:   reg,
		convs: make(map[string]sdk.Converter),
	}
	defaultConvs := []sdk.Converter{
		&stringConverter{},
		&intConverter{},
		&floatConverter{},
		&boolConverter{},
		&complexConverter{},
		&listConverter{mngr},
	}
	convs = slices.Concat(defaultConvs, convs)
	for _, conv := range convs {
		if err := mngr.AddConverter(conv); err != nil {
			return nil, err
		}
	}
	// listConv := &listConverter{mngr}
	// if err := mngr.AddConverter(listConv); err != nil {
	// 	return nil, err
	// }
	return mngr, nil
}

func (c *ConverterManager) GetConverter(rtype reflect.Type) (sdk.Converter, error) {
	typeID := TypeID(rtype)
	conv, ok := c.convs[typeID]
	if !ok {
		return nil, fmt.Errorf("converter not found: %s\n", typeID)
	}
	return conv, nil
}

func (c *ConverterManager) AddConverter(conv sdk.Converter) error {
	for _, rtype := range conv.SupportedTypes() {
		typeID := TypeID(rtype)
		_, found := c.convs[typeID]
		if found {
			return fmt.Errorf("converter already exists: %s\n", typeID)
		}
		c.convs[typeID] = conv
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
	if target != def.Target {
		return nil, fmt.Errorf("marker `%s` is appliable to `%s`. Was applied to `%s`", idn, def.Target, target)
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

func (c *ConverterManager) AllConverters() map[string]sdk.Converter {
	return c.convs
}

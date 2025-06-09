package codemark

import (
	"errors"
	"fmt"
	"log/slog"
	"reflect"

	"github.com/naivary/codemark/parser"
)

type ConverterManager struct {
	reg Registry

	converters *tree

	converterList []Converter
}

func NewConvMngr(reg Registry, convs ...Converter) (*ConverterManager, error) {
	if len(reg.All()) == 0 {
		return nil, errors.New("registry is empty")
	}
	mngr := &ConverterManager{
		reg:        reg,
		converters: newTree(),
	}
	for _, conv := range convs {
		if err := mngr.AddConverter(conv); err != nil {
			return nil, err
		}
	}
	return mngr, nil
}

func (c *ConverterManager) GetConverter(rtype reflect.Type) (Converter, error) {
	typeID := TypeID(rtype)
	return c.converters.GetConverter(typeID)
}

func (c *ConverterManager) AddConverter(conv Converter) error {
	for _, rtype := range conv.SupportedTypes() {
		typeID := TypeID(rtype)
		convNode := &node{value: typeID, conv: conv}
		if err := c.converters.Add(convNode); err != nil {
			return err
		}
	}
	return nil
}

func (c *ConverterManager) Convert(mrk parser.Marker, target Target) (any, error) {
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
	typeID := TypeID(def.output)
	conv, err := c.converters.GetConverter(typeID)
	if err != nil {
		return nil, err
	}
	if err := conv.CanConvert(mrk, def); err != nil {
		return nil, err
	}
	return conv.Convert(mrk, def)
}

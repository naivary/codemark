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
}

func NewConvMngr(reg Registry, convs ...Converter) (*ConverterManager, error) {
	if len(reg.All()) == 0 {
		return nil, errors.New("registry is empty")
	}
	root := &node{value: _rootNodeValue, children: make([]*node, 0, len(convs))}
	mngr := &ConverterManager{
		reg:        reg,
		converters: &tree{root},
	}
	for _, conv := range convs {
		if err := mngr.AddConverter(conv); err != nil {
			return nil, err
		}
	}
	return mngr, nil
}

func (c *ConverterManager) GetConverter(rtype reflect.Type) (Converter, error) {
	typeID, err := TypeID(rtype)
	if err != nil {
		return nil, err
	}
	return c.converters.GetConverter(typeID)
}

func (c *ConverterManager) AddConverter(conv Converter) error {
	for _, typ := range conv.SupportedTypes() {
		rtype := reflect.TypeOf(typ)
		typeID, err := TypeID(rtype)
		if err != nil {
			return err
		}
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
	typeID, err := TypeID(def.output)
	if err != nil {
		return nil, err
	}
	conv, err := c.converters.GetConverter(typeID)
	if err := conv.CanConvert(mrk, def); err != nil {
		return nil, err
	}
	return conv.Convert(mrk, def)
}

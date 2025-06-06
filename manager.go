package codemark

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/naivary/codemark/parser"
)

type ConverterManager struct {
	reg Registry

	converters map[string]Converter
}

func NewConvMngr(reg Registry) (*ConverterManager, error) {
	if len(reg.All()) == 0 {
		return nil, errors.New("registry is empty")
	}
	mngr := &ConverterManager{
		reg: reg,
	}
	return mngr, nil
}

func (c *ConverterManager) Convert(mrk parser.Marker, target Target) (any, error) {
	idn := mrk.Ident()
	def := c.reg.Get(idn)
	if def == nil {
		return nil, fmt.Errorf("marker `%s` is not defined in the registry", idn)
	}
	if inFavorOf, isDepcrecated := def.IsDeprecated(); isDepcrecated {
		msg := fmt.Sprintf("MARKER `%s` IS DEPRECATED IN FAVOR OF `%s`\n", idn, *inFavorOf)
		slog.Warn(msg)
	}
	if target != def.Target {
		return nil, fmt.Errorf("marker `%s` is appliable to `%s`. Was applied to `%s`", idn, def.Target, target)
	}
	return nil, nil
}

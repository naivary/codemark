package converter

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/naivary/codemark"
	"github.com/naivary/codemark/marker"
)

type manager struct {
	reg codemark.Registry

	converters map[string]Converter
}

// eigentlich sollten wir keine kinds haben, sondern, wenn wir es nicht regeln
// können sollte es an einen custom manager weiter geleitet werden.
func NewManager(reg codemark.Registry) (*manager, error) {
	if len(reg.All()) == 0 {
		return nil, errors.New("registry is empty")
	}
	mngr := &manager{
		reg: reg,
	}
	// ptr.int -> intconverter
	// ptr.slice -> custom
	// ptr.string -> stringConverter
	// slice.int -> sliceConverter
	// ptr 
	defaultConverters := []any{
		string(""),
		int(0),
		int8(0),
		int16(0),
		int32(0),
		int64(0),
		uint(0),
		uint8(0),
		uint16(0),
		uint32(0),
		uint64(0),
		float32(0.0),
		float64(0.0),
		complex64(1 + 1i),
		complex128(1 + 1i),
		bool(false),
		new(string),
		new(int),
		new(int8),
		new(int16),
		new(int32),
		new(int64),
		new(uint),
		new(uint8),
		new(uint16),
		new(uint32),
		new(uint64),
		new(float32),
		new(float64),
		new(complex64),
		new(complex128),
		new(bool),
	}
	return mngr, nil
}

func (m *manager) Convert(mrk marker.Marker, target codemark.Target) (any, error) {
	idn := mrk.Ident()
	def := m.reg.Get(idn)
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
}

func (m *manager) GetConverter(kindPath string) Converter {
	// TODO: implement decision tree based on given kindPath
	return nil
}

package main

import (
	"fmt"

	"github.com/naivary/codemark/marker"
)

func errImpossibleConv(m marker.Marker, def *Definition) error {
	return fmt.Errorf("cannot convert marker of kind `%v` to definition `%v` of kind `%v`", m.Kind(), def.output, def.kind)
}

package converter

import (
	"reflect"
	"slices"

	convv1 "github.com/naivary/codemark/api/converter/v1"
	"github.com/naivary/codemark/converter/convertertest"
	"github.com/naivary/codemark/marker"
	"github.com/naivary/codemark/marker/markertest"
	"github.com/naivary/codemark/registry/registrytest"
)

func customTypesFor(conv convv1.Converter) []any {
	if _, isList := conv.(*listConverter); isList {
		return registrytest.ListTypes()
	}
	if _, isInt := conv.(*intConverter); isInt {
		return slices.Concat(registrytest.IntTypes(), registrytest.UintTypes())
	}
	if _, isFloat := conv.(*floatConverter); isFloat {
		return registrytest.FloatTypes()
	}
	if _, isComplex := conv.(*complexConverter); isComplex {
		return registrytest.ComplexTypes()
	}
	if _, isBool := conv.(*boolConverter); isBool {
		return registrytest.BoolTypes()
	}
	if _, isString := conv.(*stringConverter); isString {
		return registrytest.StringTypes()
	}
	return nil
}

func validCasesFor(conv convv1.Converter) ([]convertertest.Case, error) {
	types := customTypesFor(conv)
	cases := make([]convertertest.Case, 0, len(types))
	for _, typ := range types {
		rtype := reflect.TypeOf(typ)
		marker, err := markertest.RandMarker(rtype)
		if err != nil {
			return nil, err
		}
		cases = append(cases, makeCase(typ, *marker, true))
	}
	return cases, nil
}

func makeCase(to any, m marker.Marker, isValidCase bool) convertertest.Case {
	rto := reflect.TypeOf(to)
	return convertertest.MustNewCase(&m, rto, isValidCase, getEqualFunc(rto))
}

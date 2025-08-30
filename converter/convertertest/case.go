package convertertest

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/naivary/codemark/marker"
	"github.com/naivary/codemark/marker/markertest"
)

type Case struct {
	// Name of the test case
	Name string
	// Marker to convert by the converter
	Marker marker.Marker
	// Type to convert the marker to.
	To reflect.Type
	// If the test case is a valid or invalid case
	IsValidCase bool
	// Function to validate the value of the converter (after conversion) with
	// the value of the given marker.
	IsEqual func(got, want reflect.Value) bool
}

func NewCase(
	m *marker.Marker,
	to reflect.Type,
	isValidCase bool,
	equal func(got, want reflect.Value) bool,
) (Case, error) {
	var err error
	if m == nil {
		m, err = markertest.Rand(to)
	}
	if err != nil {
		return _casez, err
	}
	if to == nil {
		return _casez, errors.New("to cannot be nil")
	}
	if equal == nil {
		return _casez, fmt.Errorf("func(got, want reflect.Value) bool cannot be nil: %s", m)
	}

	name := fmt.Sprintf("marker[%s] to %v", m.Ident, to)
	tc := Case{
		Name:        name,
		Marker:      *m,
		To:          to,
		IsValidCase: isValidCase,
		IsEqual:     equal,
	}
	return tc, nil
}

func MustNewCase(
	m *marker.Marker,
	to reflect.Type,
	isValidCase bool,
	equal func(got, want reflect.Value) bool,
) Case {
	cas, err := NewCase(m, to, isValidCase, equal)
	if err != nil {
		panic(err)
	}
	return cas
}

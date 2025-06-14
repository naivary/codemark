package testing

import (
	"io"
	"testing"
)

type LoaderTestCase struct {
	File io.Reader
}

type LoaderTester interface {
	Run(t *testing.T, tc LoaderTestCase)
}

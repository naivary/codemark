package testing

import (
	"io"
	"testing"

	"github.com/naivary/codemark/sdk"
)

type LoaderTestCase struct {
	File io.Reader
}

type LoaderTester interface {
	Tests() ([]LoaderTestCase, error)

	Run(t *testing.T, tc ConverterTestCase, mngr sdk.ConverterManager)
}

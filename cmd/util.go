package cmd

import (
	"slices"

	genv1 "github.com/naivary/codemark/api/generator/v1"
	outv1 "github.com/naivary/codemark/api/outputer/v1"
	"github.com/naivary/codemark/generator"
	"github.com/naivary/codemark/internal/generator/openapi"
	outimpl "github.com/naivary/codemark/internal/outputer"
	"github.com/naivary/codemark/outputer"
)

func mustInit[T any](fn func() (T, error)) T {
	r, err := fn()
	if err != nil {
		panic(err)
	}
	return r
}

func newGenManager(cfgFile string, gens []genv1.Generator) (*generator.Manager, error) {
	mngr, err := generator.NewManager(cfgFile)
	if err != nil {
		return nil, err
	}
	builtinGens := []genv1.Generator{
		mustInit(openapi.New),
	}
	for _, gen := range slices.Concat(builtinGens, gens) {
		if err := mngr.Add(gen); err != nil {
			return nil, err
		}
	}
	return mngr, nil
}

func newOutManager(outs []outv1.Outputer) (*outputer.Manager, error) {
	mngr, err := outputer.NewManager(outs...)
	if err != nil {
		return nil, err
	}
	builtinOuts := []outv1.Outputer{
		mustInit(outimpl.NewFsOutputer),
		mustInit(outimpl.NewStdoutOutputer),
	}
	for _, gen := range slices.Concat(builtinOuts, outs) {
		if err := mngr.Add(gen); err != nil {
			return nil, err
		}
	}
	return mngr, nil
}

package main

import "reflect"

type MarkerValues map[string]reflect.Type

type Collector interface {
    Collect(files ...string) (MarkerValues, error)
}

func NewCollector(reg Registry) Collector {
    return &collector{}
}

var _ Collector = (*collector)(nil)

type collector struct{}

func (c collector) Collect(files ...string) (MarkerValues, error) {
    return nil, nil
}

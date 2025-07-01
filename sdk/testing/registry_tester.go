package testing

import "testing"

type RegistryTestCase struct{}

type RegistryTester interface {
	Run(t *testing.T, tc RegistryTestCase)
}

type registryTester struct{}

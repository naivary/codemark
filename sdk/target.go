package sdk

import "fmt"

// Target defines to which type of
// expression a marker is appliable
type Target int

const (
	TargetField Target = iota + 1
	TargetType
	TargetPackage
	TargetFunc
	TargetConst
	TargetVar
	TargetMethod
	TargetInterface
	TargetImport
	TargetAlias
	TargetInterfaceSignature
	TargetStruct
)

var targetNames = map[Target]string{
	TargetField:              "TargetField",
	TargetType:               "TargetType",
	TargetPackage:            "TargetPackage",
	TargetFunc:               "TargetFunc",
	TargetConst:              "TargetConst",
	TargetVar:                "TargetVar",
	TargetMethod:             "TargetMethod",
	TargetInterface:          "TargetInterface",
	TargetImport:             "TargetImport",
	TargetAlias:              "TargetAlias",
	TargetInterfaceSignature: "TargetInterfaceSignature",
	TargetStruct:             "TargetStruct",
}

func (t Target) String() string {
	if name, ok := targetNames[t]; ok {
		return name
	}
	return fmt.Sprintf("Target<%d>", t)
}

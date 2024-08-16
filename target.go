package main

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
	TargetImportStmt
	TargetImportPackage
)

var targetNames = map[Target]string{
	TargetField:         "TargetField",
	TargetType:          "TargetType",
	TargetPackage:       "TargetPackage",
	TargetFunc:          "TargetFunc",
	TargetConst:         "TargetConst",
	TargetVar:           "TargetVar",
	TargetMethod:        "TargetMethod",
	TargetInterface:     "TargetInterface",
	TargetImportStmt:    "TargetImportStmt",
	TargetImportPackage: "TargetImportPackage",
}

func (t Target) String() string {
	if name, ok := targetNames[t]; ok {
        return name
	}
	return fmt.Sprintf("Target<%d>", t)
}

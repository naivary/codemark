package main

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
	TargetImportPackage
)

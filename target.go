package main

// Target defines to which type of
// expression a marker is appliable
type Target int

const (
    // Appliable to struct fields
	TargetField Target = iota + 1
    // Appliable to `type`
	TargetType
    // Appliable to `package`
	TargetPackage
    // Appliable to functions and methods
	TargetFunc
    // Appliable to global constants
	TargetConst
    // Appliable to global variables
	TargetVar

    TargetMethod
)

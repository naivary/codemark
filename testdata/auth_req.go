// this is thed ocumentation for the pkg
package testdata

// import docs
import (
	// this is the fmt docs
	"context"
	"fmt"

	// this is the docs for a named package
	p "go/parser"
	"time"
)

// this is the documentation for multiple vars on the var keyword
var (
	// this is the documentati for one variable in multiple for `multiple`
	multiple int
	vars     string
)

// this is the documentation of the type keyword
type (
	// this is the doc of int
	this int
	are  float32
	many string
	decl complex128
)

type emptyStruct struct{}
type basic int
type slice []int
type ptrMap map[int]int
type array [2]int

type (
	out string
	als = string
	ptr = *struct{}
)

// documentation for interface
// +path:to:maxh=3
type Interface interface {
	// docs for a
	A(context.Context) error
	// this is docs for b
	B() string
	C() AuthRequest
	D() int
}

// this should be gendecl
const (
	// +path:to:max=3
	A int = 2
	B     = ""
)

const (
	Block  = 3
	Of     = '\n'
	Consts = 2 + 3i
)

// this is a doc line for in line
const In, Line = 2, "something"

// name is very important
const NewName = "test-name"

// this is an alias
type Alias = string

// This is the size docs
type Size int

func (s Size) Max() int {
	return 812031283
}

var Comment = p.ParseComments

func (s *Size) Ptr() {}

// this is the doc for the const doc
const (
	// this is the sierequestmax doc
	SizeReququestMax = iota + 1
	SizeReququestMin = iota + 1
)

// +jsonschema:validation=231
// some docs
var Now = time.Now()

// AuthRequest is a request to authenticate
// a user using email and password
type AuthRequest struct {
	// Size of the request
	Size Size `json:"size,omitempty"`

	// Email of the user
	Email string `json:"email,omitempty"`

	// Password is the raw password of the user
	Password string `json:"password,omitempty"`

	// age and length do tell exactly that
	Age, Length int

	Max int
}

func (as *AuthRequest) Ptr() {}

func (AuthRequest) Back(name string) error {
	return nil
}

type AuthEmbedded struct {
	AuthRequest

	Rule string
}

type AuthEmbeddedPtr struct {
	*AuthRequest
}

// some kind of map
var comp = map[complex128]string{
	2 + 3i: "something",
}

var StringMy = fmt.Sprintf("somethign")

// NewAuthReq is creating a new
// authentication request
func NewAuthReq() AuthRequest {
	return AuthRequest{}
}

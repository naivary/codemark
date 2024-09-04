package testdata

// import docs
import (
	// this is the fmt docs
	"context"
	"fmt"

	// this is the time docs
	p "go/parser"
	"time"
)

// TODO: pointer to non basic
// type ptrNonBasicSlice *[]int
//type ptrNonBasic *struct{}
//type PointerBasic *int

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

const in, line = 3, 2

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

const (
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

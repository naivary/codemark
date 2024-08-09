package testdata

import "time"

// name is very important
const ConstName = "test-name"

const (
    block = 3
    of = '\n'
    consts = 2+3i
)

// This is the size docs
type Size int

const (
    SizeReququestMax = iota + 1
    SizeReququestMin = iota + 1
)

// +jsonschema:validation=231
var Now = time.Now()

// AuthRequest is a request to authenticate
// a user using email and password 
type AuthRequest struct {
    // Size of the request
    //
    // +jsonschema:validation:maximum=3
    Size Size `json:"size,omitempty"`

    // Email of the user 
    //
    // +jsonschema:validation:format=email
    Email string `json:"email,omitempty"`

    // Password is the raw password of the user
    // +jsonschema:validation:items={something: 3}
    Password string `json:"password,omitempty"`
}

// some documentatin for send
func (as AuthRequest) Send(name string) error {
    return nil
}

// some kind of map
var comp = map[complex128]string{
    2+3i: "something",
}

// NewAuthReq is creating a new
// authentication request
func NewAuthReq() AuthRequest {
    return AuthRequest{}
}

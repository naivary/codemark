package testdata

import "time"

// name is very important
const Name = "test-name"

// +jsonschema:validation=231
var Now = time.Now()

// AuthRequest is a request to authenticate
// a user using email and password 
type AuthRequest struct {
    // Size of the request
    //
    // +jsonschema:validation:maximum=3
    Size int `json:"size,omitempty"`

    // Email of the user 
    //
    // +jsonschema:validation:format=email
    Email string `json:"email,omitempty"`

    // Password is the raw password of the user
    // +jsonschema:validation:items={something: 3}
    Password string `json:"password,omitempty"`
}

var comp = map[complex128]string{
    2+3i: "something",
}

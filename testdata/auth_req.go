// This is the package documentatino for testdata
package testdata

import "time"

// name is very important
const Name = "test-name"

// current time of the package
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
    Password string `json:"password,omitempty"`
}

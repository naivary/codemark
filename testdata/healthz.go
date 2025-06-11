// docs for package
package testdata

// import keyword
import (
	// fmt package
	"fmt"
	// strings package
	"strings"
)

type (
	Component int
	Name      string
)

type healthStatus int

// const keyword
const (
	// something
	HEALTHY healthStatus = iota + 1
	// next
	DOWN
	PENDING
)

// method doc
func (h healthStatus) String() string {
	return "HEALTHY"
}

// +openapi_v3:general:description="some description"
type HealthzRequest struct {
	// +openapi_v3:validation:required
	SystemComponent string
}

type healthzResponse struct {
	// +openapi_v3:validation:required
	Status, Amount int
}

// isvalid is validating the requst
func (h HealthzRequest) IsValid() bool {
	return false
}

// func doc
func P() {
	strings.HasPrefix("", "")
	fmt.Println("codemark")
}

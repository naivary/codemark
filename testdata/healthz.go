package testdata

// +openapi_v3:general:description="some description"
type HealthzRequest struct {
	// +openapi_v3:validation:required
	SystemComponent string
}

type healthzResponse struct {
	// +openapi_v3:validation:required
	Status, Amount int
}

type healthStatus int

const (
	HEALTHY healthStatus = iota + 1
	DOWN
	PENDING
)

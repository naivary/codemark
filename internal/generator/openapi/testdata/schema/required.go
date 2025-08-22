package schema

// +openapi:schema:description="authentication request data type"
// +openapi:schema:title="authentication request"
type AuthRequest struct {
	// +openapi:schema:required
	Email string

	// +openapi:schema:required
	Password string
}

package schema

// +openapi:schema:description="authentication request data type"
// +openapi:schema:title="authentication request"
type AuthRequest struct {
	// +openapi:schema:format="email"
	// +openapi:schema:required
	Email string

	// +openapi:schema:minLength=12
	// +openapi:schema:maxLength=32
	// +openapi:schema:required
	Password string
}

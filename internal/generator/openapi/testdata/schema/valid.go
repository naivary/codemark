package schema

// +openapi:schema:description="authentication request data type"
type AuthRequest struct {
	// +openapi:schema:format="email"
	Email string

	// +openapi:schema:minLength=12
	// +openapi:schema:maxLength=32
	// +openapi:schema:required
	Password string

	// +openapi:schema:minimum=18
	// +openapi:schema:exclusiveMaximum=99
	// +openapi:schema:required
	Age int
}

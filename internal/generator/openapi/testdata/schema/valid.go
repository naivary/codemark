package schema

// +openapi:schema:description="authentication request data type"
type AuthRequest struct {
	Email string

	// +openapi:schema:minLength=12
	// +openapi:schema:maxLength=32
	Password string

	// +openapi:schema:maximum=99
	// +openapi:schema:minimum=18
	Age int
}

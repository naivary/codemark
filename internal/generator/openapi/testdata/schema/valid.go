package schema

// +openapi:schema:title=""
// +openapi:schema:description=""
type AuthRequest struct {
	// +openapi:schema:required
	// +openapi:schema:format="email"
	Email string

	// +openapi:schema:required
	// +openapi:schema:minLength=12
	// +openapi:schema:maxLength=32
	Password string
}

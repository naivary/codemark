package schema

// +openapi:schema:description="authentication request data type"
type AuthRequest struct {
	// +openapi:schema:format="email"
	Email string

	// +openapi:schema:minLength=12
	// +openapi:schema:maxLength=32
	// +openapi:schema:required
	Password string

	Name string

	// +openapi:schema:minimum=18
	// +openapi:schema:exclusiveMaximum=99
	// +openapi:schema:required
	// +openapi:schema:title="age of the person"
	// +openapi:schema:description="longer desc"
	Age int

	// +openapi:schema:dependentRequired=["name"]
	// +openapi:schema:mutuallyExclusive=["email"]
	// +openapi:schema:required
	UserName string
}

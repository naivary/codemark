package schema

type Email = string

// +openapi:schema:description="authentication request data type"
type AuthRequest struct {
	// +openapi:schema:format="email"
	Email Email

	// +openapi:schema:minLength=12
	// +openapi:schema:maxLength=32
	// +openapi:schema:required
	Password string

	// +openapi:schema:minimum=18
	// +openapi:schema:exclusiveMaximum=99
	// +openapi:schema:required
	// +openapi:schema:title="age of the person"
	// +openapi:schema:description="longer desc"
	Age int

	Args AuthArgs

	// +openapi:schema:uniqueItems
	// +openapi:schema:maxItems=3
	// +openapi:schema:minItems=1
	Slice []string

	AuthArgs []AuthArgs
}

type AuthArgs struct {
	Example string
}

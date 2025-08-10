package schema

// +openapi:schema:description="authentication request data type"
// +openapi:schema:title="authentication request"
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

	// +openapi:schema:dependentRequired=["Name"]
	// +openapi:schema:mutuallyExclusive=["Email", "Name"]
	UserName string

	// +openapi:schema:enum=[3,2,3, "null"]
	NonArr int

	// +openapi:schema:enum=["3", 2, "null"]
	Any []any

	// +openapi:schema:format="email"
	Strings []string

	Iface any

	// +openapi:schema:dependentRequired=["UserName"]
	MapArr map[string][]string

	// +openapi:schema:minItems=5
	Args []AuthArgs
}

type AuthArgs struct{}

package proj

type AuthnRequest struct {
	// +openapi_v3:validation:required
	Username string

	// +openapi_v3:validation:required=true
	Password string
}

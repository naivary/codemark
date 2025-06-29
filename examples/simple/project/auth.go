package proj

type AuthnRequest struct {
	// +openapi_v3:validation:required
	Username string

	Password string
}

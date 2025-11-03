package schema

// +openapi:schema:title="tagging test"
type Tagging struct {
	Email string `json:"EMAIL"`

	Password string `json:"-"`
}

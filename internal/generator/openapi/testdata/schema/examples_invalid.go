package schema

// +openapi:schema:title="invalid examples"
type Examples struct {
	// +openapi:schema:examples=["something"]
	F1 map[string]string
}

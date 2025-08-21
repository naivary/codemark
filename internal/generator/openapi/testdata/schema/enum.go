package schema

// +openapi:schema:title="enum-test"
type Enum struct {
	// +openapi:schema:enum=["e1", "e2"]
	F1 string

	// +openapi:schema:enum=[1, 2]
	F2 int

	// +openapi:schema:enum=[1.1, 2.2]
	F3 float32

	// +openapi:schema:enum=["e1", "e2"]
	F4 []string

	// +openapi:schema:enum=[1, 2]
	F5 []int

	// +openapi:schema:enum=[1.1, 2.2]
	F6 []float32
}

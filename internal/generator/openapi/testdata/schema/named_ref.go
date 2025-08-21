package schema

type NamedType int

// +openapi:schema:description="example"
type Struct struct {
	F1 string
	F2 NamedType
}

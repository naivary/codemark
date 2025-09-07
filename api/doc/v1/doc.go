package v1

type Domain struct {
	Name string
	Desc string
}

type Resource struct {
	Name string
	Desc string
}

type Option struct {
	Desc    string
	Default string
	Type    string
}

type Outputer struct {
	Name string
	Desc string
}

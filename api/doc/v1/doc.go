package v1

type Domain struct {
	Name string
	Desc string
}

type Resource struct {
	Desc string
}

type Option struct {
	Desc    string
	Default string
}

type Outputer struct {
	Name string
	Desc string
}

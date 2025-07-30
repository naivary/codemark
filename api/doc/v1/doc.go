package v1

type Domain struct {
	Desc string
	// Name of the resources of this domain
	Resources []string

	Examples string
}

type Resource struct {
	Desc string
	// Name of the options of the resource
	Options []string
}

type Option struct {
	Desc    string
	Default string
}

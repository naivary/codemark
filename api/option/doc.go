package option

type DomainDoc struct {
	Desc      string
	Resources []ResourceDoc
}

type ResourceDoc struct {
	Desc    string
	Options []Option
}

type OptionDoc struct {
	Desc    string
	Default string
}

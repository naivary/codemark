package v1

type Domain struct {
	Name string
	Desc string
}

type Resource struct {
	Desc string
}

type Option struct {
	// Desc is the detailed description of the option
	Desc string
	// Short summart of the option to show in the listing of the options in the
	// explain coimmand
	Summary string
}

type Outputer struct {
	Name string

	Summary string

	Desc string
}

type Config struct {
	Default     any
	Description string
	Options     map[string]Config
}

package codemark

type config struct {
	// +codemark:config:description="foo"
	// +codemark:config:default="foo"
	foo string

	// +codemark:config:description="foo"
	// +codemark:config:default="bar"
	bar int

	// +codemark:config:description="m"
	// +codemark:config:default="m"
	m map[string]foo
}

type foo struct {
	foo    string
	foobar bool
}

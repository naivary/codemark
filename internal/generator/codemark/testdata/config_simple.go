package codemark

type config struct {
	// +codemark:config:description="djaskdj"
	// +codemark:config:default="foo"
	foo string
	// +codemark:config:description="sedd"
	// +codemark:config:default="bar"
	bar int
	// +codemark:config:description="some desc"
	// +codemark:config:default="foobar"
	foobar bool
}

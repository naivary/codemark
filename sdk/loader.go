package sdk

type Project struct{}

type Loader interface {
	Load(patterns ...string) (*Project, error)
}

package sdk

type Overlayer interface {
	Overlay() (map[string][]byte, error)
}

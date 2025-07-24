package option

type Docer[T any] interface {
	Doc() T
}

type OptionDoc struct {
	Desc    string
	Default string
}

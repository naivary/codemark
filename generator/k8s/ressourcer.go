package k8s

type Ressourcer[T any] interface {
	Resource() string
}

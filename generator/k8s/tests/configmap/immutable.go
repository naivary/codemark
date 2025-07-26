package configmap

// +k8s:meta:name="codemark-test-configmap"
// +k8s:meta:namespace="codemark"
// +k8s:configmap:immutable=true
// +k8s:configmap:keyformat="snake_case"
type ImmutableButEmptyDefault struct {
	// +k8s:configmap:default="4"
	Int int
	// +k8s:configmap:default="1024"
	String string

	// +k8s:configmap:default=""
	NoDefault string

	unexported string
}

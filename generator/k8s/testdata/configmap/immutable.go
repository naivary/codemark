package configmap

// +k8s:metadata:name="codemark-test-configmap"
// +k8s:metadata:namespace="codemark"
// +k8s:configmap:immutable=true
type ImmutableButEmptyDefault struct {
	// +k8s:configmap:default="4"
	Int int
	// +k8s:configmap:default="1024"
	String string

	// +k8s:configmap:default=""
	NoDefault string

	unexported string
}

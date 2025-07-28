package configmap

// +k8s:meta:name="codemark-test-configmap"
// +k8s:meta:namespace="codemark"
type Valid struct {
	// +k8s:configmap:default="4"
	Int int
	// +k8s:configmap:default="1024"
	String string

	NoDefault string

	unexported string
}

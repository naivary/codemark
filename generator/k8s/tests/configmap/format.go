package configmap

// +k8s:meta:name="format"
// +k8s:meta:namespace="codemark"
// +k8s:configmap:format.key="snake_case"
type ValidWithFormat struct {
	// +k8s:configmap:default="4"
	Int int
	// +k8s:configmap:default="1024"
	String string

	NoDefault string

	unexported string
}

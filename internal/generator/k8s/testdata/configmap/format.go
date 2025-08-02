package configmap

// +k8s:metadata:name="format"
// +k8s:metadata:namespace="codemark"
// +k8s:configmap:format.key="snake_case"
type ValidWithFormat struct {
	// +k8s:configmap:default="4"
	Int int
	// +k8s:configmap:default="1024"
	String string

	// +k8s:configmap:default="no-default"
	NoDefault string

	unexported string
}

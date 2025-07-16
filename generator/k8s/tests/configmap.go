package tests

// +k8s:meta:name="my-name"
// +k8s:meta:namespace="default"
// +k8s:configmap:immutable=true
type ConfigMap struct {
	// +k8s:configmap:default="4"
	CPU int

	// +k8s:configmap:default="1024"
	Memory int

	// +k8s:configmap:default="/etc/app/ca.crt"
	TLSPath string

	NoDefault string
}

// +k8s:pod:image="docker.io/naivary/filevault:latest"
// +k8s:pod:imagepullpolicy="Always"
func main() {}

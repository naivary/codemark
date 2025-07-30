package pod

// +k8s:metadata:namespace="my-app"
//
// +k8s:pod:image="docker.io/naivary/filevault:latest"
// +k8s:pod:imagepullpolicy="Always"
func main() {}

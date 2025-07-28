package pod

// +k8s:metadata:namespace="my-app"
//
// +k8s:pod:image="docker.io/naivary/filevault:latest"
// +k8s:pod:imagepullpolicy="Always"
// +k8s:rbac:apigroups=[""]
// +k8s:rbac:verbs=["get", "list"]
// +k8s:rbac:resources=["pod"]
// +k8s:serviceaccount:name="my-app-svc-acc"
func main() {}

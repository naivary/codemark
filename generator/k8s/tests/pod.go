package main

// +k8s:pod:image="docker.io/naivary/filevault:latest"
// +k8s:pod:imagepullpolicy="Always"
// +k8s:rbac:apigroups=[""]
// +k8s:rbac:verbs=["get", "list"]
// +k8s:rbac:resources=["pod"]
// +k8s:meta:namespace="my-app"
func main() {}

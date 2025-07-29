package rbac

// +k8s:rbac:apigroups=[""]
// +k8s:rbac:verbs=["get", "list"]
// +k8s:rbac:resources=["pod"]
// +k8s:serviceaccount:name="my-app-svc-acc"
func main() {}

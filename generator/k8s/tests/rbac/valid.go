package rbac

// +k8s:metadata:name="application-name"
// +k8s:metadata:namespace="my-application"
// -- RBAC MARKERS
// +k8s:rbac:apigroups=[""]
// +k8s:rbac:verbs=["get", "list"]
// +k8s:rbac:resources=["pod"]
func main() {}

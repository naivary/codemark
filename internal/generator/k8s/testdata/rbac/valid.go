package rbac

// +k8s:metadata:name="codemark-rbac"
// +k8s:metadata:namespace="codemark-rbac"
// -- RBAC MARKERS
// +k8s:rbac:apigroups=[""]
// +k8s:rbac:verbs=["get", "list"]
// +k8s:rbac:resources=["pods"]
func main() {}

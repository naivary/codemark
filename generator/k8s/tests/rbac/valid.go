package rbac

// +k8s:metadata:name="my-application-name"
// +k8s:rbac:apigroups=[""]
// +k8s:rbac:verbs=["get", "list"]
// +k8s:rbac:resources=["pod"]
func main() {}

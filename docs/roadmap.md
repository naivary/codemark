# Roadmap

- generator for openapi, (jsonschema), k8s (e.g. rbac, network policies, deployment/pod,
   etc.) (this is the way to create manifests for microservices to take away
   that boilerplace and bound it to documentation (prio: 3)
- variables for marker over sdk (prio: 1)
- time converter (prio: 1)
- map converter (prio: 1)
- Itereator for *Info types of loader to itereator over the defs without nesting
- extending by using codemark as a library
- makewith doc struct not doc just as string
- restructure the project. packages are useless rn. remove some of it. Rename
  Definition to Option because it Describe it much better.
- remove the builtin converters from the manager and dont add them to the map. Just get them using `Get`

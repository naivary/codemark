package k8s

import (
	"encoding/json"
	"fmt"
	"os"

	loaderapi "github.com/naivary/codemark/api/loader"
	corev1 "k8s.io/api/core/v1"
)

type applier interface {
	apply(ident string, cm *corev1.ConfigMap) error
}

func createConfigMap(strc *loaderapi.StructInfo) error {
	cm := &corev1.ConfigMap{
		Data: make(map[string]string),
	}
	for ident, defs := range strc.Defs {
		for _, def := range defs {
			if err := def.(applier).apply(ident, cm); err != nil {
				return err
			}
		}
	}

	for _, field := range strc.Fields {
		if len(field.Defs) == 0 {
			return fmt.Errorf("missing default marker for field: %s\n", field.Idn)
		}
		for _, defs := range field.Defs {
			for _, def := range defs {
				if err := def.(applier).apply(field.Idn.Name, cm); err != nil {
					return err
				}
			}
		}
	}
	return json.NewEncoder(os.Stdout).Encode(cm)
}

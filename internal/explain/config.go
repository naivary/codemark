package explain

import (
	"fmt"
	"io"
	"strings"

	docv1 "github.com/naivary/codemark/api/doc/v1"
)

var _ = map[string]any{
	"schema": map[string]any{
		"draft":     docv1.Config{Default: "", Description: "", Summary: ""},
		"idBaseURL": docv1.Config{},
		"formats": map[string]any{
			"property": docv1.Config{},
			"filename": docv1.Config{},
		},
	},
}

func Config(w io.Writer, dotPath string, configDocs map[string]any) error {
	paths := strings.Split(dotPath, ".")
	for _, p := range paths[:len(paths)-1] {
		var isMap bool
		configDocs, isMap = configDocs[p].(map[string]any)
		if !isMap {
			return fmt.Errorf("dot path is not map: %s", p)
		}
	}
	path := paths[len(paths)-1]
	v := configDocs[path]
	return nil
}

func config(w io.Writer, v any) error {
	return nil
}

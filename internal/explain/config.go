package explain

import (
	"fmt"
	"io"
	"reflect"
	"strings"

	docv1 "github.com/naivary/codemark/api/doc/v1"
)

// openapi.schema
//
//
// OPTIONS:
//   draft <string>
//   <summary>
//
//   idBaseUrl <string>
//   <summary>
//
//   formats <map[string]any>

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
	if doc, isDoc := v.(docv1.Config); isDoc {
		return detailConfigOption(w, doc)
	}
	return configMap(w, v)
}

func configMap(w io.Writer, v any) error {
	for option, doc := range v.(map[string]any) {
		switch v := doc.(type) {
		case docv1.Config:
			fmt.Fprintf(w, "%s (%#v) <%s>\n", option, v.Default, TypeOf(reflect.TypeOf(v.Default)))
			v.Summary = trunc(v.Summary, _truncLen)
			if v.Summary == "" {
				v.Summary = _none
			}
			writeLinesInCol(w, "  %s\n", v.Summary, nil)
			fmt.Fprintf(w, "\n")
		case map[string]any:
			fmt.Fprintf(w, "%s <%s>\n", option, "map[string]any")
		}
	}
	return nil
}

func detailConfigOption(w io.Writer, doc docv1.Config) error {
	fmt.Fprintf(w, "DEFAULT: %#v\n", doc.Default)
	fmt.Fprintf(w, "TYPE: %s\n", TypeOf(reflect.TypeOf(doc.Default)))
	desc := trunc(doc.Description, _truncLen)
	format := "\n  %s"
	if desc == "" {
		desc = _none
		format = "%s\n"
	}
	fmt.Fprintf(w, "DESCRIPTION: ")
	writeLinesInCol(w, format, desc, nil)
	return nil
}

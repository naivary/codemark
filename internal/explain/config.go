package explain

import (
	"fmt"
	"io"
	"reflect"
	"strings"

	docv1 "github.com/naivary/codemark/api/doc/v1"
)

func Config(w io.Writer, dotPath string, configDocs map[string]docv1.Config) error {
	const sep = "."
	paths := strings.Split(dotPath, sep)
	for _, path := range paths[:len(paths)-1] {
		configDocs = configDocs[path].Options
	}
	optName := paths[len(paths)-1]
	doc := configDocs[optName]
	if optName == "" {
		return config(w, docv1.Config{Options: configDocs})
	}
	return config(w, doc)
}

func config(w io.Writer, doc docv1.Config) error {
	if doc.Options == nil {
		return configOptDetail(w, doc)
		// its a options documetnation not the e.g. not a leaf
	}
	// its an overview of all opts available
	for optName, doc := range doc.Options {
		typ := TypeOf(reflect.TypeOf(doc.Default))
		if doc.Options != nil {
			typ = "Object"
		}
		fmt.Fprintf(w, "%s (%#v) %s", optName, doc.Default, typ)
		format := "\n   %s"
		desc := trunc(doc.Description, _truncLen)
		if desc == "" {
			desc = _none
			format = "\n   %s"
		}
		writeLinesInCol(w, format, desc, nil)
		fmt.Fprintln(w)
	}
	return nil
}

func configOptDetail(w io.Writer, doc docv1.Config) error {
	fmt.Fprintf(w, "DEFAULT: %#v\n", doc.Default)
	fmt.Fprintf(w, "TYPE: %s\n", TypeOf(reflect.TypeOf(doc.Default)))
	fmt.Fprintf(w, "DESCRIPTION:")
	format := "\n  %s"
	desc := trunc(doc.Description, _truncLen)
	if desc == "" {
		desc = _none
		format = " %s"
	}
	writeLinesInCol(w, format, desc, nil)
	fmt.Fprintf(w, "\n")
	return nil
}

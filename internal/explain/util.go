package explain

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"slices"
	"strings"
	"text/tabwriter"

	docv1 "github.com/naivary/codemark/api/doc/v1"
	optv1 "github.com/naivary/codemark/api/option/v1"
	"github.com/naivary/codemark/optionutil"
)

const (
	_none = "<none>"
)

func TypeOf(rtype reflect.Type) string {
	var b bytes.Buffer
	return typeOf(rtype, &b)
}

// typeOf is returning the string representation of the type
func typeOf(rtype reflect.Type, b *bytes.Buffer) string {
	if rtype == nil {
		return b.String()
	}
	switch kind := rtype.Kind(); kind {
	case reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64,
		reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64,
		reflect.Float32,
		reflect.Float64,
		reflect.Bool,
		reflect.String:
		b.WriteString(kind.String())
		return typeOf(nil, b)
	case reflect.Slice:
		b.WriteString("[]")
		return typeOf(rtype.Elem(), b)
	case reflect.Array:
		fmt.Fprintf(b, "[%d]", rtype.Len())
		return typeOf(rtype.Elem(), b)
	case reflect.Struct:
		b.WriteString(rtype.String())
		return typeOf(nil, b)
	case reflect.Interface:
		b.WriteString("any")
		return typeOf(nil, b)
	default:
		return "INVALID"
	}
}

// trunc truncates teh string `s` to the length of `n` while respecting
// punctuactions and newlines. This means that there will be longer lines than
// `n` because a sentence isn't finished yet.
func trunc(s string, n int) string {
	var b bytes.Buffer
	pos := 1
	for _, r := range s {
		if pos%n == 0 && r == ' ' {
			fmt.Fprintf(&b, "\n")
			pos = 1
			continue
		}
		if r == '\n' {
			pos = 1
		}
		if pos%n != 0 {
			pos++
		}
		fmt.Fprint(&b, string(r))
	}
	return b.String()
}

func writeLinesInCol(w io.Writer, format, s string, firstLine []any) {
	lines := strings.Split(s, "\n")
	firstLine = append(firstLine, lines[0])
	fmt.Fprintf(w, format, firstLine...)
	numOfEmptyLines := strings.Count(format, "\t")
	emptyLines := []any{}
	for range numOfEmptyLines {
		emptyLines = append(emptyLines, "")
	}
	for _, line := range lines[1:] {
		row := slices.Concat(emptyLines, []any{line})
		fmt.Fprintf(w, format, row...)
	}
}

func resourceDocOf(resources []docv1.Resource, name string) *docv1.Resource {
	for _, resource := range resources {
		if resource.Name == name {
			return &resource
		}
	}
	return nil
}

func optsOf(opts map[string]*optv1.Option, resourceName string) map[string]*optv1.Option {
	res := make(map[string]*optv1.Option, len(opts))
	for ident, opt := range opts {
		if optionutil.ResourceOf(ident) == resourceName {
			res[ident] = opt
		}
	}
	return res
}

func targetsToString(targets []optv1.Target) string {
	slices.Sort(targets)
	s := ""
	for i, target := range targets {
		if i > 0 {
			s += ","
		}
		s += target.String()
	}
	return s
}

func newTabWriter(w io.Writer) *tabwriter.Writer {
	const (
		minWidth = 0
		tabWidth = 0
		padding  = 2
		padChar  = ' '
		flags    = 0
	)
	return tabwriter.NewWriter(w, minWidth, tabWidth, padding, padChar, flags)
}

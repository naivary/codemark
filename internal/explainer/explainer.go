package explainer

import "io"

type Explainer interface {
	Explain(w io.Writer, args ...string) error
}

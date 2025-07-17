package k8s

import (
	"fmt"
	"testing"

	"github.com/naivary/codemark"
	"github.com/naivary/codemark/api"
	"github.com/naivary/codemark/definition/target"
)

func TestConfigMapResourec(t *testing.T) {
	gen, err := NewGenerator()
	if err != nil {
		t.Errorf("err occured: %s\n", err)
	}
	infos, err := codemark.Load(gen.Registry(), "./tests")
	if err != nil {
		t.Errorf("err occured: %s\n", err)
	}
	if err := gen.Generate(infos); err != nil {
		t.Errorf("err occured: %s\n", err)
	}
	doc := api.OptionDoc{
		Ident: "+domain:resource:option",
		Doc: `Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum. Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet. Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum. Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet. Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum. Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet.  

Duis autem vel eum iriure dolor in hendrerit in vulputate velit esse molestie consequat, vel illum dolore eu feugiat nulla facilisis at vero eros et accumsan et iusto odio dignissim qui blandit praesent luptatum zzril delenit augue duis dolore te feugait nulla facilisi. Lorem ipsum dolor sit amet, consectetuer`,
		Default: "REQUIRED",
		Targets: []target.Target{target.ANY, target.FUNC},
	}
	fmt.Println(doc.String())
}

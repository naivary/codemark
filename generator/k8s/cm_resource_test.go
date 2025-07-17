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
		Doc: `Lorem ipsum dolor sit amet, consectetur adipiscing elit. Duis rutrum libero id urna aliquet accumsan. Quisque lectus leo, commodo vel elit sit amet, egestas ornare urna. Etiam eget lobortis ante. Quisque ut augue diam. Donec felis quam, tristique sit amet mattis ut, ullamcorper eget ligula. Morbi bibendum orci non purus pretium, sed porttitor velit porttitor. Sed commodo ipsum in lobortis dapibus. Aenean nisi augue, rutrum quis nulla a, luctus euismod massa. Suspendisse sed enim non tellus ornare vestibulum eget a ante.

In varius in justo et commodo. Aenean a imperdiet est. Nam rutrum accumsan tortor, a feugiat nunc. Donec bibendum justo vel feugiat porta. Duis tempus scelerisque nibh sed ultricies. Curabitur vel tincidunt erat. Duis malesuada finibus elit, id tempus risus dictum nec. Donec nec arcu maximus, interdum risus in, mollis purus. Sed ut ligula finibus, sagittis orci et, egestas purus. Nulla facilisi. Integer molestie vestibulum pharetra.

In sollicitudin congue est a posuere. Lorem ipsum dolor sit amet, consectetur adipiscing elit. Duis lorem lacus, euismod in accumsan sit amet, consequat ut elit. Sed interdum neque vel augue cursus mollis. Aenean efficitur neque turpis, quis mollis nunc facilisis eget. Ut varius eros ac ex facilisis pretium. Etiam venenatis, arcu et commodo.`,
		Default: "REQUIRED",
		Targets: []target.Target{target.ANY, target.FUNC},
	}
	res := api.ResourceDoc{
		Doc: `Lorem ipsum dolor sit amet, consectetur adipiscing elit. Duis rutrum libero id urna aliquet accumsan. Quisque lectus leo, commodo vel elit sit amet, egestas ornare urna. Etiam eget lobortis ante. Quisque ut augue diam. Donec felis quam, tristique sit amet mattis ut, ullamcorper eget ligula. Morbi bibendum orci non purus pretium, sed porttitor velit porttitor. Sed commodo ipsum in lobortis dapibus. Aenean nisi augue, rutrum quis nulla a, luctus euismod massa. Suspendisse sed enim non tellus ornare vestibulum eget a ante.

In varius in justo et commodo. Aenean a imperdiet est. Nam rutrum accumsan tortor, a feugiat nunc. Donec bibendum justo vel feugiat porta. Duis tempus scelerisque nibh sed ultricies. Curabitur vel tincidunt erat. Duis malesuada finibus elit, id tempus risus dictum nec. Donec nec arcu maximus, interdum risus in, mollis purus. Sed ut ligula finibus, sagittis orci et, egestas purus. Nulla facilisi. Integer molestie vestibulum pharetra.

In sollicitudin congue est a posuere. Lorem ipsum dolor sit amet, consectetur adipiscing elit. Duis lorem lacus, euismod in accumsan sit amet, consequat ut elit. Sed interdum neque vel augue cursus mollis. Aenean efficitur neque turpis, quis mollis nunc facilisis eget. Ut varius eros ac ex facilisis pretium. Etiam venenatis, arcu et commodo.`,
		Options: []api.OptionDoc{doc},
	}
	fmt.Println(res.String())
}

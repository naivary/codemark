package k8s

import (
	"reflect"

	loaderapi "github.com/naivary/codemark/api/loader"
	"github.com/naivary/codemark/definition"
	"github.com/naivary/codemark/definition/target"
	"github.com/naivary/codemark/maker"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// TODO: Create defs dynamiclly by having a Doc() function for every of the
// options and build it using reflect

func newPod() *corev1.Pod {
	return &corev1.Pod{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Pod",
		},
		Spec: corev1.PodSpec{
			Containers: make([]corev1.Container, 1, 1),
		},
	}
}

func podDefs() []*definition.Definition {
	return []*definition.Definition{
		maker.MustMakeDefWithDoc(podIdent("image"), reflect.TypeFor[Image](), "image to use for the container", target.FUNC),
		maker.MustMakeDefWithDoc(podIdent("imagepullpolicy"), reflect.TypeFor[ImagePullPolicy](), "pull policy for the container", target.FUNC),
	}
}

func createPod(fn loaderapi.FuncInfo) (*corev1.Pod, error) {
	pod := newPod()
	pod.ObjectMeta = createObjectMeta(fn)
	for _, defs := range fn.Defs {
		for _, def := range defs {
			var err error
			switch d := def.(type) {
			case Image:
				err = d.apply(&pod.Spec.Containers[0])
			case ImagePullPolicy:
				err = d.apply(&pod.Spec.Containers[0])
			}
			if err != nil {
				return nil, err
			}
		}
	}
	return pod, nil
}

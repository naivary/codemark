package k8s

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/naivary/codemark/api/core"
	loaderapi "github.com/naivary/codemark/api/loader"
)

func newPod() *corev1.Pod {
	return &corev1.Pod{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Pod",
		},
		Spec: corev1.PodSpec{
			Containers: make([]corev1.Container, 1),
		},
	}
}

func podOpts() []*core.Option {
	const resource = "pod"
	return makeDefs(resource, map[any][]core.Target{
		Image(""):           {core.TargetFunc},
		ImagePullPolicy(""): {core.TargetFunc},
	})
}

func createPod(fn loaderapi.FuncInfo) (*corev1.Pod, error) {
	pod := newPod()
	objectMeta, err := createObjectMeta(fn)
	if err != nil {
		return nil, err
	}
	pod.ObjectMeta = *objectMeta
	for _, opts := range fn.Opts {
		for _, opt := range opts {
			var err error
			switch o := opt.(type) {
			case Image:
				err = o.apply(&pod.Spec.Containers[0])
			case ImagePullPolicy:
				err = o.apply(&pod.Spec.Containers[0])
			}
			if err != nil {
				return nil, err
			}
		}
	}
	return pod, nil
}

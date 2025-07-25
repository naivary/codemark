package k8s

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	loaderapi "github.com/naivary/codemark/api/loader"
)

func newPod(fn loaderapi.FuncInfo) (corev1.Pod, error) {
	pod := corev1.Pod{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Pod",
		},
		Spec: corev1.PodSpec{
			Containers: make([]corev1.Container, 1),
		},
	}
	objectMeta, err := createObjectMeta(fn)
	if err != nil {
		return pod, err
	}
	pod.ObjectMeta = objectMeta
	return pod, nil
}

func createPod(fn loaderapi.FuncInfo) (corev1.Pod, error) {
	pod, err := newPod(fn)
	if err != nil {
		return pod, err
	}
	for _, opts := range fn.Opts {
		for _, opt := range opts {
			if err := applyOptToPod(opt, &pod); err != nil {
				return pod, err
			}
		}
	}
	return pod, nil
}

func applyOptToPod(opt any, pod *corev1.Pod) error {
	var err error
	switch o := opt.(type) {
	case Image:
		err = o.apply(&pod.Spec.Containers[0])
	case ImagePullPolicy:
		err = o.apply(&pod.Spec.Containers[0])
	}
	return err
}

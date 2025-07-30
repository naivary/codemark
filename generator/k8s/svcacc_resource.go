package k8s

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	genv1 "github.com/naivary/codemark/api/generator/v1"
	infov1 "github.com/naivary/codemark/api/info/v1"
)

func newServiceAccount(info infov1.Info) (corev1.ServiceAccount, error) {
	sva := corev1.ServiceAccount{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "ServiceAccount",
		},
	}
	objectMeta, err := createObjectMeta(info)
	if err != nil {
		return sva, err
	}
	sva.ObjectMeta = objectMeta
	return sva, nil
}

func createServiceAccount(info infov1.Info) (*genv1.Artifact, error) {
	svcAcc, err := newServiceAccount(info)
	if err != nil {
		return nil, err
	}
	for _, opts := range info.Options() {
		for _, opt := range opts {
			if err := applyOptToServiceAccount(opt, &svcAcc); err != nil {
				return nil, err
			}
		}
	}
	filename := fmt.Sprintf("%s.sva.yaml", svcAcc.Name)
	return newArtifact(filename, svcAcc)
}

func applyOptToServiceAccount(opt any, svcAcc *corev1.ServiceAccount) error {
	var err error
	switch o := opt.(type) {
	case ServiceAccountName:
		err = o.apply(svcAcc)
	}
	return err
}

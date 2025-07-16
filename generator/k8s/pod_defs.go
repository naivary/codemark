package k8s

import (
	"fmt"
	"slices"

	corev1 "k8s.io/api/core/v1"
)

func podIdent(option string) string {
	return fmt.Sprintf("k8s:pod:%s", option)
}

type Image string

func (i Image) apply(c *corev1.Container) error {
	c.Image = string(i)
	return nil
}

type ImagePullPolicy string

func (i ImagePullPolicy) apply(c *corev1.Container) error {
	options := []corev1.PullPolicy{corev1.PullAlways, corev1.PullIfNotPresent, corev1.PullNever}
	if !slices.Contains(options, corev1.PullPolicy(i)) {
		return fmt.Errorf("value for ImagePullPolicy is not valid: %s. Valid options are: %v\n", i, options)
	}
	c.ImagePullPolicy = corev1.PullPolicy(i)
	return nil
}

package k8s

import (
	"fmt"
	"slices"

	"github.com/naivary/codemark/api"
	"github.com/naivary/codemark/definition/target"
	corev1 "k8s.io/api/core/v1"
)

type Image string

func (i Image) apply(c *corev1.Container) error {
	c.Image = string(i)
	return nil
}

func (i Image) Doc() api.OptionDoc {
	return api.OptionDoc{
		Targets: []target.Target{target.FUNC},
		Doc:     "Image to use for the container",
		Default: "REQUIRED",
	}
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

func (i ImagePullPolicy) Doc() api.OptionDoc {
	return api.OptionDoc{
		Targets: []target.Target{target.FUNC},
		Doc:     "Image Pull Policy of the container image",
		Default: "Always",
	}
}

package k8s

import (
	"fmt"
	"slices"

	corev1 "k8s.io/api/core/v1"

	optionapi "github.com/naivary/codemark/api/option"
)

type Image string

func (i Image) apply(c *corev1.Container) error {
	c.Image = string(i)
	return nil
}

func (i Image) Doc() optionapi.OptionDoc {
	return optionapi.OptionDoc{
		Desc:    `Image to use as the container`,
		Default: "",
	}
}

type ImagePullPolicy string

func (i ImagePullPolicy) apply(c *corev1.Container) error {
	options := []corev1.PullPolicy{corev1.PullAlways, corev1.PullIfNotPresent, corev1.PullNever}
	if !slices.Contains(options, corev1.PullPolicy(i)) {
		return fmt.Errorf(
			"value for ImagePullPolicy is not valid: %s. Valid options are: %v",
			i,
			options,
		)
	}
	c.ImagePullPolicy = corev1.PullPolicy(i)
	return nil
}

func (i ImagePullPolicy) Doc() optionapi.OptionDoc {
	return optionapi.OptionDoc{
		Desc:    `Image Pull Policy of the container`,
		Default: "Always",
	}
}

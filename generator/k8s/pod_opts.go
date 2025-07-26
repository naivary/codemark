package k8s

import (
	"fmt"
	"slices"

	corev1 "k8s.io/api/core/v1"

	docv1 "github.com/naivary/codemark/api/doc/v1"
	optionv1 "github.com/naivary/codemark/api/option/v1"
)

const _podResource = "pod"

func podOpts() []*optionv1.Option {
	return makeOpts(_podResource,
		mustMakeOpt(_typeName, Image(""), _repetable, optionv1.TargetFunc),
		mustMakeOpt(_typeName, ImagePullPolicy(""), _repetable, optionv1.TargetFunc),
	)
}

type Image string

func (i Image) apply(c *corev1.Container) error {
	c.Image = string(i)
	return nil
}

func (i Image) Doc() docv1.Option {
	return docv1.Option{
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

func (i ImagePullPolicy) Doc() docv1.Option {
	return docv1.Option{
		Desc:    `Image Pull Policy of the container`,
		Default: "Always",
	}
}

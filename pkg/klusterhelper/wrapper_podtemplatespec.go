package klusterhelper

import (
	"maps"

	corev1 "k8s.io/api/core/v1"
)

type PodTemplateSpecWrapper struct {
	*corev1.PodTemplateSpec
}

func (t *PodTemplateSpecWrapper) MergeLabels(labels map[string]string) *PodTemplateSpecWrapper {
	if t.ObjectMeta.Labels == nil {
		t.ObjectMeta.Labels = make(map[string]string)
	}
	maps.Copy(t.ObjectMeta.Labels, labels)
	return t
}
func (t *PodTemplateSpecWrapper) AddInitContainer(c *ContainerWrapper) *PodTemplateSpecWrapper {
	t.Spec.InitContainers = append(t.Spec.InitContainers, *c.Container)
	return t
}
func (t *PodTemplateSpecWrapper) AddContainer(c *ContainerWrapper) *PodTemplateSpecWrapper {
	t.Spec.Containers = append(t.Spec.Containers, *c.Container)
	return t
}
func (t *PodTemplateSpecWrapper) AddVolumes(v ...corev1.Volume) *PodTemplateSpecWrapper {
	t.Spec.Volumes = append(t.Spec.Volumes, v...)
	return t
}
func (t *PodTemplateSpecWrapper) WithSecurityContext(ctx *corev1.PodSecurityContext) *PodTemplateSpecWrapper {
	t.Spec.SecurityContext = ctx
	return t
}

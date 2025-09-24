package klusterhelper

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

type ContainerWrapper struct {
	*corev1.Container
}

func (c *ContainerWrapper) WithPort(port int32) *ContainerWrapper {
	c.Ports = append(c.Ports, corev1.ContainerPort{
		HostPort:      port,
		ContainerPort: port,
	})
	return c
}
func (c *ContainerWrapper) WithCPULimit(cpu string) *ContainerWrapper {
	if c.Resources.Limits == nil {
		c.Resources.Limits = corev1.ResourceList{}
	}
	c.Resources.Limits[corev1.ResourceCPU] = resource.MustParse(cpu)
	return c
}
func (c *ContainerWrapper) WithMemoryLimit(memory string) *ContainerWrapper {
	if c.Resources.Limits == nil {
		c.Resources.Limits = corev1.ResourceList{}
	}
	c.Resources.Limits[corev1.ResourceMemory] = resource.MustParse(memory)
	return c
}
func (c *ContainerWrapper) WithCPURequest(cpu string) *ContainerWrapper {
	if c.Resources.Requests == nil {
		c.Resources.Requests = corev1.ResourceList{}
	}
	c.Resources.Requests[corev1.ResourceCPU] = resource.MustParse(cpu)
	return c
}
func (c *ContainerWrapper) WithMemoryRequest(memory string) *ContainerWrapper {
	if c.Resources.Requests == nil {
		c.Resources.Requests = corev1.ResourceList{}
	}
	c.Resources.Requests[corev1.ResourceMemory] = resource.MustParse(memory)
	return c
}
func (c *ContainerWrapper) AddEnvVar(name, value string) *ContainerWrapper {
	c.Env = append(c.Env, corev1.EnvVar{
		Name:  name,
		Value: value,
	})
	return c
}
func (c *ContainerWrapper) AddEnvFromSecret(name string) *ContainerWrapper {
	c.EnvFrom = append(c.EnvFrom, corev1.EnvFromSource{
		SecretRef: &corev1.SecretEnvSource{
			LocalObjectReference: corev1.LocalObjectReference{
				Name: name,
			},
		},
	})
	return c
}
func (c *ContainerWrapper) WithImagePullPolicy(policy corev1.PullPolicy) *ContainerWrapper {
	c.ImagePullPolicy = policy
	return c
}
func (c *ContainerWrapper) AddCommands(commands ...string) *ContainerWrapper {
	c.Command = append(c.Command, commands...)
	return c
}
func (c *ContainerWrapper) MountVolume(name, mountPath string) *ContainerWrapper {
	c.VolumeMounts = append(c.VolumeMounts, corev1.VolumeMount{
		Name:      name,
		MountPath: mountPath,
	})
	return c
}

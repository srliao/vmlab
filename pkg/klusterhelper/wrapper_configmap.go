package klusterhelper

import (
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/yaml"
)

type ConfigMapWrapper struct {
	*corev1.ConfigMap
}

var _ KubeResource = &ConfigMapWrapper{}

func (c *ConfigMapWrapper) validate() error          { return nil }
func (c *ConfigMapWrapper) marshal() ([]byte, error) { return yaml.Marshal(c.ConfigMap) }
func (c *ConfigMapWrapper) WithData(data map[string]string) *ConfigMapWrapper {
	c.Data = data
	return c
}

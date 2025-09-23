package klusterhelper

import (
	"maps"
	"path/filepath"

	fluxv1 "github.com/fluxcd/kustomize-controller/api/v1"
	"github.com/fluxcd/pkg/apis/meta"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/kustomize/api/types"
	"sigs.k8s.io/yaml"
)

type (
	// TOP LEVEL RESOURCES
	ConfigMapWrapper struct {
		*corev1.ConfigMap
	}
	DeploymentWrapper struct {
		*appsv1.Deployment
	}
	FluxKustomizationWrapper struct {
		*fluxv1.Kustomization
		subpath string
	}
	IngressWrapper struct {
		*networkingv1.Ingress
	}
	KustomizationWrapper struct {
		*types.Kustomization
	}
	ReplicaSetWrapper struct {
		*appsv1.ReplicaSet
	}
	ServiceWrapper struct {
		*corev1.Service
	}

	// NESTED RESOURCES
	ContainerWrapper struct {
		*corev1.Container
	}
	PodTemplateSpecWrapper struct {
		*corev1.PodTemplateSpec
	}
)

// TOP LEVEL RESOURCES
func (c *ConfigMapWrapper) validate() error          { return nil }
func (c *ConfigMapWrapper) marshal() ([]byte, error) { return yaml.Marshal(c.ConfigMap) }
func (c *ConfigMapWrapper) WithData(data map[string]string) *ConfigMapWrapper {
	c.Data = data
	return c
}

func (d *DeploymentWrapper) validate() error          { return nil }
func (d *DeploymentWrapper) marshal() ([]byte, error) { return marshalCleanYAML(d.Deployment) }
func (d *DeploymentWrapper) WithAnnotations(annotations map[string]string) *DeploymentWrapper {
	if d.Annotations == nil {
		d.Annotations = make(map[string]string)
	}
	maps.Copy(d.Annotations, annotations)
	return d
}
func (d *DeploymentWrapper) WithReplicas(replicas int32) *DeploymentWrapper {
	d.Spec.Replicas = &replicas
	return d
}
func (d *DeploymentWrapper) MergeLabels(labels map[string]string) *DeploymentWrapper {
	if d.Labels == nil {
		d.Labels = make(map[string]string)
	}
	maps.Copy(labels, d.Labels)
	return d
}
func (d *DeploymentWrapper) WithPodTemplate(t *PodTemplateSpecWrapper) *DeploymentWrapper {
	d.Spec.Template = *t.PodTemplateSpec
	return d
}

func (f *FluxKustomizationWrapper) marshal() ([]byte, error) {
	return marshalCleanYAML(f.Kustomization)
}
func (f *FluxKustomizationWrapper) validate() error { return nil }
func (f *FluxKustomizationWrapper) WithPath(base, subpath string) *FluxKustomizationWrapper {
	f.Kustomization.Spec.Path = filepath.Join(base, subpath)
	f.subpath = subpath
	return f
}
func (f *FluxKustomizationWrapper) WithDependsOn(deps []meta.NamespacedObjectReference) *FluxKustomizationWrapper {
	f.Kustomization.Spec.DependsOn = deps
	return f
}

func (i *IngressWrapper) validate() error          { return nil }
func (i *IngressWrapper) marshal() ([]byte, error) { return marshalCleanYAML(i.Ingress) }
func (i *IngressWrapper) WithAnnotations(annotations map[string]string) *IngressWrapper {
	if i.Annotations == nil {
		i.Annotations = make(map[string]string)
	}
	maps.Copy(i.Annotations, annotations)
	return i
}
func (i *IngressWrapper) WithRules(rules []networkingv1.IngressRule) *IngressWrapper {
	i.Spec.Rules = rules
	return i
}
func (i *IngressWrapper) WithTLS(tls []networkingv1.IngressTLS) *IngressWrapper {
	i.Spec.TLS = tls
	return i
}
func (i *IngressWrapper) WithClass(class string) *IngressWrapper {
	if i.Spec.IngressClassName == nil {
		i.Spec.IngressClassName = new(string)
	}
	i.Spec.IngressClassName = &class
	return i
}

func (k *KustomizationWrapper) validate() error          { return nil }
func (k *KustomizationWrapper) marshal() ([]byte, error) { return yaml.Marshal(k) }
func (k *KustomizationWrapper) AddResources(resources []string) *KustomizationWrapper {
	k.Kustomization.Resources = append(k.Kustomization.Resources, resources...)
	return k
}
func (k *KustomizationWrapper) WithAnnotations(annotations map[string]string) *KustomizationWrapper {
	if k.Kustomization.MetaData == nil {
		k.Kustomization.MetaData = &types.ObjectMeta{}
	}
	if k.Kustomization.MetaData.Annotations == nil {
		k.Kustomization.MetaData.Annotations = make(map[string]string)
	}
	maps.Copy(k.Kustomization.MetaData.Annotations, annotations)
	return k
}

func (r *ReplicaSetWrapper) validate() error          { return nil }
func (r *ReplicaSetWrapper) marshal() ([]byte, error) { return marshalCleanYAML(r.ReplicaSet) }

func (s *ServiceWrapper) validate() error          { return nil }
func (s *ServiceWrapper) marshal() ([]byte, error) { return yaml.Marshal(s.Service) }
func (s *ServiceWrapper) WithSelector(labels map[string]string) *ServiceWrapper {
	s.Spec.Selector = labels
	return s
}
func (s *ServiceWrapper) WithPort(port, targetPort int32, protocol corev1.Protocol) *ServiceWrapper {
	if s.Spec.Ports == nil {
		s.Spec.Ports = []corev1.ServicePort{}
	}
	s.Spec.Ports = append(s.Spec.Ports, corev1.ServicePort{
		Port:       port,
		TargetPort: intstr.FromInt(int(targetPort)),
		Protocol:   protocol,
	})
	return s
}
func (s *ServiceWrapper) WithServiceType(serviceType corev1.ServiceType) *ServiceWrapper {
	s.Spec.Type = serviceType
	return s
}

// NESTED RESOURCES
func (c *ContainerWrapper) WithPort(port int32) *ContainerWrapper {
	if c.Ports == nil {
		c.Ports = []corev1.ContainerPort{}
	}
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
func (c *ContainerWrapper) WithEnvVar(name, value string) *ContainerWrapper {
	if c.Env == nil {
		c.Env = []corev1.EnvVar{}
	}
	c.Env = append(c.Env, corev1.EnvVar{
		Name:  name,
		Value: value,
	})
	return c
}
func (c *ContainerWrapper) WithImagePullPolicy(policy corev1.PullPolicy) *ContainerWrapper {
	c.ImagePullPolicy = policy
	return c
}

func (t *PodTemplateSpecWrapper) MergeLabels(labels map[string]string) *PodTemplateSpecWrapper {
	if t.ObjectMeta.Labels == nil {
		t.ObjectMeta.Labels = make(map[string]string)
	}
	maps.Copy(t.ObjectMeta.Labels, labels)
	return t
}
func (t *PodTemplateSpecWrapper) AddContainer(c *ContainerWrapper) *PodTemplateSpecWrapper {
	t.Spec.Containers = append(t.Spec.Containers, *c.Container)
	return t
}

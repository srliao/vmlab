package klusterhelper

import (
	"maps"

	appsv1 "k8s.io/api/apps/v1"
)

type DeploymentWrapper struct {
	*appsv1.Deployment
}

var _ KubeResource = &DeploymentWrapper{}
var _ Annotatable[*DeploymentWrapper] = &DeploymentWrapper{}

func (d *DeploymentWrapper) validate() error          { return nil }
func (d *DeploymentWrapper) marshal() ([]byte, error) { return marshalCleanYAML(d.Deployment) }
func (d *DeploymentWrapper) PodTemplate() *PodTemplateSpecWrapper {
	return &PodTemplateSpecWrapper{&d.Spec.Template}
}
func (d *DeploymentWrapper) MergeAnnotations(annotations map[string]string) *DeploymentWrapper {
	if d.Annotations == nil {
		d.Annotations = make(map[string]string)
	}
	maps.Copy(d.Annotations, annotations)
	return d
}
func (d *DeploymentWrapper) WithAnnotations(annotations map[string]string) *DeploymentWrapper {
	d.Annotations = annotations
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

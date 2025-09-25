package klusterhelper

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"sigs.k8s.io/yaml"
)

type PersistentVolumeClaimWrapper struct {
	*corev1.PersistentVolumeClaim
}

var _ KubeResource = &PersistentVolumeClaimWrapper{}

func (p *PersistentVolumeClaimWrapper) validate() error { return nil }
func (p *PersistentVolumeClaimWrapper) marshal() ([]byte, error) {
	return yaml.Marshal(p.PersistentVolumeClaim)
}

func (p *PersistentVolumeClaimWrapper) WithStorageRequest(storage string) *PersistentVolumeClaimWrapper {
	if p.Spec.Resources.Requests == nil {
		p.Spec.Resources.Requests = corev1.ResourceList{}
	}
	p.Spec.Resources.Requests[corev1.ResourceStorage] = resource.MustParse(storage)
	return p
}

func (p *PersistentVolumeClaimWrapper) WithStorageClass(class string) *PersistentVolumeClaimWrapper {
	p.Spec.StorageClassName = &class
	return p
}

func (p *PersistentVolumeClaimWrapper) WithAccessModes(modes []corev1.PersistentVolumeAccessMode) *PersistentVolumeClaimWrapper {
	p.Spec.AccessModes = modes
	return p
}

func (p *PersistentVolumeClaimWrapper) WithVolumeMode(mode corev1.PersistentVolumeMode) *PersistentVolumeClaimWrapper {
	p.Spec.VolumeMode = &mode
	return p
}

func (p *PersistentVolumeClaimWrapper) WithDataSourceRef(name, kind, apiGroup string) *PersistentVolumeClaimWrapper {
	p.Spec.DataSource = &corev1.TypedLocalObjectReference{
		APIGroup: &apiGroup,
		Kind:     kind,
		Name:     name,
	}
	return p
}

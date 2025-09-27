package klusterhelper

import (
	"path/filepath"

	fluxv1 "github.com/fluxcd/kustomize-controller/api/v1"
)

type FluxKustomizationWrapper struct {
	*fluxv1.Kustomization
	subpath string
}

var _ KubeResource = &FluxKustomizationWrapper{}

func (f *FluxKustomizationWrapper) marshal() ([]byte, error) {
	return marshalCleanYAML(f.Kustomization)
}
func (f *FluxKustomizationWrapper) validate() error { return nil }
func (f *FluxKustomizationWrapper) WithPath(base, subpath string) *FluxKustomizationWrapper {
	f.Kustomization.Spec.Path = filepath.Join(base, subpath)
	f.subpath = subpath
	return f
}
func (f *FluxKustomizationWrapper) WithDependsOn(deps []fluxv1.DependencyReference) *FluxKustomizationWrapper {
	f.Kustomization.Spec.DependsOn = deps
	return f
}

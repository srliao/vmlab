package klusterhelper

import (
	"fmt"

	"sigs.k8s.io/kustomize/api/types"
)

type App struct {
	objects   []KubeResource
	ks        *FluxKustomizationWrapper
	kust      *KustomizationWrapper
	name      string
	namespace string
	subfolder string
}

func NewApp(name, namespace, subfolder string) *App {
	a := &App{
		name:      name,
		namespace: namespace,
		subfolder: subfolder,
		kust: &KustomizationWrapper{
			Kustomization: &types.Kustomization{
				TypeMeta: types.TypeMeta{
					APIVersion: types.KustomizationVersion,
					Kind:       types.KustomizationKind,
				},
			},
		},
	}
	return a
}

func (a *App) SetKS(ks *FluxKustomizationWrapper) *App {
	a.ks = ks
	return a
}

func (a *App) SetKustomization(kust *KustomizationWrapper) *App {
	a.kust = kust
	return a
}

func (a *App) Kustomization() *KustomizationWrapper {
	return a.kust
}

func (a *App) AddObject(objects ...KubeResource) *App {
	a.objects = append(a.objects, objects...)
	return a
}

func (a *App) Validate() error {
	if len(a.objects) == 0 {
		return fmt.Errorf("no objects to validate")
	}
	return nil
}

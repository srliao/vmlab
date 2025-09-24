package klusterhelper

import (
	"fmt"

	"sigs.k8s.io/kustomize/api/types"
)

type App struct {
	resources []KubeResource
	files     []AppFile
	ks        *FluxKustomizationWrapper
	kust      *KustomizationWrapper
	name      string
	namespace string
	subfolder string
}

type AppFile interface {
	Name() string
	Content() []byte
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

func (a *App) AddObjects(objects ...KubeResource) *App {
	a.resources = append(a.resources, objects...)
	return a
}

func (a *App) AddFiles(files ...AppFile) *App {
	a.files = append(a.files, files...)
	return a
}

func (a *App) Validate() error {
	if len(a.resources) == 0 {
		return fmt.Errorf("no objects to validate")
	}
	return nil
}

type TextFile struct {
	name    string
	content []byte
}

func NewTextFile(name string, content []byte) *TextFile {
	return &TextFile{name: name, content: content}
}

func (t *TextFile) Name() string {
	return t.name
}

func (t *TextFile) Content() []byte {
	return t.content
}

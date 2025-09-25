package klusterhelper

import (
	"maps"

	"sigs.k8s.io/kustomize/api/types"
	"sigs.k8s.io/yaml"
)

type KustomizationWrapper struct {
	*types.Kustomization
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
	k.Kustomization.MetaData.Annotations = annotations
	return k
}

func (k *KustomizationWrapper) MergeAnnotations(annotations map[string]string) *KustomizationWrapper {
	if k.Kustomization.MetaData == nil {
		k.Kustomization.MetaData = &types.ObjectMeta{}
	}
	if k.Kustomization.MetaData.Annotations == nil {
		k.Kustomization.MetaData.Annotations = make(map[string]string)
	}
	maps.Copy(k.Kustomization.MetaData.Annotations, annotations)
	return k
}

func (k *KustomizationWrapper) AddConfigMapGeneratorFromFiles(name string, file ...string) *KustomizationWrapper {
	k.Kustomization.ConfigMapGenerator = append(k.Kustomization.ConfigMapGenerator, types.ConfigMapArgs{
		GeneratorArgs: types.GeneratorArgs{
			Name: name,
			KvPairSources: types.KvPairSources{
				FileSources: file,
			},
			Options: &types.GeneratorOptions{
				DisableNameSuffixHash: true,
				Annotations: map[string]string{
					"kustomize.toolkit.fluxcd.io/substitute": "disabled",
				},
			},
		},
	})
	return k
}

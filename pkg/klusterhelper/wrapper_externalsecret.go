package klusterhelper

import (
	esv1beta1 "github.com/external-secrets/external-secrets/apis/externalsecrets/v1beta1"
)

type ExternalSecretWrapper struct {
	*esv1beta1.ExternalSecret
}

var _ KubeResource = &ExternalSecretWrapper{}

func (e *ExternalSecretWrapper) validate() error          { return nil }
func (e *ExternalSecretWrapper) marshal() ([]byte, error) { return marshalCleanYAML(e.ExternalSecret) }
func (e *ExternalSecretWrapper) init() {
	if e.Spec.Target.Template == nil {
		e.Spec.Target.Template = &esv1beta1.ExternalSecretTemplate{
			EngineVersion: esv1beta1.TemplateEngineV2,
		}
	}
	if e.Spec.Target.Template.Data == nil {
		e.Spec.Target.Template.Data = make(map[string]string)
	}
}

func (e *ExternalSecretWrapper) AddDataToTemplate(name, from string) *ExternalSecretWrapper {
	e.init()
	e.Spec.Target.Template.Data[name] = from
	return e
}

func (e *ExternalSecretWrapper) AddMapDataToTemplate(d map[string]string) *ExternalSecretWrapper {
	e.init()
	for k, v := range d {
		e.Spec.Target.Template.Data[k] = v
	}
	return e
}

func (e *ExternalSecretWrapper) AddExternalDataFromKeyExtract(key string) *ExternalSecretWrapper {
	e.Spec.DataFrom = append(e.Spec.DataFrom, esv1beta1.ExternalSecretDataFromRemoteRef{
		Extract: &esv1beta1.ExternalSecretDataRemoteRef{
			Key: key,
		},
	})
	return e
}

package defaults

import (
	"path/filepath"

	esv1beta1 "github.com/external-secrets/external-secrets/apis/externalsecrets/v1beta1"
	fluxv1 "github.com/fluxcd/kustomize-controller/api/v1"
	"github.com/srliao/vmlab/pkg/klusterhelper"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func NewFluxKS(name, namespace, subpath string) *klusterhelper.FluxKustomizationWrapper {
	ks := &klusterhelper.FluxKustomizationWrapper{
		Kustomization: &fluxv1.Kustomization{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "kustomize.toolkit.fluxcd.io/v1",
				Kind:       "Kustomization",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      name,
				Namespace: namespace,
			},
			Spec: fluxv1.KustomizationSpec{
				TargetNamespace: namespace,
				CommonMetadata: &fluxv1.CommonMetadata{
					Labels: Labels(name),
				},
				SourceRef: srcRef,
				Interval:  metav1.Duration{Duration: defaultFluxReconcileInterval},
				Timeout:   &metav1.Duration{Duration: defaultFluxReconcileTimeout},
				Prune:     true,
			},
		},
	}

	ks.WithPath(filepath.Join(fluxKSBasePath, subpath), "app")
	return ks
}

func NewDeployment(name, namespace string) *klusterhelper.DeploymentWrapper {
	return &klusterhelper.DeploymentWrapper{
		Deployment: &appsv1.Deployment{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "apps/v1",
				Kind:       "Deployment",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      name,
				Namespace: namespace,
				Labels:    Labels(name),
			},
			Spec: appsv1.DeploymentSpec{
				Replicas: klusterhelper.Int32Ptr(1),
				Selector: &metav1.LabelSelector{
					MatchLabels: Labels(name),
				},
				Template: corev1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						// don't provide a name because we want k8s to auto-generate one
						Labels: Labels(name),
					},
					Spec: corev1.PodSpec{
						RestartPolicy: corev1.RestartPolicyAlways,
					},
				},
			},
		},
	}
}

func NewPodTemplate(name string) *klusterhelper.PodTemplateSpecWrapper {
	return &klusterhelper.PodTemplateSpecWrapper{
		PodTemplateSpec: &corev1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				// don't provide a name because we want k8s to auto-generate one
				Labels: Labels(name),
			},
			Spec: corev1.PodSpec{
				RestartPolicy: corev1.RestartPolicyAlways,
			},
		},
	}
}

func NewContainer(name string, imageSpec *klusterhelper.ImageSpec) *klusterhelper.ContainerWrapper {
	return &klusterhelper.ContainerWrapper{
		Container: &corev1.Container{
			Name:            name,
			Image:           imageSpec.String(),
			ImagePullPolicy: corev1.PullIfNotPresent,
		},
	}
}

func NewConfigMap(name, namespace string) *klusterhelper.ConfigMapWrapper {
	return &klusterhelper.ConfigMapWrapper{
		ConfigMap: &corev1.ConfigMap{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "v1",
				Kind:       "ConfigMap",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      name,
				Namespace: namespace,
			},
		},
	}
}

// we're defaulting to nginx
func NewIngress(name, namespace string) *klusterhelper.IngressWrapper {
	return &klusterhelper.IngressWrapper{
		Ingress: &networkingv1.Ingress{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "networking.k8s.io/v1",
				Kind:       "Ingress",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      name,
				Namespace: namespace,
				Labels:    Labels(name),
			},
		},
	}
}

func NewService(name, namespace string) *klusterhelper.ServiceWrapper {
	return &klusterhelper.ServiceWrapper{
		Service: &corev1.Service{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "v1",
				Kind:       "Service",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      name,
				Namespace: namespace,
				Labels:    Labels(name),
			},
			Spec: corev1.ServiceSpec{
				Selector: Labels(name),
			},
		},
	}
}

func NewDefaultIngressRules(host, service string, port int32) []networkingv1.IngressRule {
	pathType := networkingv1.PathTypePrefix
	return []networkingv1.IngressRule{
		{
			Host: host,
			IngressRuleValue: networkingv1.IngressRuleValue{
				HTTP: &networkingv1.HTTPIngressRuleValue{
					Paths: []networkingv1.HTTPIngressPath{
						{
							Path:     "/",
							PathType: &pathType,
							Backend: networkingv1.IngressBackend{
								Service: &networkingv1.IngressServiceBackend{
									Name: service,
									Port: networkingv1.ServiceBackendPort{
										Number: port,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func NewES(name, target string) *klusterhelper.ExternalSecretWrapper {
	return &klusterhelper.ExternalSecretWrapper{
		ExternalSecret: &esv1beta1.ExternalSecret{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "external-secrets.io/v1beta1",
				Kind:       "ExternalSecret",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name: name,
			},
			Spec: esv1beta1.ExternalSecretSpec{
				SecretStoreRef: esv1beta1.SecretStoreRef{
					Name: "onepassword-connect",
					Kind: "ClusterSecretStore",
				},
				Target: esv1beta1.ExternalSecretTarget{
					Name: target,
					Template: &esv1beta1.ExternalSecretTemplate{
						EngineVersion: esv1beta1.TemplateEngineV2,
					},
				},
			},
		},
	}
}

func NewESWithDataAndKey(name, target string, data map[string]string, keys ...string) *klusterhelper.ExternalSecretWrapper {
	es := &klusterhelper.ExternalSecretWrapper{
		ExternalSecret: &esv1beta1.ExternalSecret{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "external-secrets.io/v1beta1",
				Kind:       "ExternalSecret",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name: name,
			},
			Spec: esv1beta1.ExternalSecretSpec{
				SecretStoreRef: esv1beta1.SecretStoreRef{
					Name: "onepassword-connect",
					Kind: "ClusterSecretStore",
				},
				Target: esv1beta1.ExternalSecretTarget{
					Name: target,
					Template: &esv1beta1.ExternalSecretTemplate{
						EngineVersion: esv1beta1.TemplateEngineV2,
					},
				},
			},
		},
	}
	for k, v := range data {
		es.AddDataToTemplate(k, v)
	}
	for _, k := range keys {
		es.AddExternalDataFromKeyExtract(k)
	}
	return es
}

func NewPVC(claimName, namespace string, size string) *klusterhelper.PersistentVolumeClaimWrapper {
	filesystem := corev1.PersistentVolumeFilesystem
	return &klusterhelper.PersistentVolumeClaimWrapper{
		PersistentVolumeClaim: &corev1.PersistentVolumeClaim{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "v1",
				Kind:       "PersistentVolumeClaim",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      claimName,
				Namespace: namespace,
			},
			Spec: corev1.PersistentVolumeClaimSpec{
				AccessModes: []corev1.PersistentVolumeAccessMode{
					corev1.ReadWriteOnce,
				},
				StorageClassName: strPtr(defaultPVCStorageClass),
				VolumeMode:       &filesystem,
				Resources: corev1.VolumeResourceRequirements{
					Requests: corev1.ResourceList{
						corev1.ResourceStorage: resource.MustParse(size),
					},
				},
			},
		},
	}
}

func NewDefaultProbe() *corev1.Probe {
	return &corev1.Probe{
		ProbeHandler: corev1.ProbeHandler{
			HTTPGet: &corev1.HTTPGetAction{
				Path: "/healthcheck",
				Port: intstr.IntOrString{IntVal: 80},
			},
		},
		TimeoutSeconds: 1,
		PeriodSeconds:  10,
	}
}

func strPtr(s string) *string {
	return &s
}

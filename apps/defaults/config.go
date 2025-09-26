package defaults

import (
	"time"

	fluxv1 "github.com/fluxcd/kustomize-controller/api/v1"
	"github.com/srliao/vmlab/pkg/klusterhelper"
)

var srcRef = fluxv1.CrossNamespaceSourceReference{
	Kind: "GitRepository",
	Name: "home-kubernetes",
}

const (
	gitRepoBasePath = "vmlab/apps/main"
	fluxKSBasePath  = "./kubernetes/apps"

	defaultFluxReconcileInterval = time.Minute * 10
	defaultFluxReconcileTimeout  = time.Minute * 5

	defaultPVCStorageClass = "local-hostpath"
	defaultPVCSize         = "1Gi"
	defaultUID             = int64(568)
)

func Labels(name string) map[string]string {
	return map[string]string{
		"app":                          name,
		"app.kubernetes.io/name":       name,
		"app.kubernetes.io/instance":   name,
		"app.kubernetes.io/managed-by": "vmlab",
	}
}

func AddReloaderAnnotation[K any](a klusterhelper.Annotatable[K]) {
	a.MergeAnnotations(map[string]string{
		"reloader.stakater.com/auto": "true",
	})
}

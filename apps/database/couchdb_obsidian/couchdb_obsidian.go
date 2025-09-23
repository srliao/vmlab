package couchdb_obsidian

import (
	"log"

	"github.com/srliao/vmlab/apps/defaults"
	"github.com/srliao/vmlab/pkg/klusterhelper"
)

const (
	name        = "couchdb-obsidian"
	namespace   = "database"
	component   = "database"
	servicePort = 5984
)

var imageSpec = &klusterhelper.ImageSpec{
	Name: "couchdb",
	Tag:  "sha256:2f27666c837cb9096c95e3cbd1c125651c0fa6a92594ffc9268671e4d065d6a9",
}

func Chart() *klusterhelper.App {
	subpath, err := defaults.RelDirPath(1)
	if err != nil {
		log.Panicf("getting subpath failed for %v/%v: %v", name, namespace, err)
	}
	c := klusterhelper.NewApp(name, namespace, subpath).
		SetKS(defaults.NewFluxKS(name, namespace, subpath)).
		AddObject(deployment()).
		AddObject(service()).
		AddObject(ingress())
	defaults.AddReloaderAnnotation(c.Kustomization())
	return c
}

func deployment() klusterhelper.KubeResource {
	labels := map[string]string{
		"app.kubernetes.io/component": component,
	}

	container := defaults.NewContainer(name, imageSpec).
		WithPort(servicePort).
		WithCPURequest("20m").
		WithMemoryLimit("2Gi")

	pod := defaults.NewPodTemplate(name).
		MergeLabels(labels).
		AddContainer(container)

	dep := defaults.NewDeployment(name, namespace).
		MergeLabels(labels).
		WithPodTemplate(pod)

	defaults.AddReloaderAnnotation(dep)

	return dep
}

func service() klusterhelper.KubeResource {
	s := defaults.NewService(name, namespace).
		WithPort(5984, 5984, "TCP").
		WithServiceType("ClusterIP")
	return s
}

func ingress() klusterhelper.KubeResource {
	s := defaults.NewIngress(name, namespace).
		WithClass("internal").
		WithRules(defaults.NewDefaultIngressRules("couchdb-obsidian.winterspring.ca", name, servicePort))
	return s
}

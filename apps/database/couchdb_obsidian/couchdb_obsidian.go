package couchdb_obsidian

import (
	"log"

	_ "embed"

	"github.com/srliao/vmlab/apps/defaults"
	"github.com/srliao/vmlab/pkg/klusterhelper"
)

const (
	name               = "couchdb-obsidian"
	namespace          = "database"
	component          = "database"
	servicePort        = 5984
	configIniMapName   = "obsidian-couchdb-configmap"
	defaultLoginSecret = "couchdb-obsidian-secrets"
	dataPVCClaimName   = "couchdb-obsidian"
)

var imageSpec = &klusterhelper.ImageSpec{
	Name: "couchdb",
	Tag:  "3.5.0",
}

//go:embed config.ini
var configInit []byte

func Chart() *klusterhelper.App {
	subpath, err := defaults.RelDirPath(1)
	if err != nil {
		log.Panicf("getting subpath failed for %v/%v: %v", name, namespace, err)
	}

	c := klusterhelper.
		NewApp(name, namespace, subpath).
		SetKS(defaults.NewFluxKS(name, namespace, subpath)).
		AddObjects(
			deployment(),
			service(),
			ingress(),
			secrets(),
		).
		AddObjects(
			defaults.NewVolsyncPVCResources(name, namespace, "5Gi")...,
		).
		AddFiles(klusterhelper.NewTextFile("config.ini", configInit))

	// make sure kustomization gets reloaded if config ini changes
	defaults.AddReloaderAnnotation(c.Kustomization())

	// add config map for init file
	c.Kustomization().AddConfigMapGeneratorFromFiles(configIniMapName, "config.ini")

	return c
}

func deployment() klusterhelper.KubeResource {
	const dataPVC = "data"
	const configStorage = "config-storage"

	// init container to copy config.init file to default.d
	init := defaults.
		NewContainer(
			"init-config", &klusterhelper.ImageSpec{
				Name: "busybox",
				Tag:  "latest",
			}).
		AddCommands(
			"/bin/sh",
			"-c",
			"cp /tmp/config/*.ini /opt/couchdb/etc/default.d/; ls -lrt /opt/couchdb/etc/default.d;",
		).
		MountVolume(configIniMapName, "/tmp/config").
		MountVolume(configStorage, "/opt/couchdb/etc/default.d").
		MountVolume(dataPVC, "/opt/couchdb/data")

	// app container
	container := defaults.
		NewContainer(name, imageSpec).
		WithPort(servicePort).
		WithCPURequest("20m").
		WithMemoryRequest("512Mi").
		WithMemoryLimit("1024Mi").
		MountVolume(configStorage, "/opt/couchdb/etc/default.d").
		MountVolume(dataPVC, "/opt/couchdb/data").
		AddEnvFromSecret(defaultLoginSecret).
		WithLivenessProbe(defaults.NewDefaultProbe()).
		WithReadinessProbe(defaults.NewDefaultProbe())

	deploy := defaults.
		NewDeployment(name, namespace)

	deploy.PodTemplate().
		AddVolumes(
			defaults.NewConfigMapVolume(configIniMapName),
			defaults.NewEmptyDirVolume(configStorage),
			defaults.NewPVCVolume(dataPVC, dataPVCClaimName),
		).
		AddContainer(container).
		AddInitContainer(init)

	defaults.AddReloaderAnnotation(deploy)

	return deploy
}

func service() klusterhelper.KubeResource {
	s := defaults.NewService(name, namespace).
		AddPort(servicePort, servicePort, "TCP").
		WithServiceType("ClusterIP")
	return s
}

func ingress() klusterhelper.KubeResource {
	s := defaults.NewIngress(name, namespace).
		WithClass("internal").
		WithRules(
			defaults.NewDefaultIngressRules(
				"couchdb-obsidian.winterspring.ca",
				name,
				servicePort,
			))
	return s
}

func secrets() klusterhelper.KubeResource {
	return defaults.NewESWithDataAndKey(
		name,
		defaultLoginSecret,
		map[string]string{
			"COUCHDB_USER":     "{{ .OBSIDIAN_COUCHDB_USER }}",
			"COUCHDB_PASSWORD": "{{ .OBSIDIAN_COUCHDB_PASSWORD }}",
		},
		"obsidian",
	)
}

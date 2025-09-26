package couchdb_obsidian

import (
	"log"

	_ "embed"

	fluxmeta "github.com/fluxcd/pkg/apis/meta"
	"github.com/srliao/vmlab/apps/defaults"
	"github.com/srliao/vmlab/pkg/klusterhelper"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

const (
	name               = "couchdb-obsidian"
	namespace          = "database"
	servicePort        = 5984
	configIniMapName   = "obsidian-couchdb-configmap"
	defaultLoginSecret = "couchdb-obsidian-secrets"
	dataPVCClaimName   = "couchdb-obsidian"
)

var imageSpec = &klusterhelper.ImageSpec{
	Repository: "docker.io/library/couchdb",
	Tag:        "3.5.0",
}

//go:embed config.ini
var configInit []byte

func Chart() *klusterhelper.App {
	subpath, err := defaults.RelDirPath(1)
	if err != nil {
		log.Panicf("getting subpath failed for %v/%v: %v", name, namespace, err)
	}

	ks := defaults.NewFluxKS(name, namespace, subpath)
	ks.WithDependsOn([]fluxmeta.NamespacedObjectReference{
		{
			Name: "external-secrets-stores",
		},
		{
			Name: "volsync",
		},
		{
			Name: "democratic-csi-local-path",
		},
	})

	c := klusterhelper.
		NewApp(name, namespace, subpath).
		SetKS(ks).
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
				Repository: "docker.io/library/busybox",
				Tag:        "latest",
			}).
		AddCommands(
			"/bin/sh",
			"-c",
			"cp /tmp/config/*.ini /opt/couchdb/etc/default.d/; ls -lrt /opt/couchdb/etc/default.d;",
		).
		AddVolumeMount(configIniMapName, "/tmp/config").
		AddVolumeMount(configStorage, "/opt/couchdb/etc/default.d").
		AddVolumeMount(dataPVC, "/opt/couchdb/data")

	// app container
	container := defaults.
		NewContainer(name, imageSpec).
		AddPort(servicePort).
		WithCPURequest("20m").
		WithMemoryRequest("512Mi").
		WithMemoryLimit("1024Mi").
		AddVolumeMount(configStorage, "/opt/couchdb/etc/default.d").
		AddVolumeMount(dataPVC, "/opt/couchdb/data").
		AddEnvFromSecret(defaultLoginSecret).
		WithLivenessProbe(probe()).
		WithReadinessProbe(probe())

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

func probe() *corev1.Probe {
	return &corev1.Probe{
		ProbeHandler: corev1.ProbeHandler{
			HTTPGet: &corev1.HTTPGetAction{
				Path: "/_up",
				Port: intstr.IntOrString{IntVal: servicePort},
			},
		},
		TimeoutSeconds: 1,
		PeriodSeconds:  10,
	}
}

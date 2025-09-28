package kqmc_checker_bot

import (
	"log"

	_ "embed"

	fluxv1 "github.com/fluxcd/kustomize-controller/api/v1"
	"github.com/srliao/vmlab/apps/defaults"
	"github.com/srliao/vmlab/pkg/klusterhelper"
	corev1 "k8s.io/api/core/v1"
)

const (
	name          = "kqmc-checker-bot"
	namespace     = "gcsim"
	discordSecret = "kqmc-checker-bot-secrets"
)

var imageSpec = &klusterhelper.ImageSpec{
	Repository: "ghcr.io/charlie-zheng/gcsim-kqmc-checker-develop",
	Tag:        "rolling@sha256:9a37d56de0f34e1afd560a9d35566283d65c814a5f8f48481b2094ce71f9e6c7",
}

func Chart() *klusterhelper.App {
	subpath, err := defaults.RelDirPath(1)
	if err != nil {
		log.Panicf("getting subpath failed for %v/%v: %v", name, namespace, err)
	}

	ks := defaults.NewFluxKS(name, namespace, subpath)
	ks.WithDependsOn([]fluxv1.DependencyReference{
		{
			Name: "external-secrets-stores",
		},
	})

	c := klusterhelper.
		NewApp(name, namespace, subpath).
		SetKS(ks).
		AddObjects(
			deployment(),
			secrets(),
		)

	return c
}

func deployment() klusterhelper.KubeResource {

	// app container
	container := defaults.
		NewContainer(name, imageSpec).
		WithCPURequest("20m").
		WithMemoryRequest("512Mi").
		WithMemoryLimit("1024Mi").
		AddEnvFromSecret(discordSecret).
		WithImagePullPolicy(corev1.PullAlways)

	deploy := defaults.
		NewDeployment(name, namespace)

	deploy.PodTemplate().
		AddContainer(container)

	return deploy
}

func secrets() klusterhelper.KubeResource {
	return defaults.NewESWithDataAndKey(
		name,
		discordSecret,
		map[string]string{
			"DISCORD_TOKEN": "{{ .KQMC_CHECKER_BOT_DISCORD_TOKEN }}",
		},
		"gcsim",
	)
}

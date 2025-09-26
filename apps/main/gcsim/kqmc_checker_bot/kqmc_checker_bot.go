package kqmc_checker_bot

import (
	"log"

	_ "embed"

	fluxmeta "github.com/fluxcd/pkg/apis/meta"
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
	Tag:        "rolling@sha256:0ef7bf7890a1fc99e7c88ddb1fd63947c8dff67ee1993888a5077bf9c0bcde76",
}

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

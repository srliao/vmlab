package main

import (
	"github.com/srliao/vmlab/apps/main/database/couchdb_obsidian"
	"github.com/srliao/vmlab/apps/main/gcsim/kqmc_checker_bot"
	"github.com/srliao/vmlab/pkg/klusterhelper"
)

func main() {
	klusterhelper.NewBuilder().
		AddApp(couchdb_obsidian.Chart()).
		AddApp(kqmc_checker_bot.Chart()).
		Validate().
		Build("../kubernetes/apps", true)
}

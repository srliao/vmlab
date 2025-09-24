package main

import (
	"github.com/srliao/vmlab/apps/database/couchdb_obsidian"
	"github.com/srliao/vmlab/pkg/klusterhelper"
)

func main() {
	klusterhelper.NewBuilder().
		AddApp(couchdb_obsidian.Chart()).
		Validate().
		Build("../kubernetes/apps", true)
}

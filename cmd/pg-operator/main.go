package main

import (
	"log"

	logs "github.com/appscode/go/log/golog"
	"github.com/tekliner/postgres/pkg/cmds"
)

func main() {
	logs.InitLogs()
	defer logs.FlushLogs()
	if err := cmds.NewRootCmd(Version).Execute(); err != nil {
		log.Fatal(err)
	}
}

package main

import (
	"github.com/backup-blob/zfs-backup-blob/cmd/command"
	"os"
)

var version string

func main() {
	command.RootCmd.Version = version

	err := command.RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

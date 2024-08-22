package main

import (
	"fmt"
	"github.com/backup-blob/zfs-backup-blob/cmd/command"
	"os"
)

func main() {
	err := GenMarkdownTree(command.RootCmd, "../astro/src/content/docs/cli/")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

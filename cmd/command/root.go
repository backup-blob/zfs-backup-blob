package command

import (
	"github.com/backup-blob/zfs-backup-blob/cmd/command/backup"
	"github.com/backup-blob/zfs-backup-blob/cmd/command/group"
	"github.com/backup-blob/zfs-backup-blob/cmd/command/shared"
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "zfs-backup-blob",
	Short: "_",
}

func init() {
	RootCmd.AddCommand(group.GroupCmd)
	RootCmd.AddCommand(backup.BackupCmd)
	RootCmd.PersistentFlags().StringVarP(&shared.ConfigPath, "configPath", "c", "~/.bbackup.yaml", "Path to the config file")
	RootCmd.PersistentFlags().StringVarP(&shared.LogLevel, "logLevel", "l", "disabled", "LogLevel = debug|disabled")
}

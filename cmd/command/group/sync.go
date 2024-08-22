package group

import (
	"github.com/backup-blob/zfs-backup-blob/cmd/command/shared"
	"github.com/backup-blob/zfs-backup-blob/internal/domain"
	"github.com/spf13/cobra"
)

var backupSyncGroupName string

var backupSyncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync snapshots of a group to remote",
	RunE: func(cmd *cobra.Command, args []string) error {
		var backupSync domain.BackupSyncUsecase

		c := shared.LoadDeps(shared.ConfigPath, shared.LogLevel)
		err := c.Resolve(&backupSync)
		if err != nil {
			return err
		}

		return backupSync.Backup(cmd.Context(), backupSyncGroupName)
	},
}

func init() {
	GroupCmd.AddCommand(backupSyncCmd)

	backupSyncCmd.Flags().StringVarP(&backupSyncGroupName, "group", "g", "default", "Volume group name")
}

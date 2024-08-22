package backup

import (
	"github.com/backup-blob/zfs-backup-blob/cmd/command/shared"
	"github.com/backup-blob/zfs-backup-blob/internal/domain"
	"github.com/spf13/cobra"
)

var restoreConf domain.RestoreParams

var restoreCmd = &cobra.Command{
	Use:   "restore",
	Short: "Restore a backup",
	RunE: func(cmd *cobra.Command, args []string) error {
		var backupUse domain.BackupUsecase
		c := shared.LoadDeps(shared.ConfigPath, shared.LogLevel)
		err := c.Resolve(&backupUse)
		if err != nil {
			return err
		}

		return backupUse.Restore(cmd.Context(), &restoreConf)
	},
}

func init() {
	BackupCmd.AddCommand(restoreCmd)
	restoreCmd.Flags().StringVarP(&restoreConf.BlobKey, "blob-key", "b", "", "S3 Key to the backup to restore (excluding prefix)")
	restoreCmd.Flags().StringVarP(&restoreConf.TargetZfsLocation, "target", "t", "", "Path to zfs pool/<volume/filesystem>")
	restoreCmd.Flags().BoolVarP(&restoreConf.RestoreAll, "restoreAll", "r", true, "Restores all incremental snapshots including the full backup")
	restoreCmd.MarkFlagRequired("blob-key") //nolint:errcheck //no need
	restoreCmd.MarkFlagRequired("target")   //nolint:errcheck //no need
}

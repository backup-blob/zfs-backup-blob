package backup

import (
	"github.com/backup-blob/zfs-backup-blob/cmd/command/shared"
	"github.com/backup-blob/zfs-backup-blob/internal/domain"
	"github.com/spf13/cobra"
)

var backupListVolumeName string

var backupListCmd = &cobra.Command{
	Use:   "list",
	Short: "List the backups stored on the remote storage for a volume/fs",
	RunE: func(cmd *cobra.Command, args []string) error {
		var listUse domain.BackupListUsecase

		c := shared.LoadDeps(shared.ConfigPath, shared.LogLevel)
		err := c.Resolve(&listUse)
		if err != nil {
			return err
		}

		return listUse.List(cmd.Context(), backupListVolumeName, cmd.OutOrStdout())
	},
}

func init() {
	BackupCmd.AddCommand(backupListCmd)

	backupListCmd.Flags().StringVarP(&backupListVolumeName, "volume", "", "", "Volume name (Example: pool/vol1)")
	backupListCmd.MarkFlagRequired("volume") //nolint:errcheck // no need
}

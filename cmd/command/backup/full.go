package backup

import (
	"github.com/backup-blob/zfs-backup-blob/cmd/command/shared"
	"github.com/backup-blob/zfs-backup-blob/internal/domain"
	"github.com/spf13/cobra"
)

type BackupFullConf struct {
	Snapshot string
}

var backupFullConf BackupFullConf

var BackupFullCmd = &cobra.Command{
	Use:   "full",
	Short: "Create a full backup of a zfs snapshot",
	RunE: func(cmd *cobra.Command, args []string) error {
		var backupUse domain.BackupUsecase
		c := shared.LoadDeps(shared.ConfigPath, shared.LogLevel)
		err := c.Resolve(&backupUse)
		if err != nil {
			return err
		}

		return backupUse.BackupFull(
			cmd.Context(),
			backupFullConf.Snapshot,
			false,
		)
	},
}

func init() {
	BackupCmd.AddCommand(BackupFullCmd)

	BackupFullCmd.Flags().StringVarP(&backupFullConf.Snapshot, "snap", "s", "", "Snapshot name (Example: pool/vol@snapshot1)")
	BackupFullCmd.MarkFlagRequired("snap") //nolint:errcheck // no need
}

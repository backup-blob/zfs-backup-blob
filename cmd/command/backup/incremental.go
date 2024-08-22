package backup

import (
	"github.com/backup-blob/zfs-backup-blob/cmd/command/shared"
	"github.com/backup-blob/zfs-backup-blob/internal/domain"
	"github.com/spf13/cobra"
)

type BackupIncConf struct {
	SnapshotNew  string
	SnapshotBase string
}

var backupIncConf BackupIncConf

var BackupIncCmd = &cobra.Command{
	Use:   "incremental",
	Short: "Create a incremental backup of a zfs snapshot",
	RunE: func(cmd *cobra.Command, args []string) error {
		var backupUse domain.BackupUsecase
		c := shared.LoadDeps(shared.ConfigPath, shared.LogLevel)
		err := c.Resolve(&backupUse)
		if err != nil {
			return err
		}

		return backupUse.BackupIncremental(
			cmd.Context(),
			backupIncConf.SnapshotBase,
			backupIncConf.SnapshotNew,
			false,
		)
	},
}

func init() {
	BackupCmd.AddCommand(BackupIncCmd)

	BackupIncCmd.Flags().StringVarP(&backupIncConf.SnapshotNew, "snap", "s", "", "Snapshot name (Example: pool/vol@snapshot1)")
	BackupIncCmd.Flags().StringVarP(&backupIncConf.SnapshotBase, "base", "b", "", "The base snapshot name to base this increment on (Example: pool/vol@snapshot0)")
	BackupIncCmd.MarkFlagRequired("snap") //nolint:errcheck // no need
	BackupIncCmd.MarkFlagRequired("base") //nolint:errcheck // no need
}

package group

import (
	"github.com/backup-blob/zfs-backup-blob/cmd/command/shared"
	"github.com/backup-blob/zfs-backup-blob/internal/domain"
	"github.com/spf13/cobra"
)

type snapshotConf struct {
	groupName  string
	backupType string
}

var snapConf snapshotConf

var snapshotCmd = &cobra.Command{
	Use:   "snapshot",
	Short: "Create snapshots for all zfs volume/fs belonging to a group",
	RunE: func(cmd *cobra.Command, args []string) error {
		var snapUse domain.SnapshotUsecase
		c := shared.LoadDeps(shared.ConfigPath, shared.LogLevel)
		err := c.Resolve(&snapUse)
		if err != nil {
			return err
		}

		return snapUse.CreateByGroup(snapConf.groupName, snapConf.backupType)
	},
}

func init() {
	GroupCmd.AddCommand(snapshotCmd)

	snapshotCmd.Flags().StringVarP(&snapConf.groupName, "group", "g", "default", "Volume group name")
	snapshotCmd.Flags().StringVarP(&snapConf.backupType, "type", "t", "", "Backup type (full|incremental)")
	snapshotCmd.MarkFlagRequired("type") //nolint:errcheck // no need
}

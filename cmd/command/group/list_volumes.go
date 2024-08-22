package group

import (
	"github.com/backup-blob/zfs-backup-blob/cmd/command/shared"
	"github.com/backup-blob/zfs-backup-blob/internal/domain"
	"github.com/spf13/cobra"
)

var listVolumesCmd = &cobra.Command{
	Use:   "list-volumes",
	Short: "List all volumes which belong to a group",
	RunE: func(cmd *cobra.Command, args []string) error {
		var volumeUsecase domain.VolumeUsecase

		c := shared.LoadDeps(shared.ConfigPath, shared.LogLevel)
		err := c.Resolve(&volumeUsecase)
		if err != nil {
			return err
		}

		return volumeUsecase.List(cmd.OutOrStdout())
	},
}

func init() {
	GroupCmd.AddCommand(listVolumesCmd)
}

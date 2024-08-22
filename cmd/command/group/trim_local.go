package group

import (
	"github.com/backup-blob/zfs-backup-blob/cmd/command/shared"
	"github.com/backup-blob/zfs-backup-blob/internal/domain"
	"github.com/spf13/cobra"
)

var trimLocalParams = domain.TrimLocalParameters{GroupName: "default"}

var trimLocalCmd = &cobra.Command{
	Use:   "trim-local",
	Short: "Trim local snapshots of a group",
	RunE: func(cmd *cobra.Command, args []string) error {
		var trimUsecase domain.TrimUsecase

		c := shared.LoadDeps(shared.ConfigPath, shared.LogLevel)
		err := c.Resolve(&trimUsecase)
		if err != nil {
			return err
		}

		return trimUsecase.TrimLocal(cmd.Context(), &trimLocalParams)
	},
}

func init() {
	GroupCmd.AddCommand(trimLocalCmd)

	trimLocalCmd.Flags().StringVarP(&trimLocalParams.GroupName, "group", "g", "default", "Volume group name")
	trimLocalCmd.Flags().BoolVarP(&trimLocalParams.DryRun, "dry-run", "d", false, "When set to true, snapshots are not deleted")

	trimLocalCmd.MarkFlagRequired("group")
}

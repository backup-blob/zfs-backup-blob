package group

import (
	"github.com/backup-blob/zfs-backup-blob/cmd/command/shared"
	"github.com/backup-blob/zfs-backup-blob/internal/domain"
	"github.com/spf13/cobra"
)

var trimRemoteParams = domain.TrimRemoteParameters{GroupName: "default"}

var trimRemoteCmd = &cobra.Command{
	Use:   "trim-remote",
	Short: "Trim remote backups of a group",
	RunE: func(cmd *cobra.Command, args []string) error {
		var trimUsecase domain.TrimUsecase

		c := shared.LoadDeps(shared.ConfigPath, shared.LogLevel)
		err := c.Resolve(&trimUsecase)
		if err != nil {
			return err
		}

		return trimUsecase.TrimRemote(cmd.Context(), &trimRemoteParams)
	},
}

func init() {
	GroupCmd.AddCommand(trimRemoteCmd)

	trimRemoteCmd.Flags().StringVarP(&trimRemoteParams.GroupName, "group", "g", "default", "Volume group name")
	trimRemoteCmd.Flags().BoolVarP(&trimRemoteParams.DryRun, "dry-run", "d", false, "When set to true, blobs are not deleted")

	trimRemoteCmd.MarkFlagRequired("group")
}

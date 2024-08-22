package group

import (
	"github.com/backup-blob/zfs-backup-blob/cmd/command/shared"
	"github.com/backup-blob/zfs-backup-blob/internal/domain"
	"github.com/spf13/cobra"
)

var groupAddVolumeName string
var groupAddVolumeGroup string

var groupAddVolume = &cobra.Command{
	Use:   "add-volume",
	Short: "Add a volume to a group",
	RunE: func(cmd *cobra.Command, args []string) error {
		var volumeUsecase domain.VolumeUsecase

		c := shared.LoadDeps(shared.ConfigPath, shared.LogLevel)
		err := c.Resolve(&volumeUsecase)
		if err != nil {
			return err
		}

		return volumeUsecase.AddToGroup(groupAddVolumeName, groupAddVolumeGroup)
	},
}

func init() {
	GroupCmd.AddCommand(groupAddVolume)

	groupAddVolume.Flags().StringVarP(&groupAddVolumeGroup, "group", "g", "default", "Volume group name")
	groupAddVolume.Flags().StringVarP(&groupAddVolumeName, "volume", "", "", "Volume name (Example: pool/vol1)")
	groupAddVolume.MarkFlagRequired("volume")
}

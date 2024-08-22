package backup

import "github.com/spf13/cobra"

var BackupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Actions related to backups",
}

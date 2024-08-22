package domain

import "io"

type RenderDriver interface {
	RenderTable(writer io.Writer, headerRow []interface{}, rows [][]interface{})
}

type RenderRepository interface {
	RenderBackupTable(writer io.Writer, backups []BackupRecordWithKey)
	RenderVolumeTable(writer io.Writer, volumes []*ZfsVolume)
}

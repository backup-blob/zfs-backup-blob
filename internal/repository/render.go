package repository

import (
	"github.com/backup-blob/zfs-backup-blob/internal/domain"
	"io"
)

type Sizer func(int642 *int64) string

type render struct {
	renderDriver domain.RenderDriver
	sizer        Sizer
}

func NewRender(renderDriver domain.RenderDriver, sizer Sizer) domain.RenderRepository {
	return &render{renderDriver: renderDriver, sizer: sizer}
}

func (r *render) RenderBackupTable(writer io.Writer, backups []domain.BackupRecordWithKey) {
	headerRow := []interface{}{"Key", "Type", "Size"}
	rows := make([][]interface{}, len(backups))

	for i, b := range backups {
		rows[i] = []interface{}{b.Key, b.Type.String(), r.sizer(b.Size)}
	}

	r.renderDriver.RenderTable(writer, headerRow, rows)
}

func (r *render) RenderVolumeTable(writer io.Writer, volumes []*domain.ZfsVolume) {
	headerRow := []interface{}{"Volume", "Group"}
	rows := make([][]interface{}, len(volumes))

	for i, b := range volumes {
		rows[i] = []interface{}{b.Name, b.GroupName}
	}

	r.renderDriver.RenderTable(writer, headerRow, rows)
}

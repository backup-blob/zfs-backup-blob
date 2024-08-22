package driver

import (
	"github.com/backup-blob/zfs-backup-blob/internal/domain"
	"github.com/jedib0t/go-pretty/v6/table"
	"io"
)

type renderDriver struct {
}

func NewRender() domain.RenderDriver {
	return &renderDriver{}
}

func (td *renderDriver) RenderTable(writer io.Writer, headerRow []interface{}, rows [][]interface{}) {
	t := table.NewWriter()
	t.SetOutputMirror(writer)
	t.AppendHeader(headerRow)

	for _, i := range rows {
		t.AppendRow(i)
	}

	t.Render()
}

package usecase

import (
	"context"
	"fmt"
	"github.com/backup-blob/zfs-backup-blob/internal/domain"
	"io"
	"path"
	"sort"
)

type backupList struct {
	stateRepo      domain.BackupStateRepo
	namingStrategy domain.SnapshotNamestrategy
	renderRepo     domain.RenderRepository
}

func NewBackupList(
	stateRepo domain.BackupStateRepo,
	namingStrategy domain.SnapshotNamestrategy,
	renderRepo domain.RenderRepository,
) domain.BackupListUsecase {
	return &backupList{ //nolint:whitespace // no need
		stateRepo:      stateRepo,
		namingStrategy: namingStrategy,
		renderRepo:     renderRepo,
	}
}

func (b *backupList) List(ctx context.Context, volumeName string, writer io.Writer) error {
	stateFilePath := fmt.Sprintf("%s/%s", volumeName, domain.StateFileName)

	state, err := b.stateRepo.Download(ctx, stateFilePath)
	if err != nil {
		return fmt.Errorf("stateRepo.Download failed %w", err)
	}

	backupPaths := make([]string, 0, len(state.Backups))
	for k := range state.Backups {
		backupPaths = append(backupPaths, k)
	}

	sort.Slice(backupPaths, func(i, j int) bool {
		return b.namingStrategy.IsGreater(path.Base(backupPaths[i]), path.Base(backupPaths[j]))
	})

	backups := make([]domain.BackupRecordWithKey, 0, len(backupPaths))
	for _, i := range backupPaths {
		backups = append(backups, domain.BackupRecordWithKey{BackupRecord: state.Backups[i], Key: i})
	}

	b.renderRepo.RenderBackupTable(writer, backups)

	return nil
}

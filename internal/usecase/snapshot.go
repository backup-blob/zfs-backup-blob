package usecase

import (
	"fmt"
	"github.com/backup-blob/zfs-backup-blob/internal/domain"
)

type snapshot struct {
	snapRepo   domain.SnapshotRepository
	volumeRepo domain.VolumeRepository
	logger     domain.LogRepository
}

func NewSnapshot(
	snapRepo domain.SnapshotRepository,
	volumeRepo domain.VolumeRepository,
	logger domain.LogRepository,
) domain.SnapshotUsecase {
	return &snapshot{ //nolint:whitespace // no need
		snapRepo:   snapRepo,
		volumeRepo: volumeRepo,
		logger:     logger,
	}
}

// todo: move to volume usecase.

func (s *snapshot) CreateByVolume(volume *domain.ZfsVolume, backupType domain.BackupType) error {
	snap, errC := s.snapRepo.CreateWithType(volume, backupType)

	s.logger.Infof("created snapshot %s for volume %s with backuptype %s", snap, volume.Name, backupType)

	if errC != nil {
		return fmt.Errorf("creating snapshot failed for volume=%s, %w", volume.Name, errC)
	}

	return nil
}

func (s *snapshot) CreateByGroup(groupName, backupType string) error {
	if groupName == "" {
		groupName = domain.DefaultVolumeGroup
	}

	volumes, err := s.volumeRepo.ListVolumesByGroup(groupName)
	if err != nil {
		return fmt.Errorf("volumeRepo.ListVolumesByGroup failed %w", err)
	}

	if len(volumes) == 0 {
		return fmt.Errorf("no volumes matched")
	}

	enumBackupType := domain.StringToBackupType(backupType)
	if enumBackupType == domain.Unknown {
		return fmt.Errorf("backupType %s unknown", backupType)
	}

	for _, volume := range volumes {
		errC := s.CreateByVolume(volume, enumBackupType)
		if errC != nil {
			return fmt.Errorf("create one snapshot failed for volume=%s, %w", volume.Name, errC)
		}
	}

	return nil
}

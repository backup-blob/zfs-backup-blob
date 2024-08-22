package repository

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/backup-blob/zfs-backup-blob/internal/domain"
	"gopkg.in/yaml.v3"
)

type backupStateRepo struct {
	sd domain.StorageDriver
}

func NewBackupStateRepo(sd domain.StorageDriver) domain.BackupStateRepo {
	return &backupStateRepo{sd: sd}
}

func (b *backupStateRepo) Download(ctx context.Context, key string) (*domain.BackupState, error) {
	var backupState domain.BackupState

	buf := bytes.NewBuffer([]byte{})

	err := b.sd.Download(ctx, &domain.DownloadParameters{Key: key}, buf)
	if err != nil {
		return nil, fmt.Errorf("storageDriver.Download failed: %w", err)
	}

	errU := yaml.Unmarshal(buf.Bytes(), &backupState)
	if errU != nil {
		return nil, fmt.Errorf("yaml.Unmarshal of state failed: %w", errU)
	}

	return &backupState, nil
}

func (b *backupStateRepo) DownloadOrDefault(ctx context.Context, key string) (*domain.BackupState, error) {
	res, err := b.Download(ctx, key)
	if err != nil && errors.Is(err, domain.ErrNotFound) {
		return &domain.BackupState{Backups: map[string]domain.BackupRecord{}}, nil
	}

	return res, err
}

func (b *backupStateRepo) UpdateState(ctx context.Context, stateFileKey string, f func(state *domain.BackupState) error) error {
	state, err := b.DownloadOrDefault(ctx, stateFileKey)
	if err != nil {
		return fmt.Errorf("b.DownloadOrDefault failed: %w", err)
	}

	errF := f(state)
	if errF != nil {
		return fmt.Errorf("f failed %w", errF)
	}

	errU := b.Upload(ctx, stateFileKey, state)
	if errU != nil {
		return fmt.Errorf("b.Upload failed: %w", errU)
	}

	return nil
}

func (b *backupStateRepo) Upload(ctx context.Context, key string, state *domain.BackupState) error {
	out, err := yaml.Marshal(state)
	if err != nil {
		return fmt.Errorf("yaml.Marshal of state failed: %w", err)
	}

	_, errU := b.sd.Upload(ctx, &domain.UploadParameters{Key: key}, bytes.NewReader(out))
	if errU != nil {
		return fmt.Errorf("sd.Upload failed: %w", errU)
	}

	return nil
}

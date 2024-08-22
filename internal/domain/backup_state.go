package domain

import (
	"context"
	"fmt"
	"path"
)

type BackupState struct {
	Head    string                  `yaml:"head"`
	Backups map[string]BackupRecord `yaml:"backups"`
}

func (b *BackupState) DeleteBackup(key string) {
	delete(b.Backups, key)
}

func (b *BackupState) VisitParent(startNode *BackupRecordWithKey, f func(r *BackupRecordWithKey) bool) error {
	currentNode := startNode
	loop := f(currentNode)

	for loop {
		om, err := b.GetRecordByKey(currentNode.ParentBackupKey)
		if err != nil {
			return fmt.Errorf("bs.GetRecordByKey failed %w", err)
		}

		currentNode = om
		loop = f(currentNode)
	}

	return nil
}

func (b *BackupState) GetRecordByKey(key string) (r *BackupRecordWithKey, err error) {
	if val, ok := b.Backups[key]; ok {
		return &BackupRecordWithKey{BackupRecord: val, Key: key}, nil
	}

	return nil, fmt.Errorf("backup for key %s does not exists", key)
}

type BackupRecord struct {
	Type            BackupType `yaml:"type"`
	ParentBackupKey string     `yaml:"parent-backup-key"`
	Size            *int64     `yaml:"size,omitempty"`
}

type BackupRecordWithKey struct {
	BackupRecord
	Key string
}

func (b *BackupRecordWithKey) GetFileName() string {
	return path.Base(b.Key)
}

type BackupStateRepo interface {
	Download(ctx context.Context, key string) (*BackupState, error)
	DownloadOrDefault(ctx context.Context, key string) (*BackupState, error)
	UpdateState(ctx context.Context, stateFileKey string, f func(state *BackupState) error) error
	Upload(ctx context.Context, key string, state *BackupState) error
}

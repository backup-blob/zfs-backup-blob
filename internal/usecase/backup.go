package usecase

import (
	"context"
	"errors"
	"fmt"
	"github.com/backup-blob/zfs-backup-blob/internal/domain"
	"github.com/backup-blob/zfs-backup-blob/internal/domain/config"
	"path"
	"strings"
)

type backup struct {
	backupRepo domain.BackupRepository
	stateRepo  domain.BackupStateRepo
	logger     domain.LogRepository
	configRepo config.ConfigRepo
}

func NewBackup(
	backupRepo domain.BackupRepository,
	stateRepo domain.BackupStateRepo,
	logger domain.LogRepository,
	configRepo config.ConfigRepo,
) domain.BackupUsecase {
	return &backup{ //nolint:whitespace // no need
		backupRepo: backupRepo,
		stateRepo:  stateRepo,
		logger:     logger,
		configRepo: configRepo,
	}
}

func (b *backup) BackupFull(ctx context.Context, fullSnapName string, updateHead bool) error {
	snap := domain.NewZfsSnapshot(fullSnapName)
	if snap == nil {
		return fmt.Errorf("snapshot name invalid %s", fullSnapName)
	}

	b.logger.Infof("creating full backup for snapshot %s with %d middlewares", fullSnapName, len(b.configRepo.GetMiddlewares()))

	res, err := b.backupRepo.Create(ctx, &domain.BackupCreate{
		Snapshot:        *snap,
		ProxyReaderFunc: domain.ChainMiddlewareRead(b.configRepo.GetMiddlewares()),
	})
	if err != nil {
		return fmt.Errorf("backupRepo.Create failed %w", err)
	}

	backupR := domain.BackupRecord{
		Type: domain.Full,
		Size: &res.Size,
	}

	return b.updateState(ctx, snap, backupR, updateHead)
}

func (b *backup) updateState(ctx context.Context, snap *domain.ZfsSnapshot, backupR domain.BackupRecord, updateHead bool) error {
	stateFileKey := fmt.Sprintf("%s/%s", snap.VolumeName, domain.StateFileName)

	b.logger.Infof("updating state file %s", stateFileKey)

	err := b.stateRepo.UpdateState(ctx, stateFileKey, func(state *domain.BackupState) error {
		state.Backups[snap.NormalizedFullPath()] = backupR
		if updateHead {
			state.Head = snap.NormalizedFullPath()
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("stateRepo.UpdateState failed %w", err)
	}

	return nil
}

func (b *backup) BackupIncremental(ctx context.Context, baseSnapName, newSnapName string, updateHead bool) error {
	snapNew := domain.NewZfsSnapshot(newSnapName)
	if snapNew == nil {
		return fmt.Errorf("new snapshot name invalid %s", newSnapName)
	}

	snapBase := domain.NewZfsSnapshot(baseSnapName)
	if snapBase == nil {
		return fmt.Errorf("base snapshot name invalid %s", baseSnapName)
	}

	b.logger.Infof("creating incremental backup for snapshot %s (base: %s) with %d middlewares", newSnapName, baseSnapName, len(b.configRepo.GetMiddlewares()))

	res, err := b.backupRepo.Create(ctx, &domain.BackupCreate{
		Snapshot:         *snapNew,
		PreviousSnapshot: snapBase,
		ProxyReaderFunc:  domain.ChainMiddlewareRead(b.configRepo.GetMiddlewares()),
	})

	if err != nil {
		return fmt.Errorf("backupRepo.Create failed %w", err)
	}

	backupR := domain.BackupRecord{
		Type:            domain.Incremental,
		ParentBackupKey: snapBase.NormalizedFullPath(),
		Size:            &res.Size,
	}

	return b.updateState(ctx, snapNew, backupR, updateHead)
}

func (b *backup) Restore(ctx context.Context, params *domain.RestoreParams) error {
	if strings.Contains(params.TargetZfsLocation, "@") {
		return errors.New("target needs to be volume/fs an not snapshot")
	}

	stateFilePath := fmt.Sprintf("%s/%s", path.Dir(params.BlobKey), domain.StateFileName)

	b.logger.Infof("reading state file %s for restore", stateFilePath)

	state, err := b.stateRepo.DownloadOrDefault(ctx, stateFilePath)
	if err != nil {
		return fmt.Errorf("stateRepo.DownloadOrDefault failed %w", err)
	}

	om, errK := state.GetRecordByKey(params.BlobKey)
	if errK != nil {
		return fmt.Errorf("state.GetRecordByKey failed: %w", errK)
	}

	switch om.Type {
	case domain.Full:
		return b.restoreOne(ctx, params, nil)
	case domain.Incremental:
		if params.RestoreAll {
			return b.restoreMany(ctx, params, om, state)
		} else {
			return b.restoreOne(ctx, params, nil)
		}
	default:
		return fmt.Errorf(`type %s unknown`, om.Type.String())
	}
}

func (b *backup) restoreMany(ctx context.Context, p *domain.RestoreParams, br *domain.BackupRecordWithKey, bs *domain.BackupState) error {
	restoreList, err := b.getRestoreList(br, bs)
	if err != nil {
		return fmt.Errorf("getRestoreList failed %w", err)
	}

	for _, item := range restoreList {
		errL := b.restoreOne(ctx, p, &item.Key)
		if errL != nil {
			return fmt.Errorf("restoreOne failed %w", errL)
		}
	}

	return nil
}

// todo: use VisitParent function on BackupState
func (b *backup) getRestoreList(startNode *domain.BackupRecordWithKey, bs *domain.BackupState) ([]*domain.BackupRecordWithKey, error) {
	var list []*domain.BackupRecordWithKey

	list = append([]*domain.BackupRecordWithKey{startNode}, list...)

	currentNode := startNode
	for currentNode.Type != domain.Full {
		om, err := bs.GetRecordByKey(currentNode.ParentBackupKey)
		if err != nil {
			return nil, fmt.Errorf("bs.GetRecordByKey failed %w", err)
		}

		currentNode = om
		list = append([]*domain.BackupRecordWithKey{om}, list...)
	}

	return list, nil
}

func (b *backup) restoreOne(ctx context.Context, p *domain.RestoreParams, overrideKey *string) error {
	key := p.BlobKey
	if overrideKey != nil {
		key = *overrideKey
	}

	b.logger.Infof("restoring %s to %s with %d middleware", key, p.TargetZfsLocation, len(b.configRepo.GetMiddlewares()))

	return b.backupRepo.Restore(ctx, &domain.BackupRestore{
		TargetZfsLocation: p.TargetZfsLocation,
		BlobKey:           key,
		ProxyWriterFunc:   domain.ChainMiddlewareWrite(b.configRepo.GetMiddlewares()),
	})
}

package usecase

import (
	"context"
	"fmt"
	"github.com/backup-blob/zfs-backup-blob/internal/domain"
	"github.com/backup-blob/zfs-backup-blob/internal/domain/config"
	"sort"
)

type trimUsecase struct {
	logger       domain.LogRepository
	volumeRepo   domain.VolumeRepository
	stateRepo    domain.BackupStateRepo
	configRepo   config.ConfigRepo
	backupRepo   domain.BackupRepository
	snapshotRepo domain.SnapshotRepository
}

func NewTrimUseCase(
	logger domain.LogRepository,
	volumeRepo domain.VolumeRepository,
	stateRepo domain.BackupStateRepo,
	configRepo config.ConfigRepo,
	backupRepo domain.BackupRepository,
	snapshotRepo domain.SnapshotRepository,
) domain.TrimUsecase {
	return &trimUsecase{ //nolint:whitespace // no need
		logger,
		volumeRepo,
		stateRepo,
		configRepo,
		backupRepo,
		snapshotRepo,
	}
}

func (t *trimUsecase) TrimRemote(ctx context.Context, pa *domain.TrimRemoteParameters) error {
	if err := t.validateRemotePolicy(); err != nil {
		return err
	}

	vols, err := t.volumeRepo.ListVolumesByGroup(pa.GroupName)
	if err != nil {
		return fmt.Errorf("volumeRepo.ListVolumesByGroup failed %w", err)
	}

	for _, vol := range vols {
		errT := t.trimVolumeRemote(ctx, vol, pa.DryRun)
		if errT != nil {
			return fmt.Errorf("trimVolumeRemote failed for volume %s: %w", vol.Name, errT)
		}
	}

	return nil
}

func (t *trimUsecase) TrimLocal(_ context.Context, pa *domain.TrimLocalParameters) error {
	if err := t.validateLocalPolicy(); err != nil {
		return err
	}

	vols, err := t.volumeRepo.ListVolumesByGroup(pa.GroupName)
	if err != nil {
		return fmt.Errorf("volumeRepo.ListVolumesByGroup failed %w", err)
	}

	for _, vol := range vols {
		errT := t.trimVolumeLocal(vol, pa.DryRun)
		if errT != nil {
			return fmt.Errorf("trimVolumeRemote failed for volume %s: %w", vol.Name, errT)
		}
	}

	return nil
}

func (t *trimUsecase) validateRemotePolicy() error {
	remotePolicy := t.configRepo.GetConfig().RemoteTrimPolicy
	if remotePolicy.GetFullCount() == 0 && remotePolicy.GetIncrementalCount() == 0 {
		return domain.ErrNoRemoteTrimPolicy
	}

	return nil
}

func (t *trimUsecase) validateLocalPolicy() error {
	policy := t.configRepo.GetConfig().LocalTrimPolicy
	if policy.GetFullCount() == 0 && policy.GetIncrementalCount() == 0 {
		return domain.ErrNoLocalTrimPolicy
	}

	return nil
}

func (t *trimUsecase) trimVolumeLocal(vol *domain.ZfsVolume, dryRun bool) error {
	snapshots, err := t.snapshotRepo.ListFilter(&domain.FilterCriteria{VolumeName: vol.Name, IgnoreInvalidSnapshotNames: true})
	if err != nil {
		return fmt.Errorf("snapshotRepo.ListFilter failed for volume %s: %w", vol.Name, err)
	}

	// sort desc
	sort.Slice(snapshots, func(i, j int) bool {
		return snapshots[i].Name > snapshots[j].Name
	})

	snapshotsToDelete, err := t.calcDeletableSnapshots(snapshots)
	if err != nil {
		return fmt.Errorf("calcDeletableSnapshots failed for volume %s: %w", vol.Name, err)
	}

	errD := t.deleteSnapshots(snapshotsToDelete, dryRun)
	if errD != nil {
		return fmt.Errorf("deleteSnapshots failed for volume %s: %w", vol.Name, errD)
	}

	return nil
}

func (t *trimUsecase) deleteSnapshots(snapshotsToDelete []*domain.ZfsSnapshot, dryRun bool) error {
	for _, snap := range snapshotsToDelete {
		t.logger.Debugf(fmt.Sprintf("Deleting snapshot (dryrun=%v) %s", dryRun, snap.FullName()))

		if !dryRun {
			errD := t.snapshotRepo.Delete(snap)
			if errD != nil {
				return fmt.Errorf("snapshotRepo.Delete failed for snapshot %s: %w", snap.FullName(), errD)
			}
		}
	}

	return nil
}

func (t *trimUsecase) trimVolumeRemote(ctx context.Context, vol *domain.ZfsVolume, dryRun bool) error {
	stateFilePath := fmt.Sprintf("%s/%s", vol.Name, domain.StateFileName)

	state, err := t.stateRepo.Download(ctx, stateFilePath)
	if err != nil {
		return fmt.Errorf("getState failed for volume %s: %w", vol.Name, err)
	}

	backups := t.getRemoteBackups(state)
	deleteResult := t.calcDeletableBackups(backups, state)
	deleteResult = t.undeleteDependants(deleteResult)

	for _, bup := range deleteResult.Delete {
		t.logger.Debugf(fmt.Sprintf("Deleting backup (dryrun=%v) %s", dryRun, bup.Key))

		if dryRun == false {
			errD := t.backupRepo.Delete(ctx, &domain.BackupDelete{BlobKey: bup.Key})
			if errD != nil {
				return fmt.Errorf("backupRepo.Delete failed for volume %s: %w", vol.Name, errD)
			}

			state.DeleteBackup(bup.Key)
		}
	}

	if dryRun == false {
		errU := t.stateRepo.Upload(ctx, stateFilePath, state)
		if errU != nil {
			return fmt.Errorf("stateRepo.Upload failed for volume %s: %w", vol.Name, errU)
		}
	}

	return nil
}

func (t *trimUsecase) undeleteDependants(tb *domain.DeleteResult) *domain.DeleteResult {
	var backupsToDelete []*domain.BackupRecordWithKey

	copyTb := domain.DeleteResult{}
	copyTb.Delete = append(copyTb.Delete, tb.Delete...)

	for _, b := range copyTb.Delete {
		if _, exists := tb.DependentsKeys[b.Key]; exists {
			t.logger.Debugf(fmt.Sprintf("not deleting %s since it has dependents", b.Key))
		} else {
			backupsToDelete = append(backupsToDelete, b)
		}
	}

	copyTb.Delete = backupsToDelete

	return &copyTb
}

func (t *trimUsecase) calcDeletableSnapshots(snapshots []*domain.ZfsSnapshot) ([]*domain.ZfsSnapshot, error) {
	conf := t.configRepo.GetConfig()

	fullBackupCounter := conf.LocalTrimPolicy.GetFullCount()
	incrementalBackupCounter := conf.LocalTrimPolicy.GetIncrementalCount()

	var snapshotsToDelete []*domain.ZfsSnapshot

	for _, snap := range snapshots {
		snapType, err := t.snapshotRepo.GetType(snap.FullName())
		if err != nil {
			return nil, fmt.Errorf("snapshotRepo.GetType failed for snapshot %s: %w", snap.FullName(), err)
		}

		switch snapType {
		case domain.Full:
			if fullBackupCounter > 0 {
				fullBackupCounter -= 1
			} else {
				snapshotsToDelete = append(snapshotsToDelete, snap)
			}
		case domain.Incremental:
			if incrementalBackupCounter > 0 {
				incrementalBackupCounter -= 1
			} else {
				snapshotsToDelete = append(snapshotsToDelete, snap)
			}
		}
	}

	return snapshotsToDelete, nil
}

func (t *trimUsecase) calcDeletableBackups(backups []*domain.BackupRecordWithKey, state *domain.BackupState) *domain.DeleteResult {
	trimBuckets := domain.DeleteResult{DependentsKeys: map[string]bool{}}
	conf := t.configRepo.GetConfig()

	fullBackupCounter := conf.RemoteTrimPolicy.GetFullCount()
	incrementalBackupCounter := conf.RemoteTrimPolicy.GetIncrementalCount()

	for _, bup := range backups {
		switch bup.Type {
		case domain.Full:
			if fullBackupCounter > 0 {
				fullBackupCounter -= 1
			} else {
				trimBuckets.Delete = append(trimBuckets.Delete, bup)
			}
		case domain.Incremental:
			if incrementalBackupCounter > 0 {
				incrementalBackupCounter -= 1

				state.VisitParent(bup, func(r *domain.BackupRecordWithKey) bool {
					trimBuckets.DependentsKeys[r.Key] = true
					return r.Type != domain.Full
				})
			} else {
				trimBuckets.Delete = append(trimBuckets.Delete, bup)
			}
		}
	}

	return &trimBuckets
}

func (t *trimUsecase) getRemoteBackups(state *domain.BackupState) []*domain.BackupRecordWithKey {
	var backups []*domain.BackupRecordWithKey

	for key, backupRecord := range state.Backups {
		backups = append(backups, &domain.BackupRecordWithKey{BackupRecord: backupRecord, Key: key})
	}

	// sort desc
	sort.Slice(backups, func(i, j int) bool {
		return backups[i].GetFileName() > backups[j].GetFileName()
	})

	return backups
}

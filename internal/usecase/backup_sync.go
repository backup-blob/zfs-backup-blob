package usecase

import (
	"context"
	"fmt"
	"github.com/backup-blob/zfs-backup-blob/internal/domain"
	"path/filepath"
	"sort"
)

type backupSync struct {
	backupUsecase  domain.BackupUsecase
	snapshotRepo   domain.SnapshotRepository
	namingStrategy domain.SnapshotNamestrategy
	stateRepo      domain.BackupStateRepo
	volumeRepo     domain.VolumeRepository
}

func NewBackupSync(
	backupUsecase domain.BackupUsecase,
	snapshotRepo domain.SnapshotRepository,
	stateRepo domain.BackupStateRepo,
	naming domain.SnapshotNamestrategy,
	volumeRepo domain.VolumeRepository,
) domain.BackupSyncUsecase { //
	return &backupSync{ //nolint:whitespace // no need
		backupUsecase:  backupUsecase,
		snapshotRepo:   snapshotRepo,
		stateRepo:      stateRepo,
		namingStrategy: naming,
		volumeRepo:     volumeRepo,
	}
}

func (b *backupSync) Backup(ctx context.Context, groupName string) error {
	if groupName == "" {
		groupName = domain.DefaultVolumeGroup
	}

	volumes, errV := b.volumeRepo.ListVolumesByGroup(groupName)
	if errV != nil {
		return fmt.Errorf("snapshotRepo.ListVolumesByGroup %w", errV)
	}

	for _, volume := range volumes {
		err := b.backupSnapshots(ctx, volume)
		if err != nil {
			return fmt.Errorf("b.backupSnapshots failed %w", err)
		}
	}

	return nil
}

func (b *backupSync) getSnapshotsSorted(volume *domain.ZfsVolume) ([]*domain.ZfsSnapshot, error) {
	filter := domain.FilterCriteria{VolumeName: volume.Name, IgnoreInvalidSnapshotNames: true}

	snapshots, err := b.snapshotRepo.ListFilter(&filter)
	if err != nil {
		return nil, fmt.Errorf("snapshotRepo.ListFilter failed %w", err)
	}

	sort.Slice(snapshots, func(i, j int) bool {
		return b.namingStrategy.IsGreater(snapshots[i].Name, snapshots[j].Name)
	})

	return snapshots, nil
}

func (b *backupSync) getStateFile(ctx context.Context, volume *domain.ZfsVolume) (*domain.BackupState, error) {
	stateFileKey := fmt.Sprintf("%s/%s", volume.Name, domain.StateFileName)

	state, err := b.stateRepo.DownloadOrDefault(ctx, stateFileKey)
	if err != nil {
		return nil, fmt.Errorf("stateRepo.DownloadOrDefault failed %w", err)
	}

	return state, nil
}

func (b *backupSync) backupSnapshots(ctx context.Context, volume *domain.ZfsVolume) error {
	snapshots, err := b.getSnapshotsSorted(volume)
	if err != nil {
		return err
	}

	state, errS := b.getStateFile(ctx, volume)
	if errS != nil {
		return errS
	}

	requests, errB := b.CalcSnapsToBackup(state, snapshots)
	if errB != nil {
		return errB
	}

	return b.doBackup(ctx, requests)
}

func (b *backupSync) doBackup(ctx context.Context, requests []*domain.BackupRequest) error {
	for index, req := range requests {
		if req.IsHead {
			continue // skip head -> its already available remote
		}

		switch req.Type {
		case domain.Incremental:
			previousIndex := index - 1
			if previousIndex < 0 {
				return fmt.Errorf("previous snapshot does not exists %s", req.Snapshot.FullName())
			}

			previousSnap := requests[previousIndex].Snapshot.FullName()

			errI := b.backupUsecase.BackupIncremental(ctx, previousSnap, req.Snapshot.FullName(), true)
			if errI != nil {
				return errI
			}
		case domain.Full:
			errF := b.backupUsecase.BackupFull(ctx, req.Snapshot.FullName(), true)
			if errF != nil {
				return errF
			}
		default:
			return fmt.Errorf("type not found %s", req.Type.String())
		}
	}

	return nil
}

func (b *backupSync) getBackupReq(snap *domain.ZfsSnapshot) (*domain.BackupRequest, error) {
	req := domain.BackupRequest{Snapshot: snap}

	backupType, errT := b.snapshotRepo.GetType(snap.FullName())
	if errT != nil {
		return nil, fmt.Errorf("snapshotRepo.GetType failed %s", snap.FullName())
	}

	req.Type = backupType

	return &req, nil
}

func (b *backupSync) validateBackupReq(isFirstBackup bool, remoteHead string, backupRequests []*domain.BackupRequest) error {
	var foundRemoteHead bool

	for index, req := range backupRequests {
		if req.IsHead {
			foundRemoteHead = true
		}

		if isFirstBackup && req.Type != domain.Full && index == 0 {
			return fmt.Errorf("first backup needs to be a full backup")
		}
	}

	if !foundRemoteHead && !isFirstBackup {
		return fmt.Errorf("remote head '%s' could not be found locally", remoteHead)
	}

	return nil
}

// todo: maybe move to own struct
func (b *backupSync) CalcSnapsToBackup(bs *domain.BackupState, snaps []*domain.ZfsSnapshot) ([]*domain.BackupRequest, error) {
	if len(snaps) == 0 {
		return nil, nil // not snaps to backup
	}

	remoteHead := filepath.Base(bs.Head)
	localHead := snaps[0].Name
	isFirstBackup := bs.Head == ""

	if remoteHead == localHead {
		return nil, nil // heads match
	}

	if b.namingStrategy.IsGreater(remoteHead, localHead) {
		return nil, fmt.Errorf("remote head is newer then local head %s > %s", remoteHead, localHead)
	}

	var (
		backupRequests []*domain.BackupRequest
	)

	for _, snap := range snaps {
		req, err := b.getBackupReq(snap)
		if err != nil {
			return nil, fmt.Errorf("b.getBackupReq failed %w", err)
		}

		if req.Type == domain.Unknown {
			continue // skip snapshot without tag
		}

		backupRequests = append([]*domain.BackupRequest{req}, backupRequests...)

		if snap.Name == remoteHead { // breaks only if remoteHead is set
			req.IsHead = true

			break
		}
	}

	errV := b.validateBackupReq(isFirstBackup, remoteHead, backupRequests)
	if errV != nil {
		return nil, fmt.Errorf("b.validateBackupReq failed %w", errV)
	}

	return backupRequests, nil
}

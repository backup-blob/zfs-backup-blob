package repository

import (
	"bytes"
	"fmt"
	"github.com/backup-blob/zfs-backup-blob/internal/domain"
	"strings"
)

const BackupTypeFlag = "backup_blob::type"

type snapshot struct {
	Zfs          domain.ZfsDriver
	Namestrategy domain.SnapshotNamestrategy
}

func NewSnapshot(zfs domain.ZfsDriver, namer domain.SnapshotNamestrategy) domain.SnapshotRepository {
	return &snapshot{
		Zfs:          zfs,
		Namestrategy: namer,
	}
}

func (s *snapshot) ListFilter(filter *domain.FilterCriteria) ([]*domain.ZfsSnapshot, error) {
	snaps := []*domain.ZfsSnapshot{}

	allSnaps, err := s.List()
	if err != nil {
		return nil, err
	}

	for _, snap := range allSnaps {
		if filter.VolumeName != "" && snap.VolumeName != filter.VolumeName {
			continue // skip snapshots which don't match volume name
		}

		if filter.IgnoreInvalidSnapshotNames && !s.Namestrategy.IsMatching(snap.Name) {
			continue // skip snapshots with invalid name
		}

		snaps = append(snaps, snap)
	}

	return snaps, nil
}

// TODO: use volname str instead of struct.
func (s *snapshot) Create(v *domain.ZfsVolume) (string, error) {
	snapName := v.Name + "@" + s.Namestrategy.GetName()

	cmd := s.Zfs.Snapshot(snapName)

	return snapName, cmd.Run()
}

func (s *snapshot) Delete(snap *domain.ZfsSnapshot) error {
	return s.Zfs.Destroy(snap.FullName())
}

func (s *snapshot) CreateWithType(v *domain.ZfsVolume, t domain.BackupType) (string, error) {
	snapFullName, err := s.Create(v)
	if err != nil {
		return snapFullName, fmt.Errorf("snapshot create failed %w", err)
	}

	errS := s.Zfs.SetField(BackupTypeFlag, t.String(), snapFullName)
	if errS != nil {
		return snapFullName, fmt.Errorf("set field on snapshot failed %w", errS)
	}

	return snapFullName, nil
}

func (s *snapshot) GetType(zfsEntity string) (domain.BackupType, error) {
	cmd := s.Zfs.GetField(BackupTypeFlag, zfsEntity)

	output, err := cmd.Output()
	if err != nil {
		return domain.Unknown, err
	}

	return domain.StringToBackupType(strings.TrimSpace(string(output))), nil
}

func (s *snapshot) List() ([]*domain.ZfsSnapshot, error) {
	var stdout bytes.Buffer

	cmd := s.Zfs.List(&domain.ListParameters{Type: []string{"snapshot"}, Fields: []string{"name"}})
	cmd.Stdout = &stdout

	err := cmd.Run()
	if err != nil {
		return nil, err
	}

	snaps := []*domain.ZfsSnapshot{}
	lines := strings.Split(stdout.String(), "\n")

	for _, line := range lines[1:] { // Skip the header line
		fields := strings.Fields(line)
		if len(fields) != 1 {
			continue
		}

		snap := domain.NewZfsSnapshot(fields[0])
		if snap == nil {
			return nil, fmt.Errorf("snapshot name invalid")
		}

		snaps = append(snaps, snap)
	}

	return snaps, nil
}

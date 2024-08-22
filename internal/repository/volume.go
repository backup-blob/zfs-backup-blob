package repository

import (
	"bytes"
	"fmt"
	"github.com/backup-blob/zfs-backup-blob/internal/domain"
	"os"
	"strings"
)

const BackupGroupFlag = "backup_blob::group"

type volumeRepo struct {
	zfs domain.ZfsDriver
}

func NewVolume(zfs domain.ZfsDriver) domain.VolumeRepository {
	return &volumeRepo{zfs: zfs}
}

func (s *volumeRepo) ListVolumes() ([]*domain.ZfsVolume, error) {
	var stdout bytes.Buffer

	params := domain.ListParameters{Type: []string{"volume", "filesystem"}, Fields: []string{"name", BackupGroupFlag}}
	cmd := s.zfs.List(&params)
	cmd.Stdout = &stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("zfs list volumes failed %w", err)
	}

	volumes := []*domain.ZfsVolume{}
	lines := strings.Split(stdout.String(), "\n")

	for _, line := range lines[1:] { // Skip the header line
		fields := strings.Fields(line)
		if len(fields) == 2 {
			name := fields[0]
			volume := &domain.ZfsVolume{Name: name, GroupName: fields[1]}
			volumes = append(volumes, volume)
		}
	}

	return volumes, nil
}

func (s *volumeRepo) ListVolumesByGroup(groupName string) ([]*domain.ZfsVolume, error) {
	volumes, err := s.ListVolumes()
	if err != nil {
		return nil, fmt.Errorf("listing volumes failed %w", err)
	}

	var volumesByGroup []*domain.ZfsVolume

	for _, volume := range volumes {
		if volume.GroupName == groupName {
			volumesByGroup = append(volumesByGroup, volume)
		}
	}

	return volumesByGroup, nil
}

func (s *volumeRepo) TagVolumeWithGroup(volume *domain.ZfsVolume) error {
	err := s.zfs.SetField(BackupGroupFlag, volume.GroupName, volume.Name)
	if err != nil {
		return fmt.Errorf("set zfs field failed %w", err)
	}

	return nil
}

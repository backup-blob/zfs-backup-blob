package usecase

import (
	"fmt"
	"github.com/backup-blob/zfs-backup-blob/internal/domain"
	"io"
)

type volumeUsecase struct {
	volumeRepository domain.VolumeRepository
	renderRepository domain.RenderRepository
	snapUsecase      domain.SnapshotUsecase
}

func NewVolumeUsecase(volumeRepository domain.VolumeRepository, renderRepository domain.RenderRepository, snapUsecase domain.SnapshotUsecase) domain.VolumeUsecase {
	return &volumeUsecase{
		volumeRepository: volumeRepository,
		renderRepository: renderRepository,
		snapUsecase:      snapUsecase,
	}
}

func (v *volumeUsecase) AddToGroup(volumeName, groupName string) error {
	if groupName == "" {
		groupName = domain.DefaultVolumeGroup
	}

	volume := &domain.ZfsVolume{Name: volumeName, GroupName: groupName}

	err := v.volumeRepository.TagVolumeWithGroup(volume)
	if err != nil {
		return fmt.Errorf("tagging volume failed: %w", err)
	}

	errC := v.snapUsecase.CreateByVolume(volume, domain.Full)
	if errC != nil {
		return fmt.Errorf("creating snapshot failed: %w", errC)
	}

	return nil
}

func (v *volumeUsecase) List(writer io.Writer) error {
	volumes, err := v.volumeRepository.ListVolumes()

	var volumesFiltered []*domain.ZfsVolume

	if err != nil {
		return fmt.Errorf("list volumes: %w", err)
	}

	for _, volume := range volumes {
		if volume.GroupName != "-" {
			volumesFiltered = append(volumesFiltered, volume)
		}
	}

	v.renderRepository.RenderVolumeTable(writer, volumesFiltered)

	return nil
}

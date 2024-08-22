package domain

import "io"

var DefaultVolumeGroup = "default"

type VolumeRepository interface {
	ListVolumes() ([]*ZfsVolume, error)
	ListVolumesByGroup(groupName string) ([]*ZfsVolume, error)
	TagVolumeWithGroup(volume *ZfsVolume) error
}

type VolumeUsecase interface {
	List(writer io.Writer) error
	AddToGroup(volumeName, groupName string) error
}

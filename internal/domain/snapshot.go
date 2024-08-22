package domain

type FilterCriteria struct {
	VolumeName                 string
	IgnoreInvalidSnapshotNames bool
}

type SnapshotRepository interface {
	Create(v *ZfsVolume) (string, error)
	CreateWithType(v *ZfsVolume, t BackupType) (string, error)
	List() ([]*ZfsSnapshot, error)
	ListFilter(filter *FilterCriteria) ([]*ZfsSnapshot, error)
	GetType(zfsEntity string) (BackupType, error)
	Delete(snap *ZfsSnapshot) error
}

type SnapshotUsecase interface {
	CreateByGroup(groupName string, backupType string) error
	CreateByVolume(volume *ZfsVolume, backupType BackupType) error
}

type SnapshotNamestrategy interface {
	GetName() string
	IsMatching(snapshotName string) bool
	IsGreater(snapNameA, snapNameB string) bool
}

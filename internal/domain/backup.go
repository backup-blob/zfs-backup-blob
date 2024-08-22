package domain

import (
	"context"
	"io"
)

const StateFileName = ".backupstate.yaml"

type BackupCreate struct {
	Snapshot         ZfsSnapshot
	PreviousSnapshot *ZfsSnapshot
	ProxyReaderFunc  func(r io.Reader) (io.Reader, error)
}

type BackupDelete struct {
	BlobKey string
}

type BackupRestore struct {
	TargetZfsLocation string
	BlobKey           string
	ProxyWriterFunc   func(w io.Writer) (io.Writer, error)
}

type RestoreParams struct {
	TargetZfsLocation string
	BlobKey           string
	RestoreAll        bool
}

type BackupRequest struct {
	Snapshot *ZfsSnapshot
	IsHead   bool
	Type     BackupType
}

type BackupRepository interface {
	Create(ctx context.Context, p *BackupCreate) (*UploadResponse, error)
	Restore(ctx context.Context, p *BackupRestore) error
	Delete(ctx context.Context, bd *BackupDelete) error
}

type BackupUsecase interface {
	BackupFull(ctx context.Context, fullSnapName string, updateHead bool) error
	BackupIncremental(ctx context.Context, baseSnapName, newSnapName string, updateHead bool) error
	Restore(ctx context.Context, params *RestoreParams) error
}

type BackupSyncUsecase interface {
	Backup(ctx context.Context, groupName string) error
	CalcSnapsToBackup(bs *BackupState, snaps []*ZfsSnapshot) ([]*BackupRequest, error)
}

type BackupListUsecase interface {
	List(ctx context.Context, volumeName string, writer io.Writer) error
}

type BackupType int

const (
	Unknown BackupType = iota
	Full
	Incremental
)

func (s BackupType) String() string {
	switch s {
	case Full:
		return "full"
	case Incremental:
		return "incremental"
	}

	return "unknown"
}

func StringToBackupType(str string) BackupType {
	switch str {
	case "full":
		return Full
	case "incremental":
		return Incremental
	default:
		return Unknown
	}
}

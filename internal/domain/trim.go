package domain

import (
	"context"
	"errors"
)

var ErrNoRemoteTrimPolicy = errors.New("no remote trim policy configured")
var ErrNoLocalTrimPolicy = errors.New("no local trim policy configured")

type DeleteResult struct {
	Delete         []*BackupRecordWithKey
	DependentsKeys map[string]bool
}

type TrimRemoteParameters struct {
	GroupName string
	DryRun    bool
}

type TrimLocalParameters struct {
	GroupName string
	DryRun    bool
}

type TrimUsecase interface {
	TrimRemote(ctx context.Context, pa *TrimRemoteParameters) error
	TrimLocal(_ context.Context, pa *TrimLocalParameters) error
}

// Code generated by MockGen. DO NOT EDIT.
// Source: internal/domain/backup_state.go
//
// Generated by this command:
//
//	mockgen -source=internal/domain/backup_state.go -destination=internal/domain/mocks/backup_state.go -package mocks
//

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	domain "github.com/backup-blob/zfs-backup-blob/internal/domain"
	gomock "go.uber.org/mock/gomock"
)

// MockBackupStateRepo is a mock of BackupStateRepo interface.
type MockBackupStateRepo struct {
	ctrl     *gomock.Controller
	recorder *MockBackupStateRepoMockRecorder
}

// MockBackupStateRepoMockRecorder is the mock recorder for MockBackupStateRepo.
type MockBackupStateRepoMockRecorder struct {
	mock *MockBackupStateRepo
}

// NewMockBackupStateRepo creates a new mock instance.
func NewMockBackupStateRepo(ctrl *gomock.Controller) *MockBackupStateRepo {
	mock := &MockBackupStateRepo{ctrl: ctrl}
	mock.recorder = &MockBackupStateRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBackupStateRepo) EXPECT() *MockBackupStateRepoMockRecorder {
	return m.recorder
}

// Download mocks base method.
func (m *MockBackupStateRepo) Download(ctx context.Context, key string) (*domain.BackupState, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Download", ctx, key)
	ret0, _ := ret[0].(*domain.BackupState)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Download indicates an expected call of Download.
func (mr *MockBackupStateRepoMockRecorder) Download(ctx, key any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Download", reflect.TypeOf((*MockBackupStateRepo)(nil).Download), ctx, key)
}

// DownloadOrDefault mocks base method.
func (m *MockBackupStateRepo) DownloadOrDefault(ctx context.Context, key string) (*domain.BackupState, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DownloadOrDefault", ctx, key)
	ret0, _ := ret[0].(*domain.BackupState)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DownloadOrDefault indicates an expected call of DownloadOrDefault.
func (mr *MockBackupStateRepoMockRecorder) DownloadOrDefault(ctx, key any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DownloadOrDefault", reflect.TypeOf((*MockBackupStateRepo)(nil).DownloadOrDefault), ctx, key)
}

// UpdateState mocks base method.
func (m *MockBackupStateRepo) UpdateState(ctx context.Context, stateFileKey string, f func(*domain.BackupState) error) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateState", ctx, stateFileKey, f)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateState indicates an expected call of UpdateState.
func (mr *MockBackupStateRepoMockRecorder) UpdateState(ctx, stateFileKey, f any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateState", reflect.TypeOf((*MockBackupStateRepo)(nil).UpdateState), ctx, stateFileKey, f)
}

// Upload mocks base method.
func (m *MockBackupStateRepo) Upload(ctx context.Context, key string, state *domain.BackupState) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Upload", ctx, key, state)
	ret0, _ := ret[0].(error)
	return ret0
}

// Upload indicates an expected call of Upload.
func (mr *MockBackupStateRepoMockRecorder) Upload(ctx, key, state any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Upload", reflect.TypeOf((*MockBackupStateRepo)(nil).Upload), ctx, key, state)
}

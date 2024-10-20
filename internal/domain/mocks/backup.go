// Code generated by MockGen. DO NOT EDIT.
// Source: internal/domain/backup.go
//
// Generated by this command:
//
//	mockgen -source=internal/domain/backup.go -destination=internal/domain/mocks/backup.go -package mocks
//

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	io "io"
	reflect "reflect"

	domain "github.com/backup-blob/zfs-backup-blob/internal/domain"
	gomock "go.uber.org/mock/gomock"
)

// MockBackupRepository is a mock of BackupRepository interface.
type MockBackupRepository struct {
	ctrl     *gomock.Controller
	recorder *MockBackupRepositoryMockRecorder
}

// MockBackupRepositoryMockRecorder is the mock recorder for MockBackupRepository.
type MockBackupRepositoryMockRecorder struct {
	mock *MockBackupRepository
}

// NewMockBackupRepository creates a new mock instance.
func NewMockBackupRepository(ctrl *gomock.Controller) *MockBackupRepository {
	mock := &MockBackupRepository{ctrl: ctrl}
	mock.recorder = &MockBackupRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBackupRepository) EXPECT() *MockBackupRepositoryMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockBackupRepository) Create(ctx context.Context, p *domain.BackupCreate) (*domain.UploadResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, p)
	ret0, _ := ret[0].(*domain.UploadResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockBackupRepositoryMockRecorder) Create(ctx, p any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockBackupRepository)(nil).Create), ctx, p)
}

// Delete mocks base method.
func (m *MockBackupRepository) Delete(ctx context.Context, bd *domain.BackupDelete) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, bd)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockBackupRepositoryMockRecorder) Delete(ctx, bd any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockBackupRepository)(nil).Delete), ctx, bd)
}

// Restore mocks base method.
func (m *MockBackupRepository) Restore(ctx context.Context, p *domain.BackupRestore) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Restore", ctx, p)
	ret0, _ := ret[0].(error)
	return ret0
}

// Restore indicates an expected call of Restore.
func (mr *MockBackupRepositoryMockRecorder) Restore(ctx, p any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Restore", reflect.TypeOf((*MockBackupRepository)(nil).Restore), ctx, p)
}

// MockBackupUsecase is a mock of BackupUsecase interface.
type MockBackupUsecase struct {
	ctrl     *gomock.Controller
	recorder *MockBackupUsecaseMockRecorder
}

// MockBackupUsecaseMockRecorder is the mock recorder for MockBackupUsecase.
type MockBackupUsecaseMockRecorder struct {
	mock *MockBackupUsecase
}

// NewMockBackupUsecase creates a new mock instance.
func NewMockBackupUsecase(ctrl *gomock.Controller) *MockBackupUsecase {
	mock := &MockBackupUsecase{ctrl: ctrl}
	mock.recorder = &MockBackupUsecaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBackupUsecase) EXPECT() *MockBackupUsecaseMockRecorder {
	return m.recorder
}

// BackupFull mocks base method.
func (m *MockBackupUsecase) BackupFull(ctx context.Context, fullSnapName string, updateHead bool) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BackupFull", ctx, fullSnapName, updateHead)
	ret0, _ := ret[0].(error)
	return ret0
}

// BackupFull indicates an expected call of BackupFull.
func (mr *MockBackupUsecaseMockRecorder) BackupFull(ctx, fullSnapName, updateHead any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BackupFull", reflect.TypeOf((*MockBackupUsecase)(nil).BackupFull), ctx, fullSnapName, updateHead)
}

// BackupIncremental mocks base method.
func (m *MockBackupUsecase) BackupIncremental(ctx context.Context, baseSnapName, newSnapName string, updateHead bool) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BackupIncremental", ctx, baseSnapName, newSnapName, updateHead)
	ret0, _ := ret[0].(error)
	return ret0
}

// BackupIncremental indicates an expected call of BackupIncremental.
func (mr *MockBackupUsecaseMockRecorder) BackupIncremental(ctx, baseSnapName, newSnapName, updateHead any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BackupIncremental", reflect.TypeOf((*MockBackupUsecase)(nil).BackupIncremental), ctx, baseSnapName, newSnapName, updateHead)
}

// Restore mocks base method.
func (m *MockBackupUsecase) Restore(ctx context.Context, params *domain.RestoreParams) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Restore", ctx, params)
	ret0, _ := ret[0].(error)
	return ret0
}

// Restore indicates an expected call of Restore.
func (mr *MockBackupUsecaseMockRecorder) Restore(ctx, params any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Restore", reflect.TypeOf((*MockBackupUsecase)(nil).Restore), ctx, params)
}

// MockBackupSyncUsecase is a mock of BackupSyncUsecase interface.
type MockBackupSyncUsecase struct {
	ctrl     *gomock.Controller
	recorder *MockBackupSyncUsecaseMockRecorder
}

// MockBackupSyncUsecaseMockRecorder is the mock recorder for MockBackupSyncUsecase.
type MockBackupSyncUsecaseMockRecorder struct {
	mock *MockBackupSyncUsecase
}

// NewMockBackupSyncUsecase creates a new mock instance.
func NewMockBackupSyncUsecase(ctrl *gomock.Controller) *MockBackupSyncUsecase {
	mock := &MockBackupSyncUsecase{ctrl: ctrl}
	mock.recorder = &MockBackupSyncUsecaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBackupSyncUsecase) EXPECT() *MockBackupSyncUsecaseMockRecorder {
	return m.recorder
}

// Backup mocks base method.
func (m *MockBackupSyncUsecase) Backup(ctx context.Context, groupName string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Backup", ctx, groupName)
	ret0, _ := ret[0].(error)
	return ret0
}

// Backup indicates an expected call of Backup.
func (mr *MockBackupSyncUsecaseMockRecorder) Backup(ctx, groupName any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Backup", reflect.TypeOf((*MockBackupSyncUsecase)(nil).Backup), ctx, groupName)
}

// CalcSnapsToBackup mocks base method.
func (m *MockBackupSyncUsecase) CalcSnapsToBackup(bs *domain.BackupState, snaps []*domain.ZfsSnapshot) ([]*domain.BackupRequest, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CalcSnapsToBackup", bs, snaps)
	ret0, _ := ret[0].([]*domain.BackupRequest)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CalcSnapsToBackup indicates an expected call of CalcSnapsToBackup.
func (mr *MockBackupSyncUsecaseMockRecorder) CalcSnapsToBackup(bs, snaps any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CalcSnapsToBackup", reflect.TypeOf((*MockBackupSyncUsecase)(nil).CalcSnapsToBackup), bs, snaps)
}

// MockBackupListUsecase is a mock of BackupListUsecase interface.
type MockBackupListUsecase struct {
	ctrl     *gomock.Controller
	recorder *MockBackupListUsecaseMockRecorder
}

// MockBackupListUsecaseMockRecorder is the mock recorder for MockBackupListUsecase.
type MockBackupListUsecaseMockRecorder struct {
	mock *MockBackupListUsecase
}

// NewMockBackupListUsecase creates a new mock instance.
func NewMockBackupListUsecase(ctrl *gomock.Controller) *MockBackupListUsecase {
	mock := &MockBackupListUsecase{ctrl: ctrl}
	mock.recorder = &MockBackupListUsecaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBackupListUsecase) EXPECT() *MockBackupListUsecaseMockRecorder {
	return m.recorder
}

// List mocks base method.
func (m *MockBackupListUsecase) List(ctx context.Context, volumeName string, writer io.Writer) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", ctx, volumeName, writer)
	ret0, _ := ret[0].(error)
	return ret0
}

// List indicates an expected call of List.
func (mr *MockBackupListUsecaseMockRecorder) List(ctx, volumeName, writer any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockBackupListUsecase)(nil).List), ctx, volumeName, writer)
}

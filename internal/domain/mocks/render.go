// Code generated by MockGen. DO NOT EDIT.
// Source: internal/domain/render.go
//
// Generated by this command:
//
//	mockgen -source=internal/domain/render.go -destination=internal/domain/mocks/render.go -package mocks
//

// Package mocks is a generated GoMock package.
package mocks

import (
	io "io"
	reflect "reflect"

	domain "github.com/backup-blob/zfs-backup-blob/internal/domain"
	gomock "go.uber.org/mock/gomock"
)

// MockRenderDriver is a mock of RenderDriver interface.
type MockRenderDriver struct {
	ctrl     *gomock.Controller
	recorder *MockRenderDriverMockRecorder
}

// MockRenderDriverMockRecorder is the mock recorder for MockRenderDriver.
type MockRenderDriverMockRecorder struct {
	mock *MockRenderDriver
}

// NewMockRenderDriver creates a new mock instance.
func NewMockRenderDriver(ctrl *gomock.Controller) *MockRenderDriver {
	mock := &MockRenderDriver{ctrl: ctrl}
	mock.recorder = &MockRenderDriverMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRenderDriver) EXPECT() *MockRenderDriverMockRecorder {
	return m.recorder
}

// RenderTable mocks base method.
func (m *MockRenderDriver) RenderTable(writer io.Writer, headerRow []any, rows [][]any) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "RenderTable", writer, headerRow, rows)
}

// RenderTable indicates an expected call of RenderTable.
func (mr *MockRenderDriverMockRecorder) RenderTable(writer, headerRow, rows any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RenderTable", reflect.TypeOf((*MockRenderDriver)(nil).RenderTable), writer, headerRow, rows)
}

// MockRenderRepository is a mock of RenderRepository interface.
type MockRenderRepository struct {
	ctrl     *gomock.Controller
	recorder *MockRenderRepositoryMockRecorder
}

// MockRenderRepositoryMockRecorder is the mock recorder for MockRenderRepository.
type MockRenderRepositoryMockRecorder struct {
	mock *MockRenderRepository
}

// NewMockRenderRepository creates a new mock instance.
func NewMockRenderRepository(ctrl *gomock.Controller) *MockRenderRepository {
	mock := &MockRenderRepository{ctrl: ctrl}
	mock.recorder = &MockRenderRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRenderRepository) EXPECT() *MockRenderRepositoryMockRecorder {
	return m.recorder
}

// RenderBackupTable mocks base method.
func (m *MockRenderRepository) RenderBackupTable(writer io.Writer, backups []domain.BackupRecordWithKey) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "RenderBackupTable", writer, backups)
}

// RenderBackupTable indicates an expected call of RenderBackupTable.
func (mr *MockRenderRepositoryMockRecorder) RenderBackupTable(writer, backups any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RenderBackupTable", reflect.TypeOf((*MockRenderRepository)(nil).RenderBackupTable), writer, backups)
}

// RenderVolumeTable mocks base method.
func (m *MockRenderRepository) RenderVolumeTable(writer io.Writer, volumes []*domain.ZfsVolume) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "RenderVolumeTable", writer, volumes)
}

// RenderVolumeTable indicates an expected call of RenderVolumeTable.
func (mr *MockRenderRepositoryMockRecorder) RenderVolumeTable(writer, volumes any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RenderVolumeTable", reflect.TypeOf((*MockRenderRepository)(nil).RenderVolumeTable), writer, volumes)
}

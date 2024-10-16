// Code generated by MockGen. DO NOT EDIT.
// Source: internal/domain/storage.go
//
// Generated by this command:
//
//	mockgen -source=internal/domain/storage.go -destination=internal/domain/mocks/storage.go -package mocks
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

// MockStorageDriver is a mock of StorageDriver interface.
type MockStorageDriver struct {
	ctrl     *gomock.Controller
	recorder *MockStorageDriverMockRecorder
}

// MockStorageDriverMockRecorder is the mock recorder for MockStorageDriver.
type MockStorageDriverMockRecorder struct {
	mock *MockStorageDriver
}

// NewMockStorageDriver creates a new mock instance.
func NewMockStorageDriver(ctrl *gomock.Controller) *MockStorageDriver {
	mock := &MockStorageDriver{ctrl: ctrl}
	mock.recorder = &MockStorageDriverMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStorageDriver) EXPECT() *MockStorageDriverMockRecorder {
	return m.recorder
}

// Delete mocks base method.
func (m *MockStorageDriver) Delete(ctx context.Context, dp *domain.DeleteParameters) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, dp)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockStorageDriverMockRecorder) Delete(ctx, dp any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockStorageDriver)(nil).Delete), ctx, dp)
}

// Download mocks base method.
func (m *MockStorageDriver) Download(ctx context.Context, dp *domain.DownloadParameters, writer io.Writer) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Download", ctx, dp, writer)
	ret0, _ := ret[0].(error)
	return ret0
}

// Download indicates an expected call of Download.
func (mr *MockStorageDriverMockRecorder) Download(ctx, dp, writer any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Download", reflect.TypeOf((*MockStorageDriver)(nil).Download), ctx, dp, writer)
}

// Upload mocks base method.
func (m *MockStorageDriver) Upload(ctx context.Context, up *domain.UploadParameters, reader io.Reader) (*domain.UploadResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Upload", ctx, up, reader)
	ret0, _ := ret[0].(*domain.UploadResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Upload indicates an expected call of Upload.
func (mr *MockStorageDriverMockRecorder) Upload(ctx, up, reader any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Upload", reflect.TypeOf((*MockStorageDriver)(nil).Upload), ctx, up, reader)
}

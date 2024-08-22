package repository_test

import (
	"context"
	"errors"
	"github.com/backup-blob/zfs-backup-blob/internal/domain"
	"github.com/backup-blob/zfs-backup-blob/internal/domain/mocks"
	"github.com/backup-blob/zfs-backup-blob/internal/repository"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/mock/gomock"
	"io"
	"strings"
	"testing"
)

var fixtureBackups = "head: \"\"\nbackups:\n    folder/folder1/key:\n        type: 1\n        parent-backup-key: folder/folder1/parent\n        size: 1\n"

func TestSpecBackupState(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockStorage := mocks.NewMockStorageDriver(ctrl)
	backupState := repository.NewBackupStateRepo(mockStorage)
	ctx := context.Background()
	key := "folder/folder1/key"
	parentKey := "folder/folder1/parent"

	Convey("Given i call the UpdateState function", t, func() {
		Convey("When everything works as expected", func() {
			Convey("It should update the state", func() {
				mockStorage.EXPECT().Download(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				mockStorage.EXPECT().Upload(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)

				err := backupState.UpdateState(ctx, "key", func(state *domain.BackupState) error {
					So(state, ShouldNotBeNil)
					return nil
				})

				So(err, ShouldBeNil)
			})
		})
	})

	Convey("Given i call the DownloadOrDefault function", t, func() {
		Convey("When everything works as expected", func() {
			Convey("It downloads the state", func() {
				mockStorage.EXPECT().Download(gomock.Any(), gomock.Any(), gomock.Any()).Times(1).DoAndReturn(func(_ context.Context, dp *domain.DownloadParameters, w io.Writer) error {
					io.Copy(w, strings.NewReader(fixtureBackups))
					return nil
				})

				state, err := backupState.DownloadOrDefault(ctx, key)

				So(err, ShouldBeNil)
				So(len(state.Backups), ShouldEqual, 1)
			})
		})
		Convey("When the file does not exist", func() {
			Convey("It should return the default state", func() {
				mockStorage.EXPECT().Download(gomock.Any(), gomock.Any(), gomock.Any()).Return(domain.ErrNotFound)

				state, err := backupState.DownloadOrDefault(ctx, key)

				So(err, ShouldBeNil)
				So(state.Backups, ShouldNotBeNil)
			})
		})
	})

	Convey("Given i call the Download function", t, func() {
		Convey("When everything works as expected", func() {
			Convey("It downloads the state", func() {
				mockStorage.EXPECT().Download(gomock.Any(), gomock.Any(), gomock.Any()).Times(1).DoAndReturn(func(_ context.Context, dp *domain.DownloadParameters, w io.Writer) error {
					io.Copy(w, strings.NewReader(fixtureBackups))
					return nil
				})

				state, err := backupState.Download(ctx, key)

				So(err, ShouldBeNil)
				So(len(state.Backups), ShouldEqual, 1)
				So(state.Backups[key].ParentBackupKey, ShouldEqual, parentKey)
				So(state.Backups[key].Type, ShouldEqual, domain.Full)
			})
		})
	})

	Convey("Given i call the Upload function", t, func() {
		Convey("When everything works as expected", func() {
			Convey("It should upload to the storage", func() {
				mockStorage.EXPECT().Upload(gomock.Any(), gomock.Any(), gomock.Any()).Times(1).DoAndReturn(func(_ context.Context, up *domain.UploadParameters, r io.Reader) (*domain.UploadResponse, error) {
					b, _ := io.ReadAll(r)
					So(string(b), ShouldEqual, fixtureBackups)
					So(up.Key, ShouldEqual, key)
					return nil, nil
				})
				size := int64(1)
				backups := map[string]domain.BackupRecord{
					key: {ParentBackupKey: parentKey, Type: domain.Full, Size: &size},
				}

				err := backupState.Upload(ctx, key, &domain.BackupState{Backups: backups})

				So(err, ShouldBeNil)
			})
		})
		Convey("When the download fails", func() {
			Convey("It should error", func() {
				mockStorage.EXPECT().Download(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("error"))

				_, err := backupState.Download(ctx, key)

				So(err, ShouldNotBeNil)
			})
		})
	})
}

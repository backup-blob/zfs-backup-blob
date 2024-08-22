package repository_test

import (
	"bytes"
	"context"
	"github.com/backup-blob/zfs-backup-blob/internal/domain"
	"github.com/backup-blob/zfs-backup-blob/internal/domain/mocks"
	"github.com/backup-blob/zfs-backup-blob/internal/repository"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/mock/gomock"
	"io"
	"os/exec"
	"strings"
	"testing"
)

func TestSpecBackup(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockZfs := mocks.NewMockZfsDriver(ctrl)
	mockStorage := mocks.NewMockStorageDriver(ctrl)
	backup := repository.NewBackup(mockZfs, mockStorage)
	volumeName := "zfs/folder1/folder2"
	payload := "hello"
	snapName := "backup_123"
	blobKey := "zfs/folder1/folder2/backup_123"
	uploadResp := domain.UploadResponse{}

	Convey("Given i call the Delete function", t, func() {
		Convey("When everything works as expected", func() {
			Convey("It should delete the backup", func() {
				mockStorage.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(nil)
				err := backup.Delete(context.Background(), &domain.BackupDelete{
					BlobKey: "path/path1/key",
				})

				So(err, ShouldBeNil)
			})
		})
	})

	Convey("Given i call the Create function", t, func() {
		Convey("With snapshot parameters", func() {
			Convey("It should create a backup", func() {
				mockZfs.EXPECT().Send(gomock.Any()).Times(1).Return(fakeCommand(payload))
				mockStorage.EXPECT().Upload(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, up *domain.UploadParameters, reader io.Reader) (*domain.UploadResponse, error) {
					So(up.Key, ShouldEqual, blobKey)
					data, err := io.ReadAll(reader)
					So(len(data), ShouldEqual, len(payload))
					return &uploadResp, err
				})
				res, err := backup.Create(context.Background(), &domain.BackupCreate{
					Snapshot: domain.ZfsSnapshot{Name: snapName, VolumeName: volumeName},
				})

				So(err, ShouldBeNil)
				So(res, ShouldEqual, &uploadResp)
			})
		})
		Convey("With a proxyreader parameters", func() {
			Convey("It should pass the data through the proxy", func() {
				mockZfs.EXPECT().Send(gomock.Any()).Times(1).Return(fakeCommand(payload))
				mockStorage.EXPECT().Upload(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, up *domain.UploadParameters, reader io.Reader) (*domain.UploadResponse, error) {
					_, err := io.ReadAll(reader)
					return &uploadResp, err
				})
				var fReader fakeReader

				res, err := backup.Create(context.Background(), &domain.BackupCreate{
					Snapshot: domain.ZfsSnapshot{Name: snapName, VolumeName: volumeName},
					ProxyReaderFunc: func(r io.Reader) (io.Reader, error) {
						fReader = fakeReader{R: r}
						return &fReader, nil
					},
				})

				So(err, ShouldBeNil)
				So(fReader.ByteRead, ShouldEqual, len(payload))
				So(res, ShouldEqual, &uploadResp)
			})
		})
	})

	Convey("Given i call the Restore function", t, func() {
		Convey("With restore parameters", func() {
			Convey("It should restore", func(c C) {
				byteWriter := bytes.NewBuffer([]byte{})
				mockZfs.EXPECT().Receive(gomock.Any()).Times(1).Return(fakeCommandWriter(byteWriter))
				mockStorage.EXPECT().Download(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, dp *domain.DownloadParameters, writer io.Writer) error {
					c.So(dp.Bucket, ShouldEqual, "")
					c.So(dp.Key, ShouldEqual, blobKey)
					_, err := io.Copy(writer, strings.NewReader(payload))
					return err
				})

				err := backup.Restore(context.Background(), &domain.BackupRestore{
					TargetZfsLocation: volumeName,
					BlobKey:           blobKey,
				})

				So(err, ShouldBeNil)
				So(byteWriter.String(), ShouldEqual, payload)
			})
		})
		Convey("With writeProxy parameter", func() {
			Convey("It should pipe writes through proxy", func(c C) {
				byteWriter := bytes.NewBuffer([]byte{})
				mockZfs.EXPECT().Receive(gomock.Any()).Times(1).Return(fakeCommandWriter(byteWriter))
				mockStorage.EXPECT().Download(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, dp *domain.DownloadParameters, writer io.Writer) error {
					_, err := io.Copy(writer, strings.NewReader(payload))
					return err
				})
				var writer fakeWriter
				err := backup.Restore(context.Background(), &domain.BackupRestore{
					TargetZfsLocation: volumeName,
					BlobKey:           blobKey,
					ProxyWriterFunc: func(w io.Writer) (io.Writer, error) {
						writer = fakeWriter{W: w}
						return &writer, nil
					},
				})

				So(err, ShouldBeNil)
				So(byteWriter.String(), ShouldEqual, payload)
				So(writer.BytesWritten, ShouldEqual, len(payload))
			})
		})
	})
}

func fakeCommandWriter(w io.Writer) *exec.Cmd {
	cmd := exec.Command("cat", "-")
	cmd.Stdout = w
	return cmd
}

type fakeWriter struct {
	W            io.Writer
	BytesWritten int
}

func (w *fakeWriter) Write(p []byte) (int, error) {
	n, err := w.W.Write(p)
	w.BytesWritten += n
	return n, err
}

type fakeReader struct {
	R        io.Reader
	ByteRead int
}

func (r *fakeReader) Read(p []byte) (int, error) {
	n, err := r.R.Read(p)
	r.ByteRead += n
	return n, err
}

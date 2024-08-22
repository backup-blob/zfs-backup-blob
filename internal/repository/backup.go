package repository

import (
	"context"
	"github.com/backup-blob/zfs-backup-blob/internal/domain"
	"io"
	"os"
)

type backup struct {
	Zfs     domain.ZfsDriver
	Storage domain.StorageDriver
}

func NewBackup(zfs domain.ZfsDriver, storage domain.StorageDriver) domain.BackupRepository {
	return &backup{
		Zfs:     zfs,
		Storage: storage,
	}
}

func (b *backup) Delete(ctx context.Context, bd *domain.BackupDelete) error {
	return b.Storage.Delete(ctx, &domain.DeleteParameters{Key: bd.BlobKey})
}

func (b *backup) Create(ctx context.Context, param *domain.BackupCreate) (*domain.UploadResponse, error) {
	var reader io.Reader

	uploadParam := domain.UploadParameters{
		Key: param.Snapshot.NormalizedFullPath(),
	}

	sendParams := domain.SendParameters{
		WithParameters: true,
		SnapshotName:   param.Snapshot.FullName(),
	}
	if param.PreviousSnapshot != nil {
		sendParams.PreviousSnapshotName = param.PreviousSnapshot.FullName()
	}

	send := b.Zfs.Send(&sendParams)
	sender, err := send.StdoutPipe()
	send.Stderr = os.Stdout

	if err != nil {
		return nil, err
	}

	err = send.Start()

	if err != nil {
		return nil, err
	}

	defer sender.Close()

	reader = sender

	if param.ProxyReaderFunc != nil {
		proxiedReader, errP := param.ProxyReaderFunc(reader)
		reader = proxiedReader

		if errP != nil {
			return nil, errP
		}
	}

	uploadResp, errU := b.Storage.Upload(ctx, &uploadParam, reader)
	if errU != nil {
		return nil, errU
	}

	if errW := send.Wait(); errW != nil {
		return nil, errW
	}

	return uploadResp, nil
}

func (b *backup) Restore(ctx context.Context, param *domain.BackupRestore) error {
	var (
		writer      io.Writer
		customError error
	)

	receive := b.Zfs.Receive(&domain.ReceiveParameters{TargetName: param.TargetZfsLocation})
	receiver, err := receive.StdinPipe()
	receive.Stderr = os.Stdout
	writer = receiver

	if err != nil {
		return err
	}

	err = receive.Start()

	if err != nil {
		return err
	}

	if param.ProxyWriterFunc != nil {
		proxyWriter, errP := param.ProxyWriterFunc(receiver)
		writer = proxyWriter

		if errP != nil {
			return errP
		}
	}

	go func() {
		defer receiver.Close()

		customError = b.Storage.Download(ctx, &domain.DownloadParameters{Key: param.BlobKey}, writer)
	}()

	err = receive.Wait()

	if customError != nil {
		return customError
	}

	return err
}

package domain

import (
	"context"
	"errors"
	"io"
)

var ErrNotFound = errors.New("object not found")

type UploadParameters struct {
	Bucket string
	Key    string
}

type DeleteParameters struct {
	Bucket string
	Key    string
}

type DownloadParameters struct {
	Bucket string
	Key    string
}

type UploadResponse struct {
	Size int64
}

type StorageDriver interface {
	Delete(ctx context.Context, dp *DeleteParameters) error
	Upload(ctx context.Context, up *UploadParameters, reader io.Reader) (*UploadResponse, error)
	Download(ctx context.Context, dp *DownloadParameters, writer io.Writer) error
}

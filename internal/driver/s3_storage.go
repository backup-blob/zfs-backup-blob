package driver

import (
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/backup-blob/zfs-backup-blob/internal/domain"
	config2 "github.com/backup-blob/zfs-backup-blob/internal/domain/config"
	"io"
	"net/http"
	"os"
)

type S3Storage struct {
	Client         *s3.Client
	Logger         domain.LogDriver
	Bucket         *string
	Prefix         string
	uploadPartSize int64
}

var responseError interface {
	HTTPStatusCode() int
}

func NewS3Storage(client *s3.Client, logger domain.LogDriver) domain.StorageDriver {
	return &S3Storage{Client: client, Logger: logger}
}

func NewS3StorageFromConfig(s3c *config2.S3Config, logger domain.LogDriver) (domain.StorageDriver, error) {
	instance := S3Storage{Logger: logger}

	setCredentials(s3c)

	sdkConfig, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return nil, fmt.Errorf("loading storage config error %w", err)
	}

	s3Client := s3.NewFromConfig(sdkConfig, func(options *s3.Options) {
		if s3c.BaseEndpoint != "" {
			options.BaseEndpoint = aws.String(s3c.BaseEndpoint)
		}

		if s3c.UsePathStyle {
			options.UsePathStyle = true
		}

		if s3c.Region != "" {
			options.Region = s3c.Region
		}

		options.RetryMaxAttempts = s3c.MaxRetries
	})

	instance.Client = s3Client
	instance.Prefix = s3c.Prefix
	instance.uploadPartSize = int64(s3c.UploadPartSize)

	if s3c.Bucket != "" {
		instance.Bucket = &s3c.Bucket
	}

	return &instance, nil
}

func setCredentials(s3c *config2.S3Config) {
	if s3c.AccessKey != "" {
		os.Setenv("AWS_ACCESS_KEY_ID", s3c.AccessKey)
	}

	if s3c.AccessSecret != "" {
		os.Setenv("AWS_SECRET_ACCESS_KEY", s3c.AccessSecret)
	}
}

func (s *S3Storage) getBucket(bucketOverride string) string {
	if bucketOverride != "" {
		return bucketOverride
	}

	if s.Bucket != nil {
		return *s.Bucket
	}

	return ""
}

func (s *S3Storage) withPrefix(keyPath string) string {
	return s.Prefix + keyPath
}

func (s *S3Storage) Delete(ctx context.Context, dp *domain.DeleteParameters) error {
	_, err := s.Client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.getBucket(dp.Bucket)),
		Key:    aws.String(s.withPrefix(dp.Key)),
	})

	return err
}

func (s *S3Storage) Upload(ctx context.Context, up *domain.UploadParameters, reader io.Reader) (*domain.UploadResponse, error) {
	s.Logger.Debugf("starting upload %v", up)

	uploader := manager.NewUploader(s.Client, func(uploader *manager.Uploader) {
		uploader.Concurrency = 1
		uploader.PartSize = s.uploadPartSize
	})

	fakeReader := FakeReader{r: reader}
	_, err := uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.getBucket(up.Bucket)),
		Key:    aws.String(s.withPrefix(up.Key)),
		Body:   &fakeReader,
		//ChecksumAlgorithm: types.ChecksumAlgorithmSha1, TODO: add back
	})

	if err != nil {
		return nil, fmt.Errorf("upload to storage error %w", err)
	}

	size, errS := s.getSize(ctx, up.Bucket, up.Key)
	if errS != nil {
		return nil, fmt.Errorf("get size error %w", errS)
	}

	return &domain.UploadResponse{Size: size}, nil
}

func (s *S3Storage) getSize(ctx context.Context, bucket, key string) (int64, error) {
	res, err := s.Client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(s.getBucket(bucket)),
		Key:    aws.String(s.withPrefix(key)),
	})
	if err != nil {
		return 0, fmt.Errorf("head object error %w", err)
	}

	if res.ContentLength != nil {
		return *res.ContentLength, nil
	}

	return 0, nil
}

func (s *S3Storage) Download(ctx context.Context, dp *domain.DownloadParameters, writer io.Writer) error {
	s.Logger.Debugf("starting download %v", dp)

	downloader := manager.NewDownloader(s.Client, func(downloader *manager.Downloader) {
		downloader.Concurrency = 1
		downloader.PartBodyMaxRetries = 1 // todo: increase?
		downloader.LogInterruptedDownloads = true
	})

	fakeWriter := FakeWriterAt{w: writer}
	_, err := downloader.Download(ctx, &fakeWriter, &s3.GetObjectInput{
		Bucket: aws.String(s.getBucket(dp.Bucket)),
		Key:    aws.String(s.withPrefix(dp.Key)),
	})

	if errors.As(err, &responseError) && responseError.HTTPStatusCode() == http.StatusNotFound {
		return domain.ErrNotFound
	}

	if err != nil {
		return fmt.Errorf("downloading from storage error %w", err)
	}

	return nil
}

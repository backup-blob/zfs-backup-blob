package driver_test

import (
	"bytes"
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/backup-blob/zfs-backup-blob/internal/domain"
	config2 "github.com/backup-blob/zfs-backup-blob/internal/domain/config"
	"github.com/backup-blob/zfs-backup-blob/internal/domain/mocks"
	"github.com/backup-blob/zfs-backup-blob/internal/driver"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/testcontainers/testcontainers-go/modules/minio"
	"log"
	"strings"
	"testing"
)

func TestSpecS3(t *testing.T) {
	bucket := "test"
	fileName := "filename.txt"
	payload := "xoxoxoxoxoxo"

	minioContainer, closFunc := setupContainer(context.Background())
	s3Client := getClient(t, minioContainer)
	s3Conf := config2.S3Config{
		Bucket:       bucket,
		Region:       s3Client.Options().Region,
		BaseEndpoint: *s3Client.Options().BaseEndpoint,
		UsePathStyle: s3Client.Options().UsePathStyle,
		Prefix:       "/folder1/folder2",
	}
	ctx := context.Background()
	createBucket(bucket, s3Client)
	s := driver.NewS3Storage(s3Client, mocks.NewMockLogger())

	defer closFunc()

	Convey("Given a config is loaded with all fields", t, func() {
		Convey("When everything works as expected", func() {
			Convey("It returns a valid config-store and no error", func() {
				store, err := driver.NewS3StorageFromConfig(&s3Conf, mocks.NewMockLogger())

				So(store, ShouldNotBeNil)
				So(err, ShouldBeNil)
			})
		})
		Convey("When access key and secret are set", func() {
			Convey("It uses the access key and secret", func() {
				t.Setenv("AWS_ACCESS_KEY_ID", "nil")
				t.Setenv("AWS_SECRET_ACCESS_KEY", "nil")
				s3ConfCopy := s3Conf
				s3ConfCopy.AccessKey = minioContainer.Username
				s3ConfCopy.AccessSecret = minioContainer.Password
				reader := strings.NewReader(payload)

				store, err := driver.NewS3StorageFromConfig(&s3ConfCopy, mocks.NewMockLogger())
				So(err, ShouldBeNil)
				_, errU := store.Upload(ctx, &domain.UploadParameters{
					Bucket: bucket, Key: fileName,
				}, reader)
				So(errU, ShouldBeNil)
			})
		})
	})

	Convey("Given a file is upload", t, func() {
		Convey("When i delete the file", func() {
			Convey("It should be removed from the remote", func() {
				reader := strings.NewReader(payload)

				_, errUp := s.Upload(ctx, &domain.UploadParameters{
					Bucket: bucket, Key: fileName,
				}, reader)
				So(errUp, ShouldBeNil)

				errD := s.Delete(ctx, &domain.DeleteParameters{Bucket: bucket, Key: fileName})
				So(errD, ShouldBeNil)

				errDD := s.Delete(ctx, &domain.DeleteParameters{Bucket: bucket, Key: fileName})
				So(errDD, ShouldBeNil)
			})
		})
	})

	Convey("Given i upload/download a file", t, func() {
		Convey("When everything works as expected", func() {
			Convey("It should not have errors", func() {
				reader := strings.NewReader(payload)
				var buf bytes.Buffer

				res, errUp := s.Upload(ctx, &domain.UploadParameters{
					Bucket: bucket, Key: fileName,
				}, reader)
				So(errUp, ShouldBeNil)
				So(res.Size, ShouldEqual, len(payload))

				errDown := s.Download(ctx, &domain.DownloadParameters{
					Bucket: bucket, Key: fileName,
				}, &buf)

				So(errDown, ShouldBeNil)
				So(buf.String(), ShouldEqual, payload)
			})
		})
		Convey("When i use a prefix path", func() {
			Convey("It should not have errors", func() {
				sp := driver.S3Storage{Client: s3Client, Prefix: "/folder1/folder2", Logger: mocks.NewMockLogger()}
				reader := strings.NewReader(payload)
				var buf bytes.Buffer

				_, errUp := sp.Upload(ctx, &domain.UploadParameters{
					Bucket: bucket, Key: fileName,
				}, reader)
				So(errUp, ShouldBeNil)

				errDown := sp.Download(ctx, &domain.DownloadParameters{
					Bucket: bucket, Key: fileName,
				}, &buf)

				So(errDown, ShouldBeNil)
				So(buf.String(), ShouldEqual, payload)
			})
		})
	})
}

func createBucket(name string, s3Client *s3.Client) {
	_, err := s3Client.CreateBucket(context.TODO(), &s3.CreateBucketInput{
		Bucket: aws.String(name),
	})
	if err != nil {
		panic(err)
	}
}

func getClient(t *testing.T, container *minio.MinioContainer) *s3.Client {
	endpoint, err := container.ConnectionString(context.Background())
	if err != nil {
		panic(err)
	}
	t.Setenv("AWS_ACCESS_KEY_ID", container.Username)
	t.Setenv("AWS_SECRET_ACCESS_KEY", container.Password)
	sdkConfig, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		panic(err)
	}
	client := s3.NewFromConfig(sdkConfig, func(options *s3.Options) {
		options.BaseEndpoint = aws.String("http://" + endpoint)
		options.UsePathStyle = true
		options.Region = "us-east-1"
	})
	return client
}

func setupContainer(ctx context.Context) (*minio.MinioContainer, func()) {
	minioContainer, err := minio.Run(ctx, "minio/minio:RELEASE.2024-07-16T23-46-41Z")
	if err != nil {
		log.Fatalf("failed to start container: %s", err)
	}

	//Clean up the container
	return minioContainer, func() {
		if err := minioContainer.Terminate(ctx); err != nil {
			log.Fatalf("failed to terminate container: %s", err)
		}
	}
}

package main_test

import (
	"context"
	"errors"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/cucumber/godog"
	"github.com/testcontainers/testcontainers-go/modules/minio"
	"net/http"
	"os"
)

const StorageSettingsKey = "storageSettings"

type StorageSettings struct {
	Url      string
	Username string
	Password string
}

func setRegion() {
	os.Setenv("AWS_DEFAULT_REGION", "us-east-1")
}

func getClient(ctx context.Context) *s3.Client {
	sdkConfig, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		panic(err)
	}
	path, err := replacePlaceholder(ctx, "<path>")
	if err != nil {
		panic(err)
	}
	client := s3.NewFromConfig(sdkConfig, func(options *s3.Options) {
		options.BaseEndpoint = aws.String(path)
		options.UsePathStyle = true
		options.Region = "us-east-1"
	})
	return client
}

func getObject(ctx context.Context, key string, bucket string) error {
	client := getClient(ctx)
	_, err := client.HeadObject(ctx, &s3.HeadObjectInput{Bucket: aws.String(bucket), Key: aws.String(key)})
	return err
}

var responseError interface {
	HTTPStatusCode() int
}

func objectDoesNotExist(ctx context.Context, key string, bucket string) error {
	client := getClient(ctx)
	_, err := client.HeadObject(ctx, &s3.HeadObjectInput{Bucket: aws.String(bucket), Key: aws.String(key)})
	if errors.As(err, &responseError) && responseError.HTTPStatusCode() == http.StatusNotFound {
		return nil
	}
	if err != nil {
		return err
	}
	return errors.New("error exists")
}

func createBucket(ctx context.Context, name string) error {
	client := getClient(ctx)
	_, err := client.CreateBucket(ctx, &s3.CreateBucketInput{
		Bucket: aws.String(name),
	})
	return err
}

func createStorage(ctx context.Context) (context.Context, error) {
	con, err := minio.Run(ctx, "minio/minio:RELEASE.2024-07-16T23-46-41Z")
	if err != nil {
		return ctx, err
	}

	conStr, err := con.ConnectionString(ctx)

	if err != nil {
		return ctx, err
	}

	settings := StorageSettings{
		Url:      "http://" + conStr,
		Username: con.Username,
		Password: con.Password,
	}

	go func() {
		select {
		case <-ctx.Done():
			con.Terminate(ctx)
			return
		}
	}()
	return context.WithValue(ctx, StorageSettingsKey, settings), err
}

func setUsername(ctx context.Context, envName string) error {
	v, ok := ctx.Value(StorageSettingsKey).(StorageSettings)
	if ok {
		os.Setenv(envName, v.Username)
		return nil
	}
	return errors.New("could not set username ")
}

func setPassword(ctx context.Context, envName string) error {
	v, ok := ctx.Value(StorageSettingsKey).(StorageSettings)
	if ok {
		os.Setenv(envName, v.Password)
		return nil
	}
	return errors.New("could not set password ")
}

func setBaseurl(ctx context.Context) (context.Context, error) {
	v, ok := ctx.Value(StorageSettingsKey).(StorageSettings)
	if ok {
		return setPlaceholder(ctx, "<path>", v.Url)
	}
	return ctx, errors.New("could not set baseurl ")
}

func InitStorage(sc *godog.ScenarioContext) {
	sc.Step(`the s3 base url is set as placeholder '<path>'`, setBaseurl)
	sc.Step(`the s3 key secret is set at \'([_a-zA-Z-\/0-9@]+)\'`, setPassword)
	sc.Step(`the s3 key id is set at \'([_a-zA-Z-\/0-9@]+)\'`, setUsername)
	sc.Step(`create storage`, createStorage)
	sc.Step(`a bucket with the name \'([<>_a-zA-Z-\/0-9@\.]+)\' exists`, createBucket)
	sc.Step(`the key \'([\+\s<>_a-zA-Z-\/0-9@\.]+)\' does not exists in bucket \'([\+\s<>_a-zA-Z-\/0-9@\.]+)\'`, objectDoesNotExist)
	sc.Step(`the key \'([\+\s<>_a-zA-Z-\/0-9@\.]+)\' exists in bucket \'([\+\s<>_a-zA-Z-\/0-9@\.]+)\'`, getObject)
	sc.Step(`^the \'AWS_DEFAULT_REGION\' env var is set to \'us-east-1\'$`, setRegion)
}
